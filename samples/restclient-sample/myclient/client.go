package myclient

import (
	"context"
)

//xgo:generate go tool shoot new -getset -json -file=$GOFILE
//go:generate go tool shoot enum -json -file=$GOFILE
//xgo:generate go tool shoot rest -type=Client

type KV struct {
	key   string
	value string
}

type Gender int32

const (
	Unknown Gender = iota
	Male
	Femal
)

type User struct {
	id     string
	name   string
	age    int
	gender Gender
}

type Client interface {
	//shoot: Get("/get")
	Get(ctx context.Context, key string) (*KV, error)

	//shoot: Post("/set")
	Set(ctx context.Context, kv *KV) error

	//shoot: Get("/users/{id}")
	//shoot: alias={userID:id},{q1:a1}
	GetUser(ctx context.Context, userID string, q1 int, q2 string) (*User, error)
}
