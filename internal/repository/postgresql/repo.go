package postgresql

import (
	"ChaikaGoods/internal/models"
	"ChaikaGoods/internal/myerr"
	"ChaikaGoods/internal/repository"
	"ChaikaGoods/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/jackc/pgx/v5"
	"strconv"
)

type GoodsPGRepository struct {
	client Client
	log    log.Logger
}

// NewGoodsRepository создает новый экземпляр GoodsRepository.
func NewGoodsRepository(client Client, logger log.Logger) repository.GoodsRepository {
	return &GoodsPGRepository{
		client: client,
		log:    logger,
	}
}

// GetProductByID возвращает продукт по его ID.
func (r *GoodsPGRepository) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	sql := `SELECT id, name, description, price, imageurl, sku FROM public.product WHERE id = $1;`
	row := r.client.QueryRow(ctx, sql, id)

	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.SKU)
	if errors.Is(err, pgx.ErrNoRows) {
		return p, &myerr.NotFound{ID: fmt.Sprintf("%d", id)}
	} else if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get product by ID: %v", err))
		return p, err
	}

	return p, nil
}

// GetAllProducts возвращает список всех продуктов.
func (r *GoodsPGRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	sql := `SELECT id, name, description, price, imageurl, sku FROM public.product;`
	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get all products: %v", err))
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.SKU); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan product: %v", err))
			continue
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed during rows iteration: %v", err))
		return nil, err
	}
	return products, nil
}

// AddQueryToCreateProduct добавляет запрос на создание продукта в базе данных.
func (r *GoodsPGRepository) AddQueryToCreateProduct(ctx context.Context, data *map[string]interface{}) (changeID int64, err error) {
	// добавляем новое изменение в базу
	newValue, err := json.Marshal(data)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to marshal product data: %v", err))
		return 0, err
	}
	sql := `INSERT INTO public.changes(operation, new_value) VALUES ( $1, $2) RETURNING change_id;` //version_id подставляется автоматически
	jsonValue := string(newValue)
	err = r.client.QueryRow(ctx, sql, models.OperationTypeInsert, jsonValue).Scan(&changeID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to add change: %v", err))
		return 0, err

	}
	return changeID, nil
}

// AddQueryToUpdateProduct добавление запроса на обновления продукта в базе данных.
func (r *GoodsPGRepository) AddQueryToUpdateProduct(ctx context.Context, data *map[string]interface{}) (changeID int64, err error) {
	newValue, err := json.Marshal(data)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to marshal product data: %v", err))
		return 0, err
	}
	jsonValue := string(newValue)
	sql := `INSERT INTO public.changes(operation, new_value) VALUES ( $1, $2) RETURNING change_id;` //version_id подставляется автоматически
	err = r.client.QueryRow(ctx, sql, models.OperationTypeUpdate, jsonValue).Scan(&changeID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to add change: %v", err))
		return 0, err
	}
	return changeID, nil
}

// AddQueryToDeleteProduct удаляет продукт из базы данных по его ID.
func (r *GoodsPGRepository) AddQueryToDeleteProduct(ctx context.Context, id int64) (changeID int64, err error) {
	sql := `INSERT INTO public.changes(operation, new_value) VALUES ($1, $2) RETURNING change_id;` //version_id подставляется автоматически
	str := "{\"id\":" + strconv.FormatInt(id, 10) + "}"                                            // Преобразование id в строку
	err = r.client.QueryRow(ctx, sql, models.OperationTypeDelete, str).Scan(&changeID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to add change: %v", err))
		return 0, err
	}
	return changeID, nil
}

