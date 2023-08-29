package integration

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type simpleStruct struct {
	Id      uint64
	Bool    bool
	BoolPtr *bool

	Int    int
	IntPtr *int

	Int8    int8
	Int8Ptr *int8

	Int16    int16
	Int16Ptr *int16

	Int32    int32
	Int32Ptr *int32

	Int64    int64
	Int64Ptr *int64

	Uint    uint
	UintPtr *uint

	Uint8    uint8
	Uint8Ptr *uint8

	Uint16    uint16
	Uint16Ptr *uint16

	Uint32    uint32
	Uint32Ptr *uint32

	Uint64    uint64
	Uint64Ptr *uint64

	Float32    float32
	Float32Ptr *float32

	Float64    float64
	Float64Ptr *float64

	Byte      byte
	BytePtr   *byte
	ByteArray []byte

	String string

	NullStringPtr  *sql.NullString
	NullInt16Ptr   *sql.NullInt16
	NullInt32Ptr   *sql.NullInt32
	NullInt64Ptr   *sql.NullInt64
	NullBoolPtr    *sql.NullBool
	NullFloat64Ptr *sql.NullFloat64
	JsonColumn     *simpleJson
}

type simpleUser struct {
	Name string
}

type simpleJson struct {
	Val   simpleUser
	Valid bool
}

func (j *simpleJson) Scan(src any) error {
	if src == nil {
		return nil
	}

	var bs []byte

	switch val := src.(type) {
	case string:
		bs = []byte(val)
	case []byte:
		bs = val
	case *[]byte:
		bs = *val
	default:
		return errors.New("invalid type")
	}

	if len(bs) == 0 {
		return nil
	}

	err := json.Unmarshal(bs, &j.Val)
	if err != nil {
		return err
	}

	j.Valid = true
	return nil
}

func (j *simpleJson) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}

	if !j.Valid {
		return nil, nil
	}

	bs, err := json.Marshal(j.Val)
	if err != nil {
		return nil, err
	}

	return bs, nil
}
