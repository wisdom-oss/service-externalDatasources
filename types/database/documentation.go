package databaseTypes

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
	doc.Type = strings.ReplaceAll(values[0], `''`, `'`)
	doc.Type = strings.ReplaceAll(doc.Type, `""`, `"`)
	doc.Type = strings.ReplaceAll(doc.Type, `\"`, ``)
	doc.Location = strings.ReplaceAll(values[1], `''`, `'`)
	doc.Location = strings.ReplaceAll(doc.Location, `""`, `"`)
	doc.Location = strings.ReplaceAll(doc.Location, `\"`, ``)
	doc.Verbosity = enums.NoneHighRange(values[2])
	*d = doc
	return nil
}

func (d *Documentation) Value() (driver.Value, error) {
	escapedType := strings.ReplaceAll(d.Type, `'`, `''`)
	escapedType = strings.ReplaceAll(escapedType, `"`, `""`)
	escapedLocation := strings.ReplaceAll(d.Location, `'`, `''`)
	escapedLocation = strings.ReplaceAll(escapedLocation, `"`, `""`)
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
		entry = strings.Trim(entry, `\{}/'"`)
		regex := regexp.MustCompile(`^\(([^[:cntrl:],]+),([^[:cntrl:],]+),(none|low|medium|high)\)$`)
		submatches := regex.FindStringSubmatch(entry)
		values := submatches[1:]
		if len(values) < 3 {
			return errors.New("not enough items in object")
		}
		doc.Type = strings.ReplaceAll(values[0], `''`, `'`)
		doc.Type = strings.ReplaceAll(doc.Type, `""`, `"`)
		doc.Type = strings.ReplaceAll(doc.Type, `\"`, ``)
		doc.Location = strings.ReplaceAll(values[1], `''`, `'`)
		doc.Location = strings.ReplaceAll(doc.Location, `""`, `"`)
		doc.Location = strings.ReplaceAll(doc.Location, `\"`, ``)
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
