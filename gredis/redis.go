package gredis

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisSetting struct {
	Host string
	Password string
	MaxIdle int
	MaxActive int
	IdleTimeout time.Duration
}

var RedisConn *redis.Pool

func Setup() error {

	redisSetting := RedisSetting{
		Host: "127.0.0.1:6379",
		Password: "",
		MaxIdle: 30,
		MaxActive: 30,
		IdleTimeout: 300 * time.Second,
	}

	RedisConn = &redis.Pool {
		MaxIdle:    redisSetting.MaxIdle,
		MaxActive:  redisSetting.MaxActive,
		IdleTimeout:redisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisSetting.Host)
			if err != nil {
				return nil, err
			}
			if redisSetting.Password != "" {
				if _, err := c.Do("AUTH", redisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		}, 
	}
	return nil
}

func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}
	return nil
}

func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func LikeGets(key string) ([][]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return nil, err
	}
	var result [][]byte
	for _, key := range keys {
		reply, err := Get(key)
		if err != nil {
			return nil, err
		}
		result = append(result, reply)
	}
	return result, nil
}

func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}
	for _, key := range keys {
		_, err := Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}