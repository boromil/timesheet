package slashdb

import (
	"context"
	"fmt"

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
func (s *Service) UserConfigs(ctx context.Context) (map[string]types.UserConfig, error) {
	var data = map[string]types.UserConfig{}
	if err := s.Get(ctx, &Request{Kind: "/userdef"}, &data); err != nil {
		return nil, fmt.Errorf("error retriving user configs: %w", err)
	}
	return data, nil
}

// UserConfig retrives a single user config
func (s *Service) UserConfig(ctx context.Context, id string) (types.UserConfig, error) {
	var data = types.UserConfig{}
	// Name: id is a workaround, as the config API
	// doesn't follow the standar request data paters
	if err := s.Get(ctx, &Request{Kind: "/userdef", Parts: []Part{{Name: id}}}, &data); err != nil {
		return types.UserConfig{}, fmt.Errorf("error retriving user config: %w", err)
	}
	return data, nil
}

// CreateUserConfig creates a new user config
func (s *Service) CreateUserConfig(ctx context.Context, u types.UserConfig) error {
	_, err := s.Create(ctx, &Request{Kind: "/userdef"}, u)
	return err
}

// UpdateUserConfig updates an existing user config
func (s *Service) UpdateUserConfig(ctx context.Context, id string, u types.UserConfig) error {
	// Name: id is a workaround, as the config API
	// doesn't follow the standar request data paters
	return s.Update(ctx, &Request{Kind: "/userdef", Parts: []Part{{Name: id}}}, u)
}

// DeleteUserConfig deletes a single user config
func (s *Service) DeleteUserConfig(ctx context.Context, id string) error {
	return s.Delete(ctx, &Request{Kind: "/userdef", Parts: []Part{{Name: id}}})
}
