package cache

import "time"

type Cache interface {
	SetCache (key string, value interface{}, time time.Duration)
	GetCache (key string) (interface{}, bool)
}
