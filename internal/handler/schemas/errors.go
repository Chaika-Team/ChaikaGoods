package schemas

import (
	_ "github.com/Chaika-Team/ChaikaGoods/docs"
)

// ErrorResponse представляет собой стандартный ответ об ошибке
type ErrorResponse struct {
	Code    int    `json:"code"`    // Код ошибки
	Message string `json:"message"` // Сообщение об ошибке
}

// DetailedErrorResponse Новая структура для детализированных ошибок
type DetailedErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason,omitempty"` // Дополнительное объяснение
}
