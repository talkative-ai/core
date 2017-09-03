package common

import (
	"database/sql"
	"database/sql/driver"
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
		redis.Set(key, bytes, 0)
	}
}

func RedisHSET(key, field string, bytes []byte) RedisCommand {
	return func(redis *redis.Client) {
		redis.HSet(key, field, bytes)
	}
}

func RedisSADD(key string, members ...interface{}) RedisCommand {
	return func(redis *redis.Client) {
		redis.SAdd(key, members...)
	}
}

type ProjectItem struct {
	ProjectID            uint64
	Title                string
	ZoneID               uint64
	ActorID              uint64
	DialogID             uint64
	DialogEntry          StringArray
	ParentDialogID       uint64
	ChildDialogID        uint64
	LogicalSetAlways     string
	LogicalSetStatements sql.NullString
	LogicalSetID         uint64
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
