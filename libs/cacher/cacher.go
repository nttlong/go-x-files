// this file is declaring the cacher interface and its methods
package cacher

import "time"

type SetExpireOptions struct {
	Expiry time.Duration
}

type SetExpireOption func(*SetExpireOptions)

func WithExpiry(expiry time.Duration) SetExpireOption {
	return func(o *SetExpireOptions) {
		o.Expiry = expiry
	}
}

type Cacher interface {

	//get text from cache
	GetText(key string) string
	// set text in cache with expiry time
	// Params:
	// key: key to store
	// value: value to store
	// expiry: duration after which cache should expire
	// Returns:
	// error: if any error occurs during set operation
	SetText(key string, value string, optiosn ...SetExpireOption)
	// delete key from cache
	// Paramters:
	// key: key to delete
	// Returns:
	Delete(key string)
	GetDict(key string) map[string]interface{}
	SetDict(key string, value map[string]interface{}, options ...SetExpireOption)
	SetStruct(key string, value interface{}, options ...SetExpireOption)
	GetStruct(key string, value interface{})
}
