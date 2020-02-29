package slashdb

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Doer - a simple interface for a http.Client
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Service - main SlashDB API service
type Service struct {
	host        string
	apiKeyName  string
	apiKeyValue string

	refIDPrefix string
	resources   map[string]string

	echoMode bool

	client Doer
}

// NewService - returns a new instance of a SlashDB service
func NewService(
	host, apiKeyName, apiKeyValue, refIDPrefix string,
	echoMode bool,
	httpClient Doer,
) (*Service, error) {
	if _, err := url.Parse(host); err != nil {
		return nil, fmt.Errorf("malformed SlashDB host URL: %w", err)
	}

	return &Service{
		host:        host,
		apiKeyName:  apiKeyName,
		apiKeyValue: apiKeyValue,
		refIDPrefix: refIDPrefix,
		resources:   map[string]string{},
		echoMode:    echoMode,
		client:      httpClient,
	}, nil
}

// Init - initializes the service i.e. get the base resource mapping
func (s *Service) Init(ctx context.Context) error {
	if err := s.Get(ctx, NewDataRequest(""), &s.resources); err != nil {
		return fmt.Errorf("error retriving resource info: %w", err)
	}

	// rm refIDPrefix from resources map
	delete(s.resources, s.refIDPrefix)

	return nil
}

// Resources - returns a copy of the services resources map
func (s *Service) Resources() map[string]string {
	r := make(map[string]string, len(s.resources))
	for k, v := range s.resources {
		r[k] = v
	}
	return r
}

func (s *Service) echoRequest(method string, endpoint string, body []byte) {
	if !s.echoMode {
		return
	}
	msg := fmt.Sprint("\ndoing a: ", method, "\nrequest to: ", endpoint)
	if body != nil {
		msg = msg + "\nwith body: " + string(body)
	}
	log.Println(msg)
}
