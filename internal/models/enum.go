package models

// OperationType описывает тип операции над продуктом.
type OperationType int

const (
	OperationTypeInsert OperationType = 0
	OperationTypeUpdate               = 1
	OperationTypeDelete               = 2
)
