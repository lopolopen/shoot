package myclient

import (
	"context"

	"github.com/lopolopen/shoot"
)

//go:generate go tool shoot new -getset -json -file=$GOFILE
//go:generate go tool shoot rest -type=Client

type KV struct {
	key   string
	value string
}

type Book struct {
	name  string
	price int
}

type Client interface {
	shoot.RestClient

	//shoot: Get(/get)
	Get(ctx context.Context, key string) (*KV, error)

	//shoot: Post(/set)
	Set(ctx context.Context, kv *KV) error

	//shoot: Get(/users/{id})
	//shoot: alias={userID:id},{q1:a1}
	GetUser(ctx context.Context, userID string, q1 int, q2 string) (*Book, error)
}
