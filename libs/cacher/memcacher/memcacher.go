package memcacher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	cacher "github.com/unvs/libs/cacher"
)

type MemcacheCacher struct {
	server string
	client *memcache.Client
	prefix string
	expiry time.Duration
}
type Cacher cacher.Cacher

// NewMemcacheCacher creates a new MemcacheCacher
func (c *MemcacheCacher) Init() error {
	mc := memcache.New(c.server)
	err := mc.Ping()
	if err != nil {
		return fmt.Errorf("memcache ping error: %w", err)
	}
	return nil
}

func (c *MemcacheCacher) SetPrefix(prefix string) error {
	c.prefix = prefix
	return nil
}

func (c *MemcacheCacher) GetPrefix() string {
	return c.prefix
}

func (c *MemcacheCacher) GetKey(key string) string {
	return c.prefix + key
}

func (c *MemcacheCacher) Get(key string) (interface{}, error) {
	item, err := c.client.Get(c.GetKey(key))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil // Return nil, nil for cache miss
		}
		return nil, fmt.Errorf("memcache get error: %w", err)
	}
	return item.Value, nil
}

func (c *MemcacheCacher) Set(key string, value interface{}) error {
	b, err := marshalValue(value)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:        c.GetKey(key),
		Value:      b,
		Expiration: int32(c.expiry.Seconds()),
	}
	return c.client.Set(item)
}

func (c *MemcacheCacher) Delete(key string) error {
	return c.client.Delete(c.GetKey(key))
}

func (c *MemcacheCacher) GetDict(key string) (map[string]interface{}, error) {
	val, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil // Cache miss
	}

	var dict map[string]interface{}
	err = json.Unmarshal(val.([]byte), &dict)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dict: %w", err)
	}
	return dict, nil
}

func (c *MemcacheCacher) SetDict(key string, value map[string]interface{}) error {
	return c.Set(key, value) // Use the existing Set method
}

// marshalValue marshals any value to []byte
func marshalValue(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case int, int64, uint64, float64, int32, uint32:
		return []byte(fmt.Sprintf("%v", v)), nil
	default:
		jsonByte, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("can't marshal value to json: %w", err)
		}
		return jsonByte, nil
	}
}
func (c *MemcacheCacher) HealthCheck(timeout time.Duration) error {
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	retry := 1000
	for i := 0; i < retry; i++ {
		if err := c.client.Ping(); err != nil {
			if i == retry-1 {
				return fmt.Errorf("memcached health check failed: %w", err)
			}
			time.Sleep(2000 * time.Millisecond)
			continue
		}
		return nil
	}
	return nil
}
