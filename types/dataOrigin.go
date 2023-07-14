package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type DataOrigin struct {
	Provider string `json:"provider"`
	Creator  string `json:"creator"`
	Owner    string `json:"owner"`
}

func (do *DataOrigin) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	regex := regexp.MustCompile(`^\("((?:[^"']|''|"")+)","*((?:[^"']|''|"")+)"*,"*((?:[^"']|''|"")+)"*\)$`)
	matches := regex.FindStringSubmatch(rowString)
	// now clean up the strings
	provider := strings.ReplaceAll(matches[1], `""`, `"`)
	provider = strings.ReplaceAll(provider, `''`, `'`)
	creator := strings.ReplaceAll(matches[2], `""`, `"`)
	creator = strings.ReplaceAll(creator, `''`, `'`)
	owner := strings.ReplaceAll(matches[3], `""`, `"`)
	owner = strings.ReplaceAll(owner, `''`, `'`)
	var obj DataOrigin
	obj.Provider = provider
	obj.Creator = creator
	obj.Owner = owner
	*do = obj
	return nil
}

func (do DataOrigin) Value() (driver.Value, error) {
	provider := strings.ReplaceAll(do.Provider, `"`, `""`)
	provider = strings.ReplaceAll(provider, `'`, `''`)
	creator := strings.ReplaceAll(do.Creator, `"`, `""`)
	creator = strings.ReplaceAll(creator, `'`, `''`)
	owner := strings.ReplaceAll(do.Owner, `"`, `""`)
	owner = strings.ReplaceAll(owner, `'`, `''`)
	return fmt.Sprintf("(\"%s\",\"%s\",\"%s\")", provider, creator, owner), nil
}
