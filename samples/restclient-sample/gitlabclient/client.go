package gitlabclient

import "context"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//go:generate go tool shoot rest -type=Client

type Client interface {
	// Get(/users/{userid})
	GetUser(ctx context.Context, useid string) (*User, error)
}

type restClient struct {
	baseURL string
}

func New() Client {
	return &restClient{
		baseURL: "https://gitlab.com/api/v4",
	}
}

func (c *restClient) GetUser(ctx context.Context, userid string) (*User, error) {
	// implement the logic to call GitLab API and return a User
	return &User{
		ID:   userid,
		Name: "Mock User",
	}, nil
}
