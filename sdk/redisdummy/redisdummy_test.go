package redisdummy

import (
	"testing"
	"github.com/redis/go-redis/v9"
)


func TestRedisClientPrint(t *testing.T) {
	t.Logf("Redis version: %s\n", redis.Version())
}