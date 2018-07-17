package bitly

import (
	"context"
)

type Email struct {
	Email string `json:"email"`
	IsPrimary bool `json:"is_primary"`
	IsVerified bool `json:"is_verified"`
}

type User struct {
	Name    string    `json:"name"`
	Created JSONDate `json:"created"`
	Modified JSONDate `json:"modified"`
	Login string `json:"login"`
	IsActive bool `json:"is_active"`
	Is2FAEnabled bool `json:"is_2fa_enabled"`
	Emails []Email `json:"emails"`
	IsSSOUser bool `json:"is_sso_user"`
}

type UserUpdateOptions struct {
	Name string `json:"name"`
}

type UserClient struct {
	client *Client
}

func (s *UserClient) Get(ctx context.Context) (*User, error) {
	path := versioned("user")
	u := &User{}

	_, err := s.client.get(path, u)
	return u, err
}

func (s *UserClient) Update(ctx context.Context, options *UserUpdateOptions) (*User, error) {
	if options == nil {
		return nil, errOptionsRequired
	}
	path := versioned("user")
	u := &User{}

	_, err := s.client.patch(path, options, u)
	return u, err
}

func (s *UserClient) GetGroups(ctx context.Context, login string) (*groupsResponse, error) {
	if login == "" {
		return nil, &errorParameter{paramName: "login"}
	}
	path := versioned("user")
	groupsResp := &groupsResponse{}

	_, err := s.client.get(path, groupsResp)
	return groupsResp, err
}

type UserService interface {
	Get(ctx context.Context) (*User, error)
	Update(ctx context.Context) (*User, error)
	GetGroups(ctx context.Context, login string) (*groupsResponse, error)
}