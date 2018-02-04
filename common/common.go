package common

import (
	"database/sql/driver"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	uuid "github.com/artificial-universe-maker/go.uuid"
	"github.com/go-redis/redis"
)

type BSliceIndex struct {
	Index  int
	Bslice []byte
}

type RedisCommand struct {
	Exec  func(redis *redis.Client)
	Key   string
	Value interface{}
}

func RedisSET(key string, bytes []byte) RedisCommand {
	fn := func(redis *redis.Client) {
		result := redis.Set(key, bytes, 0)
		if err := result.Err(); err != nil {
			log.SetFlags(log.Llongfile)
			log.Println("Redis command error", err.Error())
		}
	}
	return RedisCommand{
		Exec:  fn,
		Key:   key,
		Value: bytes,
	}
}

func RedisHSET(key, field string, bytes []byte) RedisCommand {
	fn := func(redis *redis.Client) {
		result := redis.HSet(key, field, bytes)
		if err := result.Err(); err != nil {
			log.SetFlags(log.Llongfile)
			log.Println("Redis command error", err.Error())
		}
	}
	return RedisCommand{
		Exec:  fn,
		Key:   key,
		Value: bytes,
	}
}

func RedisSADD(key string, members ...interface{}) RedisCommand {
	fn := func(redis *redis.Client) {
		result := redis.SAdd(key, members...)
		if err := result.Err(); err != nil {
			log.SetFlags(log.Llongfile)
			log.Println("Redis command error", err.Error())
		}
	}
	return RedisCommand{
		Exec:  fn,
		Key:   key,
		Value: nil,
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
	if len(str) <= 2 {
		arr.Val = make([]string, 0)
		return nil
	}
	str = str[1 : len(str)-1]
	arr.Val = strings.Split(str, ",")
	for i, v := range arr.Val {
		if len(str) == 0 {
			continue
		}
		if v[0] == '"' {
			arr.Val[i] = v[1 : len(v)-1]
			// Trim wrapping quotes
		}
	}
	return nil
}

type StringArray2D []StringArray

func (arr *StringArray2D) Value() (driver.Value, error) {
	var val []string
	for _, stringArray := range *arr {
		var stringCollection []string
		for _, str := range stringArray.Val {
			stringCollection = append(val, fmt.Sprintf(`"%v"`, str))
		}
		val = append(val, fmt.Sprintf("{%v}", strings.Join(stringCollection, `,`)))
	}
	return fmt.Sprintf("{%v}", strings.Join(val, `,`)), nil
}

func PseudoRand(max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

func ChooseString(list []string) string {
	l := list
	if len(l) == 1 {
		return l[0]
	}
	return l[PseudoRand(len(l))]
}

type SyncMapUUID struct {
	Value map[uuid.UUID]bool
	Mutex sync.Mutex
}