func (r *GoodsPGRepository) ApplyChanges(ctx context.Context, version *models.Version) error {
	//TODO: Переписать метод, упростить (high complexity 160%)

	// Взять все изменения, которые не были применены
	sql := `SELECT change_id, new_value, operation FROM public.changes WHERE version_id = $1 AND considered = FALSE;`
	rows, err := r.client.Query(ctx, sql, &version.VersionID)
	if err != nil {
		return fmt.Errorf("failed to get changes: %v", err)
	}
	defer rows.Close()

	// Начало транзакции
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Завершение транзакции в случае ошибки или при успешном выполнении
	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				_ = r.log.Log("error", fmt.Sprintf("Failed to rollback transaction: %v", err))
			}
			panic(p) // Re-throw panic after rollback
		} else if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				_ = r.log.Log("error", fmt.Sprintf("Failed to rollback transaction: %v", err))
			} // Rollback transaction if error occurred
		} else {
			err = tx.Commit(ctx) // Commit transaction if all is good
		}
	}()

	// Применить все изменения
	for rows.Next() {
		var changeID int
		var productDataRaw *json.RawMessage
		var operation models.OperationType
		if err := rows.Scan(&changeID, &productDataRaw, &operation); err != nil {
			return fmt.Errorf("failed to scan change: %v", err)
		}
		// Распаковка нового значения
		var productData map[string]interface{}
		if err := json.Unmarshal(*productDataRaw, &productData); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("failed to unmarshal new value: %v", err))
			return err
		}
		// Обработка изменений в зависимости от типа операции
		switch operation {
		case models.OperationTypeInsert:
			err = r.insertProduct(ctx, &productData)
		case models.OperationTypeUpdate:
			err = r.updateProduct(ctx, &productData)
		case models.OperationTypeDelete:
			err = r.deleteProduct(ctx, &productData)
		default:
			return fmt.Errorf("unknown operation type: %v", operation)
		}
		if err != nil {
			return fmt.Errorf("failed to apply change: %v", err)
		}
		// Отметить изменение как применённое
		_, err = tx.Exec(ctx, `UPDATE public.changes SET considered = TRUE WHERE change_id = $1`, changeID)
		if err != nil {
			return fmt.Errorf("failed to mark change as considered: %v", err)
		}
	}
	// Пометить версию как применённую и не dev
	_, err = tx.Exec(ctx, `UPDATE public.version SET applied = TRUE, is_dev = FALSE WHERE version_id = $1`, version.VersionID)
	if err != nil {
		return fmt.Errorf("failed to mark version as applied: %v", err)
	}
	return nil
}

// CreateNewDevVersion создает новую версию базы данных продуктов для разработки.
func (r *GoodsPGRepository) CreateNewDevVersion(ctx context.Context) (models.Version, error) {
	sql := `INSERT INTO public.version DEFAULT VALUES  RETURNING version_id;`
	var v models.Version
	err := r.client.QueryRow(ctx, sql).Scan(&v.VersionID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to create new dev version: %v", err))
		return models.Version{}, err
	}
	return v, nil
}

// GetAllChanges возвращает все изменения в базе данных продуктов за конкретную версию.
func (r *GoodsPGRepository) GetAllChanges(ctx context.Context, version models.Version) ([]models.Change, error) {
	sql := `SELECT change_id, operation, new_value, change_timestamp, considered FROM public.changes WHERE version_id = $1;`
	rows, err := r.client.Query(ctx, sql, version.VersionID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get all changes: %v", err))
		return nil, err
	}
	defer rows.Close()
	var changes []models.Change
	for rows.Next() {
		var c models.Change
		if err := rows.Scan(&c.ID, &c.OperationType, &c.NewValue, &c.ChangeTimestamp, &c.Considered); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan change: %v", err))
			continue
		}
		var product models.Product
		if err := json.Unmarshal(c.NewValue, &product); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to unmarshal new value: %v", err))
			continue
		}
		c.VersionID = version.VersionID
		changes = append(changes, c)
	}
	return changes, nil
}

// GetCurrentDevVersion возвращает текущую версию базы данных продуктов к которой привязываются новые изменения.
func (r *GoodsPGRepository) GetCurrentDevVersion(ctx context.Context) (models.Version, error) {
	sql := `SELECT version_id, is_dev FROM public.version WHERE is_dev = TRUE ORDER BY creation_date DESC LIMIT 1;`
	var v models.Version
	err := r.client.QueryRow(ctx, sql).Scan(&v.VersionID, &v.IsDev)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get current dev version: %v", err))
		return models.Version{}, err
	}
	return v, nil
}

// GetCurrentActualVersion возвращает текущую версию базы данных продуктов c актуальными данными
func (r *GoodsPGRepository) GetCurrentActualVersion(ctx context.Context) (models.Version, error) {
	sql := `SELECT version_id, is_dev FROM public.version WHERE is_dev = FALSE ORDER BY creation_date DESC LIMIT 1;`
	var v models.Version
	err := r.client.QueryRow(ctx, sql).Scan(&v.VersionID, &v.IsDev)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get current actual version: %v", err))
		return models.Version{}, err
	}
	return v, nil
}

