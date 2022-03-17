package enum

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type Service string

const (
	ServiceS3     Service = "s3"
	ServiceSentry Service = "sentry"
)

var AllService = []Service{
	ServiceS3,
	ServiceSentry,
}

func (e Service) IsValid() bool {
	switch e {
	case ServiceS3, ServiceSentry:
		return true
	}
	return false
}

func (e Service) String() string {
	return string(e)
}

func (e *Service) Scan(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Service(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Service", str)
	}
	return nil
}

func (e *Service) Value() (driver.Value, error) {
	return strconv.Quote(e.String()), nil
}
