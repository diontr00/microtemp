package model

// use within translate to translate  request payload
type FieldError interface {
	// Return the validation tag that failed for example "required or lt"
	Tag() string
	// Returns the field name
	Field() string
	// Returns the param that send  in case needed for creating the message
	Param() string
}

// Request Error Response
type ErrorResponse struct {
	Field string `json:"field,omitempty"`
	Error string `json:"error"`
}
