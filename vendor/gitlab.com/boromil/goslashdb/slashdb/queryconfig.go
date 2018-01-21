package slashdb

import (
	"context"

	"github.com/pkg/errors"
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
func (s *service) QueryConfigs(ctx context.Context) (map[string]types.QueryConfig, error) {
	var data = map[string]types.QueryConfig{}
	if err := s.Get(ctx, &Request{Kind: "/querydef"}, &data); err != nil {
		return nil, errors.Wrap(err, "error retriving query configs")
	}
	return data, nil
}

// QueryConfig retrives a single custom quetry config
func (s *service) QueryConfig(ctx context.Context, id string) (types.QueryConfig, error) {
	var data = types.QueryConfig{}
	// Name: id is a workaround, as the config API
	// doesn't follow the standar request data paters
	if err := s.Get(ctx, &Request{Kind: "/querydef", Parts: []Part{{Name: id}}}, &data); err != nil {
		return types.QueryConfig{}, errors.Wrap(err, "error retriving query config")
	}
	return data, nil
}

// QueryConfig creates a new custom quetry config
func (s *service) CreateQueryConfig(ctx context.Context, q types.QueryConfig) error {
	_, err := s.Create(ctx, &Request{Kind: "/querydef"}, q)
	return err
}

// QueryConfig updates an existing custom quetry config
func (s *service) UpdateQueryConfig(ctx context.Context, id string, q types.QueryConfig) error {
	// Name: id is a workaround, as the config API
	// doesn't follow the standar request data paters
	return s.Update(ctx, &Request{Kind: "/querydef", Parts: []Part{{Name: id}}}, q)

}

// QueryConfig deletes a single custom quetry config
func (s *service) DeleteQueryConfig(ctx context.Context, id string) error {
	return s.Delete(ctx, &Request{Kind: "/querydef", Parts: []Part{{Name: id}}})

}
