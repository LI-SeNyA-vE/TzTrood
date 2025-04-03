package redisdb

import (
	"TzTrood/internal/server/repository"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func NewRedis(conn *redis.Options) *repository.DataBase {
	return &repository.DataBase{
		KeyResponse: newKeyResponse(conn),
	}
}

func newKeyResponse(conn *redis.Options) keyResponse {
	return keyResponse{
		client: redis.NewClient(conn),
	}
}

type keyResponse struct {
	client *redis.Client
}

func (r keyResponse) Search(ctx context.Context, key string) (response string, err error) {
	redisKey := fmt.Sprintf("intent:%s:response", key)
	result, err := r.client.Get(ctx, redisKey).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r keyResponse) Add(ctx context.Context, key, response string) error {
	redisKey := fmt.Sprintf("intent:%s:response", key)
	return r.client.Set(ctx, redisKey, response, 0).Err()
}

func (r keyResponse) Delete(ctx context.Context, key string) error {
	redisKey := fmt.Sprintf("intent:%s:response", key)
	return r.client.Del(ctx, redisKey).Err()
}
