package remote

import (
	"fmt"
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
)

var rc redis.Conn

func InitRedis() (err error, cancel func() error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	rc, err = redis.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		return err, nil
	}
	return nil, rc.Close
}

type Redis interface {
	Get(key string) (error, []byte)
	Set(key string, value []byte, expiry int) error
	Delete(key string) error
}

type RedisImp struct {}

var redisObj Redis = &RedisImp{}

func GetRedis() Redis {
	return redisObj
}

func SetRedis(r Redis) {
	redisObj = r
}

func (r *RedisImp) Get(key string) (error, []byte) {
	s, err := redis.Bytes(rc.Do("GET", key))
	if err != nil {
		log.Printf("[WARN] redis get error: %v", err)
		return err, nil
	}
	return nil, s
}

func (r *RedisImp) Set(key string, value []byte, expiry int) error {
	_, err := rc.Do("SET", key, value, "EX", expiry)
	if err != nil {
		log.Printf("[ERROR] redis set error: %v", err)
	}
	return err
}

func (r *RedisImp) Delete(key string) error {
	_, err := rc.Do("DEL", key)
	if err != nil {
		log.Printf("[ERROR] redis delete error: %v", err)
	}
	return err
}