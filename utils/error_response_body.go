package utils

type ErrorResponseBody struct {
	Detail         string              `json:"detail"`
	FieldErrors    map[string][]string `json:"field_errors"`
	NonFieldErrors []string            `json:"non_field_errors"`
}

type VaultsErrorResponseBody struct {
	ClientOperation string `json:"client_operation"`
	Message         string `json:"message"`
	ContextString   string `json:"context_string"`
	RequestBody     string `json:"request_body"`
	Detail          string `json:"detail"`
}
