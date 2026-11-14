package otp

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func NewStore(addr string) *Store {
	return &Store{client: redis.NewClient(&redis.Options{Addr: addr})}
}

func NewStoreWithClient(client *redis.Client) *Store {
	return &Store{client: client}
}

func (s *Store) Save(ctx context.Context, userID, codeHash string) error {
	return s.client.Set(ctx, key(userID), codeHash, 0).Err()
}

func (s *Store) Find(ctx context.Context, userID string) (string, bool, error) {
	codeHash, err := s.client.Get(ctx, key(userID)).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return codeHash, true, nil
}

func (s *Store) Delete(ctx context.Context, userID string) (bool, error) {
	deleted, err := s.client.Del(ctx, key(userID)).Result()
	return deleted > 0, err
}

func key(userID string) string {
	return "xotp:otp:" + userID
}
