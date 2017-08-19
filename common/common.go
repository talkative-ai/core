package common

import (
	"database/sql"
	"database/sql/driver"
	"strings"
)

type BSliceIndex struct {
	Index  int
	Bslice []byte
}

// RedisBytes is used to communicate values to be written to Redis
type RedisBytes struct {
	Key   string
	Bytes []byte
}

type ProjectItem struct {
	ProjectID            uint64
	ZoneID               uint64
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
