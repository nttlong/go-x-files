// this file is declaring the cacher interface and its methods
package cacher

import "time"

type Cacher interface {
	Init() error
	HealthCheck(timeout time.Duration) error
	SetPrefix(prefix string) error
	GetPrefix() string
	GetKey(key string) string
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Delete(key string) error
	GetDict(key string) (map[string]interface{}, error)
	SetDict(key string, value map[string]interface{}) error
}