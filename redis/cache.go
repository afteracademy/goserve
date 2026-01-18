package redis

import (
	"context"
	"encoding/json"
	"time"
)

type Cache[T any] interface {
	SetJSON(key string, value *T, expiration time.Duration) error
	GetJSON(key string) (*T, error)
	SetJSONList(key string, values []*T, expiration time.Duration) error
	GetJSONList(key string) ([]*T, error)
}

type cache[T any] struct {
	context context.Context
	store   Store
}

func NewCache[T any](store Store) Cache[T] {
	return &cache[T]{
		context: context.Background(),
		store:   store,
	}
}

func (c *cache[T]) SetJSON(key string, value *T, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.store.GetInstance().Set(c.context, key, data, expiration).Err()
}

func (c *cache[T]) GetJSON(key string) (*T, error) {
	data, err := c.store.GetInstance().Get(c.context, key).Bytes()
	if err != nil {
		return nil, err
	}

	var dest T
	err = json.Unmarshal(data, &dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func (c *cache[T]) SetJSONList(key string, values []*T, expiration time.Duration) error {
	var list []json.RawMessage
	for _, value := range values {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		list = append(list, data)
	}

	str, err := json.Marshal(list)
	if err != nil {
		return err
	}

	return c.store.GetInstance().Set(c.context, key, str, expiration).Err()
}

func (c *cache[T]) GetJSONList(key string) ([]*T, error) {
	str, err := c.store.GetInstance().Get(c.context, key).Result()
	if err != nil {
		return nil, err
	}

	var list []json.RawMessage
	if err := json.Unmarshal([]byte(str), &list); err != nil {
		return nil, err
	}

	dest := make([]*T, len(list))
	for i, data := range list {
		if err := json.Unmarshal(data, &dest[i]); err != nil {
			return nil, err
		}
	}

	return dest, nil
}
