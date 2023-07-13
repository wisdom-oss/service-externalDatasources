package types

import (
	"database/sql/driver"
	"errors"
	"external-api-service/enums"
	"fmt"
	"regexp"
	"strings"
)

// Documentation represents a single documentation entry for a data source
type Documentation struct {
	Type      string              `json:"type"`
	Location  string              `json:"location"`
	Verbosity enums.NoneHighRange `json:"verbosity"`
}

func (d *Documentation) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	var doc Documentation
	regex := regexp.MustCompile(`^\{"\(([^[:cntrl:],]+),([^[:cntrl:],]+),(none|low|medium|high)\)\}$`)
	values := regex.FindStringSubmatch(rowString)[1:]
	if len(values) < 3 {
		return errors.New("not enough items in object")
	}
	doc.Type = values[0]
	doc.Location = values[1]
	doc.Verbosity = enums.NoneHighRange(values[2])
	*d = doc
	return nil
}

func (d *Documentation) Value() (driver.Value, error) {
	escapedType := strings.ReplaceAll(d.Type, `'`, `''`)
	escapedLocation := strings.ReplaceAll(d.Type, `'`, `''`)
	return fmt.Sprintf("\"(%s,%s,%s)\"", escapedType, escapedLocation, d.Verbosity), nil
}

type Documentations []Documentation

func (d *Documentations) Scan(src interface{}) error {
	var rowString string
	switch src.(type) {
	case []byte:
		rowString = string(src.([]byte))
	case string:
		rowString = src.(string)
	default:
		return errors.New("unsupported scan input")
	}
	entries := strings.Split(rowString, "\",\"")
	var documentationElements Documentations
	for _, entry := range entries {
		var doc Documentation
		regex := regexp.MustCompile(`^{*"*\(([^[:cntrl:],]+),([^[:cntrl:],]+),(none|low|medium|high)\)"*}*$`)
		submatches := regex.FindStringSubmatch(entry)
		values := submatches[1:]
		if len(values) < 3 {
			return errors.New("not enough items in object")
		}
		doc.Type = strings.Trim(values[0], ``)
		doc.Location = values[1]
		doc.Verbosity = enums.NoneHighRange(values[2])
		documentationElements = append(documentationElements, doc)
	}

	*d = documentationElements
	return nil
}

func (d Documentations) Value() (driver.Value, error) {
	var singleParts []string
	for _, documentation := range d {
		val, err := documentation.Value()
		if err != nil {
			return nil, err
		}
		singleParts = append(singleParts, val.(string))
	}
	documentationString := strings.Join(singleParts, ",")
	documentationString = fmt.Sprintf("{%s}", documentationString)
	return documentationString, nil
}
