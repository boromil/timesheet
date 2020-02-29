package slashdb

import (
	"context"
	"fmt"

	"gitlab.com/boromil/goslashdb/types"
)

// QueryConfigManager - represents CRUD operations for the QueryConfig request
type QueryConfigManager interface {
	QueryConfigs(ctx context.Context) (map[string]types.QueryConfig, error)
	QueryConfig(ctx context.Context, id string) (types.QueryConfig, error)
	CreateQueryConfig(ctx context.Context, q types.QueryConfig) error
	UpdateQueryConfig(ctx context.Context, id string, q types.QueryConfig) error
	DeleteQueryConfig(ctx context.Context, id string) error
}

// QueryConfigs retrives all the custom quetry configs
func (s *Service) QueryConfigs(ctx context.Context) (map[string]types.QueryConfig, error) {
	var data = map[string]types.QueryConfig{}
	if err := s.Get(ctx, &Request{Kind: "/querydef"}, &data); err != nil {
		return nil, fmt.Errorf("error retriving query configs: %w", err)
	}
	return data, nil
}

// QueryConfig retrives a single custom quetry config
func (s *Service) QueryConfig(ctx context.Context, id string) (types.QueryConfig, error) {
	var data = types.QueryConfig{}
	// Name: id is a workaround, as the config API
	// doesn't follow the standar request data paters
	if err := s.Get(ctx, &Request{Kind: "/querydef", Parts: []Part{{Name: id}}}, &data); err != nil {
		return types.QueryConfig{}, fmt.Errorf("error retriving query config: %w", err)
	}
	return data, nil
}

// CreateQueryConfig creates a new custom quetry config
func (s *Service) CreateQueryConfig(ctx context.Context, q types.QueryConfig) error {
	_, err := s.Create(ctx, &Request{Kind: "/querydef"}, q)
	return err
}

// UpdateQueryConfig updates an existing custom quetry config
func (s *Service) UpdateQueryConfig(ctx context.Context, id string, q types.QueryConfig) error {
	// Name: id is a workaround, as the config API
	// doesn't follow the standar request data paters
	return s.Update(ctx, &Request{Kind: "/querydef", Parts: []Part{{Name: id}}}, q)

}

// DeleteQueryConfig deletes a single custom quetry config
func (s *Service) DeleteQueryConfig(ctx context.Context, id string) error {
	return s.Delete(ctx, &Request{Kind: "/querydef", Parts: []Part{{Name: id}}})

}
