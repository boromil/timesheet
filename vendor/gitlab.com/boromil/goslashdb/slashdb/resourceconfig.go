package slashdb

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/boromil/goslashdb/types"
)

// ResourceConfigManager - represents CRUD operations for the ResourceConfig request
type ResourceConfigManager interface {
	ResourceConfigs(ctx context.Context) (map[string]types.ResourceConfig, error)
	ResourceConfig(ctx context.Context, id string) (types.ResourceConfig, error)
	CreateResourceConfig(ctx context.Context, d types.ResourceConfig) error
	UpdateResourceConfig(ctx context.Context, id string, d types.ResourceConfig) error
	DeleteResourceConfig(ctx context.Context, id string) error
}

// ResourceConfigs retrives all the data resource configs
func (s *service) ResourceConfigs(ctx context.Context) (map[string]types.ResourceConfig, error) {
	var data = map[string]types.ResourceConfig{}
	if err := s.Get(ctx, &Request{Kind: "/dbdef"}, &data); err != nil {
		return nil, errors.Wrap(err, "error retriving resource configs")
	}
	return data, nil
}

// ResourceConfig retrives a single data resource config
func (s *service) ResourceConfig(ctx context.Context, id string) (types.ResourceConfig, error) {
	var data = types.ResourceConfig{}
	// Name: id is a workaround, as the config API
	// doesn't follow the standar request data paters
	if err := s.Get(ctx, &Request{Kind: "/dbdef", Parts: []Part{{Name: id}}}, &data); err != nil {
		return types.ResourceConfig{}, errors.Wrap(err, "error retriving resource config")
	}
	return data, nil
}

// CreateResourceConfig creates a new data resource config
func (s *service) CreateResourceConfig(ctx context.Context, d types.ResourceConfig) error {
	_, err := s.Create(ctx, &Request{Kind: "/dbdef"}, d)
	return err
}

// UpdateResourceConfig updates an existing data resource config
func (s *service) UpdateResourceConfig(ctx context.Context, id string, d types.ResourceConfig) error {
	// Name: id is a workaround, as the config API
	// doesn't follow the standar request data paters
	return s.Update(ctx, &Request{Kind: "/dbdef", Parts: []Part{{Name: id}}}, d)
}

// DeleteResourceConfig deletes a single data resource config
func (s *service) DeleteResourceConfig(ctx context.Context, id string) error {
	return s.Delete(ctx, &Request{Kind: "/dbdef", Parts: []Part{{Name: id}}})

}
