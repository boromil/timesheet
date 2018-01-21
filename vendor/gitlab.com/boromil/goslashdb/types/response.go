package types

import "strings"

// CreateResponse represents an response on create/POST request
type CreateResponse struct {
	ID   string
	Body string
}

// NewCreateResponse is the default constructor for the create response object
func NewCreateResponse(input string) CreateResponse {
	var id string
	tmp := strings.Split(input, "/")
	if len(tmp) > 0 {
		id = tmp[len(tmp)-1]
	}
	return CreateResponse{
		ID:   id,
		Body: input,
	}
}
