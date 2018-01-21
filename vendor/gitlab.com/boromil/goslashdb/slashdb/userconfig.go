package slashdb

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/boromil/goslashdb/types"
)

// UserConfigManager - represents CRUD operations for the UserConfig request configs
type UserConfigManager interface {
	UserConfigs(ctx context.Context) (map[string]types.UserConfig, error)
	UserConfig(ctx context.Context, id string) (types.UserConfig, error)
	CreateUserConfig(ctx context.Context, u types.UserConfig) error
	UpdateUserConfig(ctx context.Context, id string, u types.UserConfig) error
	DeleteUserConfig(ctx context.Context, id string) error
}

// UserConfigs retrives all the user configs
func (s *service) UserConfigs(ctx context.Context) (map[string]types.UserConfig, error) {
	var data = map[string]types.UserConfig{}
	if err := s.Get(ctx, &Request{Kind: "/userdef"}, &data); err != nil {
		return nil, errors.Wrap(err, "error retriving user configs")
	}
	return data, nil
}

// UserConfig retrives a single user config
func (s *service) UserConfig(ctx context.Context, id string) (types.UserConfig, error) {
	var data = types.UserConfig{}
	// Name: id is a workaround, as the config API
	// doesn't follow the standar request data paters
	if err := s.Get(ctx, &Request{Kind: "/userdef", Parts: []Part{{Name: id}}}, &data); err != nil {
		return types.UserConfig{}, errors.Wrap(err, "error retriving user config")
	}
	return data, nil
}

// CreateUserConfig creates a new user config
func (s *service) CreateUserConfig(ctx context.Context, u types.UserConfig) error {
	_, err := s.Create(ctx, &Request{Kind: "/userdef"}, u)
	return err
}

// UpdateUserConfig updates an existing user config
func (s *service) UpdateUserConfig(ctx context.Context, id string, u types.UserConfig) error {
	// Name: id is a workaround, as the config API
	// doesn't follow the standar request data paters
	return s.Update(ctx, &Request{Kind: "/userdef", Parts: []Part{{Name: id}}}, u)
}

// DeleteUserConfig deletes a single user config
func (s *service) DeleteUserConfig(ctx context.Context, id string) error {
	return s.Delete(ctx, &Request{Kind: "/userdef", Parts: []Part{{Name: id}}})
}
