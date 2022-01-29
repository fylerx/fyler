package enum

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type TaskType string

const (
	TaskTypeDocToPdf   TaskType = "doc_to_pdf"
	TaskTypeVideoToMp4 TaskType = "video_to_mp4"
)

var AllTaskType = []TaskType{
	TaskTypeDocToPdf,
	TaskTypeVideoToMp4,
}

func (e TaskType) IsValid() bool {
	switch e {
	case TaskTypeDocToPdf, TaskTypeVideoToMp4:
		return true
	}
	return false
}

func (e TaskType) String() string {
	return string(e)
}

func (e *TaskType) Scan(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TaskType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TaskType", str)
	}
	return nil
}

func (e *TaskType) Value() (driver.Value, error) {
	return strconv.Quote(e.String()), nil
}
