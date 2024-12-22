package memcacher

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	cacher "github.com/unvs/libs/cacher"
)

type MemcacheCacher struct {
	Server string
	client *memcache.Client
	Prefix string
	Expiry time.Duration
}
type Cacher cacher.Cacher

func makeHas256Key(prefix string, key string) string {
	h := sha256.Sum256([]byte(prefix + "/" + key))
	return hex.EncodeToString(h[:])
}
func (c *MemcacheCacher) init() {
	if c.client == nil {
		c.client = memcache.New(c.Server)
		c.client.Timeout = 10 * time.Second

	}
}
func (c *MemcacheCacher) SetText(key string, value string, options ...cacher.SetExpireOption) {
	c.init()
	var expire time.Duration
	opts := cacher.SetExpireOptions{
		Expiry: c.Expiry, // Default expiry
	}

	for _, option := range options {
		option(&opts)
	}
	if expire == 0 {
		expire = c.Expiry
	}
	err := c.client.Set(&memcache.Item{
		Key:        makeHas256Key(c.Prefix, key),
		Value:      []byte(value),
		Expiration: int32(expire.Seconds()),
	})
	if err != nil {
		panic(fmt.Sprintf("MemcacheCacher: SetText: %s", err))
	}
}

func (c *MemcacheCacher) GetText(key string) string {
	c.init()
	item, err := c.client.Get(makeHas256Key(c.Prefix, key))
	if err != nil {
		panic(fmt.Sprintf("MemcacheCacher: GetText: %s", err))
	}
	return string(item.Value)
}
func (c *MemcacheCacher) Delete(key string) {
	c.init()
	err := c.client.Delete(makeHas256Key(c.Prefix, key))
	if err != nil {
		panic(fmt.Sprintf("MemcacheCacher: Delete: %s", err))
	}
}
func (c *MemcacheCacher) SetDict(key string, value map[string]interface{}, options ...cacher.SetExpireOption) {
	c.init()
	var expire time.Duration
	opts := cacher.SetExpireOptions{
		Expiry: c.Expiry, // Default expiry
	}

	for _, option := range options {
		option(&opts)
	}
	if expire == 0 {
		expire = c.Expiry
	}
	//convert value to []byte
	jsonValue, err := json.Marshal(value)
	if err != nil {
		panic(fmt.Sprintf("MemcacheCacher: SetDict: %s", err))
	}
	err = c.client.Set(&memcache.Item{
		Key:        makeHas256Key(c.Prefix, key),
		Value:      []byte(jsonValue),
		Expiration: int32(expire.Seconds()),
	})
	if err != nil {
		panic(fmt.Sprintf("MemcacheCacher: SetDict: %s", err))
	}
}
func (c *MemcacheCacher) GetDict(key string) map[string]interface{} {
	c.init()
	item, err := c.client.Get(makeHas256Key(c.Prefix, key))
	if err != nil {
		panic(fmt.Sprintf("MemcacheCacher: GetDict: %s", err))
	}
	//convert vlue of item to map[string]interface{}
	var value map[string]interface{}
	err = json.Unmarshal(item.Value, &value)
	if err != nil {
		panic(fmt.Sprintf("MemcacheCacher: GetDict: %s", err))
	}
	return value
}
func (c *MemcacheCacher) SetStruct(key string, value interface{}, options ...cacher.SetExpireOption) {
	c.init()
	var expire time.Duration
	opts := cacher.SetExpireOptions{
		Expiry: c.Expiry, // Default expiry
	}

	for _, option := range options {
		option(&opts)
	}
	if expire == 0 {
		expire = c.Expiry
	}
	//convert value to []byte
	jsonValue, err := json.Marshal(value)
	if err != nil {
		panic(fmt.Sprintf("MemcacheCacher: SetStruct: %s", err))
	}
	err = c.client.Set(&memcache.Item{
		Key:        makeHas256Key(c.Prefix, key),
		Value:      []byte(jsonValue),
		Expiration: int32(expire.Seconds()),
	})
	if err != nil {
		panic(fmt.Sprintf("MemcacheCacher: SetStruct: %s", err))
	}
}
func (c *MemcacheCacher) GetStruct(key string, value interface{}) {
	c.init()
	item, err := c.client.Get(makeHas256Key(c.Prefix, key))
	if err != nil {
		panic(fmt.Sprintf("MemcacheCacher: GetStruct: %s", err))
	}
	//convert value of item to map[string]interface{}

	// value := new(returnType)
	err = json.Unmarshal(item.Value, &value)
	if err != nil {
		panic(fmt.Sprintf("MemcacheCacher: GetStruct: %s", err))
	}

}
