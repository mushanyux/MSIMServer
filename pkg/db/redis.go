package db

import "github.com/mushanyux/MSIMServer/pkg/redis"

func NewRedis(addr string, password string) *redis.Conn {
	return redis.New(addr, password)
}
