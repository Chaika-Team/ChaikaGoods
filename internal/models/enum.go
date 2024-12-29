package models

// OperationType описывает тип операции над продуктом.
type OperationType int

const (
	OperationTypeUnknown OperationType = 0
	OperationTypeInsert  OperationType = 1
	OperationTypeUpdate  OperationType = 2
	OperationTypeDelete  OperationType = 3
)
