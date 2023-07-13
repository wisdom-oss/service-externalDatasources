package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
)

type DataProvisioning struct {
	Type   string `json:"type"`
	Format string `json:"format"`
}

func (dp *DataProvisioning) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	regex := regexp.MustCompile(`^\(([^}()[\],]+),([^{}()[\],]+)\)$`)
	matches := regex.FindStringSubmatch(rowString)
	if len(matches) != 3 {
		return errors.New("unexpected count of matches")
	}
	dp.Type = matches[1]
	dp.Format = matches[2]
	return nil
}

func (dp DataProvisioning) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)", dp.Type, dp.Format), nil
}
