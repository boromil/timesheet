package types

// APIError - represets SlashDBs API error
type APIError struct {
	HTTPCode    int    `json:"http_code" xml:"http_code"`
	Description string `json:"description" xml:"description"`
}
