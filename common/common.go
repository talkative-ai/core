package common

import (
	"database/sql/driver"
	"log"
	"strings"

	"github.com/go-redis/redis"
)

type BSliceIndex struct {
	Index  int
	Bslice []byte
}

type RedisCommand func(redis *redis.Client)

func RedisSET(key string, bytes []byte) RedisCommand {
	return func(redis *redis.Client) {
		result := redis.Set(key, bytes, 0)
		if err := result.Err(); err != nil {
			log.SetFlags(log.Llongfile)
			log.Println("Redis command error", err.Error())
		}
	}
}

func RedisHSET(key, field string, bytes []byte) RedisCommand {
	return func(redis *redis.Client) {
		result := redis.HSet(key, field, bytes)
		if err := result.Err(); err != nil {
			log.SetFlags(log.Llongfile)
			log.Println("Redis command error", err.Error())
		}
	}
}

func RedisSADD(key string, members ...interface{}) RedisCommand {
	return func(redis *redis.Client) {
		result := redis.SAdd(key, members...)
		if err := result.Err(); err != nil {
			log.SetFlags(log.Llongfile)
			log.Println("Redis command error", err.Error())
		}
	}
}

type StringArray struct {
	Val []string
}

func (arr *StringArray) Value() (driver.Value, error) {
	return arr.Val, nil
}

func (arr *StringArray) Scan(src interface{}) error {
	str := string(src.([]byte))
	str = str[1 : len(str)-1]
	arr.Val = strings.Split(str, ",")
	return nil
}