// GetVersionsBetween возвращает все версии базы данных продуктов между двумя версиями.
func (r *GoodsPGRepository) GetVersionsBetween(ctx context.Context, from, to int) ([]models.Version, error) {
	sql := `SELECT version_id, creation_date, is_dev, applied FROM public.version WHERE version_id > $1 AND version_id < $2;`
	//TODO: Переписать запрос, не подходит для крайних значений, когда int закончится
	rows, err := r.client.Query(ctx, sql, from, to)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get versions between: %v", err))
		return nil, err
	}
	defer rows.Close()
	var versions []models.Version
	for rows.Next() {
		var v models.Version
		if err := rows.Scan(&v.VersionID, &v.CreationDate, &v.IsDev, &v.Applied); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan version: %v", err))
			continue
		}
		versions = append(versions, v)
	}
	return versions, nil

}

// DeleteChange удаляет изменение из базы данных по его ID.
func (r *GoodsPGRepository) DeleteChange(ctx context.Context, id int64) error {
	sql := `DELETE FROM public.changes WHERE change_id = $1;`
	_, err := r.client.Exec(ctx, sql, id)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to delete change: %v", err))
		return err
	}
	return nil
}

// Package queries

// GetPackageByID получает полную информацию о пакете, включая его содержимое.
func (r *GoodsPGRepository) GetPackageByID(ctx context.Context, p *models.Package) ([]models.PackageContent, error) {
	sqlPackage := `SELECT packageid, packagename, description FROM public."package" WHERE packageid = $1;`
	err := r.client.QueryRow(ctx, sqlPackage, p.ID).Scan(&p.ID, &p.PackageName, &p.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &myerr.NotFound{ID: fmt.Sprintf("%d", p.ID)}
		}
		_ = r.log.Log("error", fmt.Sprintf("Failed to get package by ID: %v", err))
		return nil, err
	}

	sqlContents := `SELECT packagecontentid, packageid, productid, quantity FROM public.packagecontent WHERE packageid = $1;`
	rows, err := r.client.Query(ctx, sqlContents, p.ID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get package contents: %v", err))
		return nil, err
	}
	defer rows.Close()

	var contents []models.PackageContent
	for rows.Next() {
		var c models.PackageContent
		if err := rows.Scan(&c.ID, &c.PackageID, &c.ProductID, &c.Quantity); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan package content: %v", err))
			return nil, err
		}
		contents = append(contents, c)
	}

	return contents, nil
}

// GetProductsByPackageID получает список продуктов в определенном пакете.
func (r *GoodsPGRepository) GetProductsByPackageID(ctx context.Context, p *models.Package) ([]models.PackageContent, error) {
	sql := `SELECT packagecontentid, packageid, productid, quantity FROM public.packagecontent WHERE packageid = $1;`
	rows, err := r.client.Query(ctx, sql, p.ID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get products by package ID: %v", err))
		return nil, err
	}
	defer rows.Close()

	var contents []models.PackageContent
	for rows.Next() {
		var c models.PackageContent
		if err := rows.Scan(&c.ID, &c.PackageID, &c.ProductID, &c.Quantity); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan package content: %v", err))
			continue
		}
		contents = append(contents, c)
	}

	if len(contents) == 0 {
		return nil, &myerr.NotFound{ID: fmt.Sprintf("%d", p.ID)}
	}

	return contents, nil
}

// ListPackages возвращает список всех пакетов.
func (r *GoodsPGRepository) ListPackages(ctx context.Context) ([]models.Package, error) {
	sql := `SELECT packageid, packagename, description FROM public.package;`
	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to list packages: %v", err))
		return nil, err
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var p models.Package
		if err := rows.Scan(&p.ID, &p.PackageName, &p.Description); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan package: %v", err))
			continue
		}
		packages = append(packages, p)
	}

	return packages, nil
}

// CreatePackage добавляет новый пустой пакет в базу данных.
func (r *GoodsPGRepository) CreatePackage(ctx context.Context, pkg *models.Package) error {
	sql := `INSERT INTO public."package" (packagename, description) VALUES ($1, $2) RETURNING packageid;`
	err := r.client.QueryRow(ctx, sql, &pkg.PackageName, &pkg.Description).Scan(&pkg.ID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to create package: %v", err))
		return err
	}
	return nil
}

