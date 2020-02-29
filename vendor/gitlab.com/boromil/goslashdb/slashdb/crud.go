package slashdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"gitlab.com/boromil/goslashdb/types"
)

// CRUDer an interface representing CRUD operations
type CRUDer interface {
	Get(
		ctx context.Context,
		sdbReq fmt.Stringer,
		container interface{},
	) error
	Create(
		ctx context.Context,
		sdbReq fmt.Stringer,
		payload interface{},
	) (types.CreateResponse, error)
	Update(
		ctx context.Context,
		sdbReq fmt.Stringer,
		payload interface{},
	) error
	Delete(
		ctx context.Context,
		sdbReq fmt.Stringer,
	) error
}

// Get gets resources using GET method
func (s *Service) Get(
	ctx context.Context,
	sdbReq fmt.Stringer,
	container interface{},
) error {
	method := http.MethodGet
	endpoint := fmt.Sprintf("%s%s", s.host, sdbReq)

	s.echoRequest(method, endpoint, nil)

	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return fmt.Errorf("error creating a request: %w", err)
	}
	req.Header.Set(s.apiKeyName, s.apiKeyValue)

	resp, err := s.client.Do(req)
	if resp != nil {
		defer func() {
			if _, err = io.Copy(ioutil.Discard, resp.Body); err != nil {
				log.Printf("error copying the resp.Body content: %v\n", err)
			}
			if err = resp.Body.Close(); err != nil {
				log.Printf("error closing the resp.Body: %v\n", err)
			}
		}()
	}
	if err != nil {
		return fmt.Errorf("error doing the request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiError types.APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err == nil {
			apiError.Description = ": " + apiError.Description
		}

		return fmt.Errorf("failed to get request%s", apiError.Description)
	}

	if err := json.NewDecoder(resp.Body).Decode(container); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	return nil
}

// Create creates resources using POST method
func (s *Service) Create(
	ctx context.Context,
	sdbReq fmt.Stringer,
	payload interface{},
) (types.CreateResponse, error) {
	data := []byte{}
	buf := bytes.NewBuffer(data)
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return types.CreateResponse{}, fmt.Errorf("error encoding data: %w", err)
	}

	method := http.MethodPost
	endpoint := fmt.Sprintf("%s%s", s.host, sdbReq)

	s.echoRequest(method, endpoint, data)

	req, err := http.NewRequestWithContext(ctx, method, endpoint, buf)
	if err != nil {
		return types.CreateResponse{}, fmt.Errorf("error creating a request: %w", err)
	}
	req.Header.Set(s.apiKeyName, s.apiKeyValue)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if resp != nil {
		defer func() {
			if _, err = io.Copy(ioutil.Discard, resp.Body); err != nil {
				log.Printf("error copying the resp.Body content: %v\n", err)
			}
			if err = resp.Body.Close(); err != nil {
				log.Printf("error closing the resp.Body: %v\n", err)
			}
		}()
	}
	if err != nil {
		return types.CreateResponse{}, fmt.Errorf("error doing the request: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		var apiError types.APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err == nil {
			apiError.Description = ": " + apiError.Description
		}

		return types.CreateResponse{}, fmt.Errorf("failed to create conten%s", apiError.Description)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.CreateResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	return types.NewCreateResponse(string(b)), nil
}

// Update updates resources using PUT method
func (s *Service) Update(
	ctx context.Context,
	sdbReq fmt.Stringer,
	payload interface{},
) error {
	data := []byte{}
	buf := bytes.NewBuffer(data)
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return fmt.Errorf("error encoding data: %w", err)
	}

	method := http.MethodPut
	endpoint := fmt.Sprintf("%s%s", s.host, sdbReq)

	s.echoRequest(method, endpoint, data)

	req, err := http.NewRequestWithContext(ctx, method, endpoint, buf)
	if err != nil {
		return fmt.Errorf("error creating a request: %w", err)
	}
	req.Header.Set(s.apiKeyName, s.apiKeyValue)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if resp != nil {
		defer func() {
			if _, err = io.Copy(ioutil.Discard, resp.Body); err != nil {
				log.Printf("error copying the resp.Body content: %v\n", err)
			}
			if err = resp.Body.Close(); err != nil {
				log.Printf("error closing the resp.Body: %v\n", err)
			}
		}()
	}
	if err != nil {
		return fmt.Errorf("error doing the request: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		var apiError types.APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err == nil {
			apiError.Description = ": " + apiError.Description
		}

		return fmt.Errorf("failed to update conten%s", apiError.Description)
	}

	return nil
}

// Delete deletes resources using DELETE method
func (s *Service) Delete(
	ctx context.Context,
	sdbReq fmt.Stringer,
) error {
	method := http.MethodDelete
	endpoint := fmt.Sprintf("%s%s", s.host, sdbReq)

	s.echoRequest(method, endpoint, nil)

	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return fmt.Errorf("error creating a request: %w", err)
	}
	req.Header.Set(s.apiKeyName, s.apiKeyValue)

	resp, err := s.client.Do(req)
	if resp != nil {
		defer func() {
			if _, err = io.Copy(ioutil.Discard, resp.Body); err != nil {
				log.Printf("error copying the resp.Body content: %v\n", err)
			}
			if err = resp.Body.Close(); err != nil {
				log.Printf("error closing the resp.Body: %v\n", err)
			}
		}()
	}
	if err != nil {
		return fmt.Errorf("error doing a request: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		var apiError types.APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err == nil {
			apiError.Description = ": " + apiError.Description
		}

		return fmt.Errorf("failed to delete content%s", apiError.Description)
	}

	return nil
}
