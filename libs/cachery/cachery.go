// this package is used memcached as a cache for go-x-files
package cachery

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	mc         *memcache.Client
	prefix_key string
	once       sync.Once
)

func CreateKey(key string) string {
	h := sha256.Sum256([]byte(prefix_key + "/" + key))
	return hex.EncodeToString(h[:])
}

func Init(servers string, prefix string) {
	once.Do(func() {
		// Set prefix key
		prefix_key = prefix
		// Allocate the pointer to mc outside the Do function
		mc = new(memcache.Client)
		// Initialize mc with dereferenced pointer
		mc = memcache.New(servers)
	})
}

// check if memcached is alive if not try to reconnect after 1 second
func HealthCheck() {
	if mc == nil {
		panic("memcached client is not initialized, please call Init() of cachery package first")
	}

	ok := false
	for !ok {
		err := mc.Ping()
		if err != nil {
			// Try to reconnect after 1 second
			fmt.Println("Error pinging memcached:", err) // Informational log
			time.Sleep(1 * time.Second)
		} else {
			ok = true
			fmt.Println("Memcached connection established.") // Informational log
		}
	}
}
func Set[T any](key string, value T, expiration time.Duration) {
	if mc == nil {
		panic(fmt.Errorf("memcached client is not initialized, please call Init() first"))
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	if err != nil {
		panic(fmt.Errorf("gob encode failed: %w", err))
	}

	item := &memcache.Item{
		Key:        CreateKey(key),
		Value:      buf.Bytes(),
		Expiration: int32(expiration.Seconds()),
	}

	mc.Set(item)
}

func Get[T any](key string, out *T) bool {
	if mc == nil {
		panic(fmt.Errorf("memcached client is not initialized"))
	}

	item, err := mc.Get(CreateKey(key))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return false
		}
		panic(fmt.Errorf("memcached Get failed: %w", err))
	}

	buf := bytes.NewReader(item.Value)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(out)
	if err != nil {
		panic(fmt.Errorf("gob decode failed: %w", err))
	}
	return true
}

func Delete(key string) error {
	if mc == nil {
		panic("memcached client is not initialized, please call Init() of cachery package first")
	}
	return mc.Delete(CreateKey(key))
}
