// this file is declaring the cacher interface and its methods
package cacher

import "time"

type Cacher interface {
	Init() error
	HealthCheck(timeout time.Duration) error
	SetPrefix(prefix string) error
	GetPrefix() string
	GetKey(key string) string
	GetString(key string) string
	SetString(key string, value string, expiry time.Duration) error
	Delete(key string) error
	GetDict(key string) map[string]interface{}
	SetDict(key string, value map[string]interface{}, expiry time.Duration) error
}
