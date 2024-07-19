package schemas

// ErrorResponse represents a standard error response
// @Description Represents a standard error response for the API
type ErrorResponse struct {
	Code    int    `json:"code"`    // Error code
	Message string `json:"message"` // Error message
}
