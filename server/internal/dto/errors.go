package dto

type GenericError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type InternalError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