// AddProductToPackage добавляет продукты в пакет.
func (r *GoodsPGRepository) AddProductToPackage(ctx context.Context, packageID int64, products []models.PackageContent) error {
	sql := `INSERT INTO public.packagecontent (packageid, productid, quantity) VALUES ($1, $2, $3);`
	batch := &pgx.Batch{}

	for _, product := range products {
		batch.Queue(sql, packageID, product.ProductID, product.Quantity)
	}

	br := r.client.SendBatch(ctx, batch)
	defer func(br pgx.BatchResults) {
		err := br.Close()
		if err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to close batch results: %v", err))
		}
	}(br)
	_, err := br.Exec()
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to add products to package: %v", err))
		return err
	}
	return nil
}

// DeletePackage удаляет пакет и его содержимое.
func (r *GoodsPGRepository) DeletePackage(ctx context.Context, packageID int64) error {
	sqlDeleteContents := `DELETE FROM public.packagecontent WHERE packageid = $1;`
	sqlDeletePackage := `DELETE FROM public."package" WHERE packageid = $1;`

	tx, err := r.client.Begin(ctx)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to begin transaction: %v", err))
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to rollback transaction: %v", err))
		}
	}(tx, ctx)

	if _, err := tx.Exec(ctx, sqlDeleteContents, packageID); err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to delete package contents: %v", err))
		return err
	}
	if _, err := tx.Exec(ctx, sqlDeletePackage, packageID); err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to delete package: %v", err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to commit transaction: %v", err))
		return err
	}
	return nil
}

func (r *GoodsPGRepository) SearchPacket(ctx context.Context, searchString string, quantity int64, offset int64) ([]models.Package, error) {
	sql := `SELECT packageid, packagename, description FROM public."package" WHERE packagename LIKE $1 OR description LIKE $1 LIMIT $2 OFFSET $3;`
	rows, err := r.client.Query(ctx, sql, searchString, quantity, offset)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to search package: %v", err))
		return nil, err
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var p models.Package
		if err := rows.Scan(&p.ID, &p.PackageName, &p.Description); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan package: %v", err))
			continue
		}
		packages = append(packages, p)
	}
	return packages, nil
}

// insertProduct добавляет новый продукт в базу данных, только для внутреннего использования.
func (r *GoodsPGRepository) insertProduct(ctx context.Context, data *map[string]interface{}) error {
	// Verify that all keys in the map correspond to the fields of the models.Product structure
	err := utils.VerifyMapFields[models.Product](*data)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to verify product fields: %v", err))
		return err
	}
	// Generate sql query with all fields in the map
	sql := `INSERT INTO public.product (`
	values := `) VALUES (`
	var args []interface{}
	idx := 1
	for key, value := range *data {
		sql += key + ", "
		values += "$" + strconv.Itoa(idx) + ", "
		args = append(args, value) // Добавление значения в список аргументов
		idx++
	}
	sql = sql[:len(sql)-2] + values[:len(values)-2] + ");"
	_, err = r.client.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to insert product: %v", err)
	}
	return nil
}

// updateProduct обновляет информацию о продукте в базе данных, только для внутреннего использования.
func (r *GoodsPGRepository) updateProduct(ctx context.Context, data *map[string]interface{}) error {
	// Verify that all keys in the map correspond to the fields of the models.Product structure
	err := utils.VerifyMapFields[models.Product](*data)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to verify product fields: %v", err))
		return err
	}
	// Generate sql query with all fields in the map
	sql := `UPDATE public.product SET `
	var args []interface{}
	idx := 1
	for key, value := range *data {
		if key == "id" {
			continue
		}
		sql += key + " = $" + strconv.Itoa(idx) + ", "
		args = append(args, value) // Добавление значения в список аргументов
		idx++
	}
	if len(args) == 0 {
		return fmt.Errorf("no fields to update")
	}
	args = append(args, (*data)["id"]) // Добавление id в список аргументов
	sql = sql[:len(sql)-2] + " WHERE id = $" + strconv.Itoa(idx) + ";"
	_, err = r.client.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update product: %v", err)
	}
	return nil
}

// deleteProduct удаляет продукт из базы данных по его ID, только для внутреннего использования.
func (r *GoodsPGRepository) deleteProduct(ctx context.Context, data *map[string]interface{}) error {
	// get id
	rawId, ok := (*data)["id"]
	if !ok {
		return fmt.Errorf("failed to get product id")
	}

	var id = int64(rawId.(float64))

	sql := `DELETE FROM public.product WHERE id = $1;`
	_, err := r.client.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}
	return nil
}
