package models

// OperationType описывает тип операции над продуктом.
type OperationType int

const (
	OperationTypeUnknown OperationType = 0
	OperationTypeInsert                = 1
	OperationTypeUpdate                = 2
	OperationTypeDelete                = 3
)
