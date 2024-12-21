package memcacher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"crypto/sha256"
	"encoding/hex"

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
	relKey := c.prefix + key
	hash256 := sha256.Sum256([]byte(relKey))
	return hex.EncodeToString(hash256[:])
}

func (c *MemcacheCacher) GetString(key string) string {
	item, err := c.client.Get(c.GetKey(key))
	if err != nil {
		panic(err)
	}
	return string(item.Value)
}

func (c *MemcacheCacher) SetString(key string, value string, expiry time.Duration) error {
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

func (c *MemcacheCacher) GetDict(key string) map[string]interface{} {
	val, err := c.client.Get(key)
	if err != nil {
		panic(err)
	}
	if val == nil {
		return nil
	}

	var dict map[string]interface{}
	err = json.Unmarshal(val.Value, &dict)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal dict: %w", err))
	}
	return dict
}

func (m *MemcacheCacher) SetDict(key string, value map[string]interface{}, expiry time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	finalExpiry := m.expiry // Use the default expiry

	if expiry != 0 { // Check if expiry is not the zero value
		finalExpiry = expiry // Override with the provided expiry
	}

	item := &memcache.Item{Key: m.GetKey(key), Value: data, Expiration: int32(finalExpiry.Seconds())}
	return m.client.Set(item)
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
func NewMemcacheCacher(server, prefix string) *MemcacheCacher {
	if server == "" {
		panic(fmt.Errorf("server address cannot be empty"))
	}

	mc := &MemcacheCacher{
		server: server,
		prefix: prefix,
	}

	client := memcache.New(server)

	// Attempt a Ping to check the connection
	err := client.Ping()
	if err != nil {
		panic(fmt.Errorf("failed to connect to memcached server: %w", err))
	}

	mc.client = client
	return mc
}
