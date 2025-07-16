package limiter

import "time"

type Storage interface {
	Increment(key string, expire time.Duration) (int64, error)
	Expire(key string, expire time.Duration) error
}
