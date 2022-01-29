package enum

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type Status string

const (
	StatusSuccess  Status = "success"
	StatusQueued   Status = "queued"
	StatusProgress Status = "progress"
	StatusFailed   Status = "failed"
)

var AllStatus = []Status{
	StatusSuccess,
	StatusQueued,
	StatusProgress,
	StatusFailed,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusSuccess, StatusQueued, StatusProgress, StatusFailed:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) Scan(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e *Status) Value() (driver.Value, error) {
	return strconv.Quote(e.String()), nil
}
