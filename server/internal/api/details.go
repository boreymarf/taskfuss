package api

type FieldErrorDetail struct {
	Field    string `json:"field"`
	Expected string `json:"expected"`
	Message  string `json:"message,omitempty"`
}
