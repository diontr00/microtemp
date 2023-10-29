package model

type FieldError interface {
	// Return the validation tag that failed for example "required or lt"
	Tag() string
	// Returns the field name
	Field() string
	// Returns the param that send  in case needed for creating the message
	Param() string
}

type RequestError struct {
	Field string `json:"field"`
	Msg   string `json:"error_msg"`
}
