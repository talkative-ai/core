package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/talkative-ai/go.uuid"
)

type BSliceIndex struct {
	Index  int
	Bslice []byte
}

type RedisCommand struct {
	Exec  func(cmd RedisCommand, redis *redis.Client)
	Key   string
	Value interface{}
}

func RedisSET(key string, bytes []byte) RedisCommand {
	fn := func(cmd RedisCommand, redis *redis.Client) {
		result := redis.Set(cmd.Key, cmd.Value, 0)
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
	fn := func(cmd RedisCommand, redis *redis.Client) {
		result := redis.HSet(cmd.Key, field, cmd.Value)
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
	fn := func(cmd RedisCommand, redis *redis.Client) {
		result := redis.SAdd(cmd.Key, members...)
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

// TODO: Clean all this mess up. StringArray, StringArray2D, StringArray2DJSON
// all of this seems really heavy and awkward

func (arr *StringArray) UnmarshalJSON(b []byte) (err error) {
	return json.Unmarshal(b, &arr.Val)
}

type StringArray2D []StringArray

func (a *StringArray2D) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &a)
}

func (arr *StringArray2D) Value() (driver.Value, error) {
	return json.Marshal(&arr)
}

type StringArray2DJSON [][]string

func (arr *StringArray2DJSON) UnmarshalJSON(b []byte) (err error) {
	val := [][]string{}
	err = json.Unmarshal(b, &val)
	if err != nil {
		return
	}
	*arr = StringArray2DJSON(val)
	return nil
}
func (arr *StringArray2DJSON) Scan(b interface{}) (err error) {
	val := [][]string{}
	bytes, ok := b.([]byte)
	if !ok {
		return fmt.Errorf("Impossible typecast in StringArray2DJSON with value: %v+", b)
	}
	err = json.Unmarshal(bytes, &val)
	if err != nil {
		return
	}
	*arr = StringArray2DJSON(val)
	return nil
}
func (arr *StringArray2DJSON) Value() (driver.Value, error) {
	return json.Marshal(&arr)
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
