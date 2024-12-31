package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// extractID извлекает и преобразует ID из переменных пути, например, /api/v1/product/{id}.
func extractID(req *http.Request, paramName string) (int64, error) {
	vars := mux.Vars(req)
	idStr, ok := vars[paramName]
	if !ok {
		return 0, errors.New("missing " + paramName + " in path")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid " + paramName)
	}

	return id, nil
}
