package repository

import "context"

type DataBase struct {
	KeyResponse StorageKeyResponse
}

type StorageKeyResponse interface {
	Search(ctx context.Context, key string) (response string, err error)
	Add(ctx context.Context, key, response string) (err error)
	Delete(ctx context.Context, key string) (err error)
}
