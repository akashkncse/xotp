package otp

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
	mu     sync.RWMutex
	ttl    time.Duration
}

func NewStore(addr string) *Store {
	return newStore(redis.NewClient(&redis.Options{Addr: addr}))
}

func NewStoreWithClient(client *redis.Client) *Store {
	return newStore(client)
}

func newStore(client *redis.Client) *Store {
	return &Store{client: client, ttl: 10 * time.Minute}
}

func (s *Store) Save(ctx context.Context, userID, codeHash string) error {
	s.mu.RLock()
	ttl := s.ttl
	s.mu.RUnlock()
	return s.client.Set(ctx, key(userID), codeHash, ttl).Err()
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

func (s *Store) SetDefaultTTL(ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ttl = ttl
}

func key(userID string) string {
	return "xotp:otp:" + userID
}
