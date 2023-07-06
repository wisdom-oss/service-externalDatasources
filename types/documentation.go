package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Documentation represents a single documentation entry for a data source
type Documentation struct {
	Type      string        `json:"type"`
	Location  string        `json:"location"`
	Verbosity NoneHighRange `json:"verbosity"`
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
	regex := regexp.MustCompile(`^\{"\('((?:[^"']|\\'|\\")+)','((?:[^"']|\\'|\\")+)',(none|low|medium|high)\)"\}$`)
	values := regex.FindStringSubmatch(rowString)[1:]
	if len(values) < 3 {
		return errors.New("not enough items in object")
	}
	doc.Type = values[0]
	doc.Location = values[1]
	doc.Verbosity = NoneHighRange(values[2])
	*d = doc
	return nil
}

func (d *Documentation) Value() (driver.Value, error) {
	escapedType := strings.ReplaceAll(d.Type, `'`, `''`)
	escapedLocation := strings.ReplaceAll(d.Type, `'`, `''`)
	return fmt.Sprintf("('%s','%s',%s::nonehighrange)::documentation", escapedType, escapedLocation, d.Verbosity), nil
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
	fmt.Println(rowString)
	entries := strings.Split(rowString, "\",\"")
	var documentationElements Documentations
	for idx, entry := range entries {
		var doc Documentation
		var regex *regexp.Regexp
		switch {
		// there is only one documentation entry
		case idx == 0 && len(entries) == 1:
			regex = regexp.MustCompile(`^\{"\('((?:[^"']|\\'|\\")+)','((?:[^"']|\\'|\\")+)',(none|low|medium|high)\)"\}$`)
			break
		// the current parsed entry is the first one of at least two
		case idx == 0 && len(entries) > 1:
			regex = regexp.MustCompile(`^\{"\(((?:[^"']|\\'|\\")+),((?:[^"']|\\'|\\")+),(none|low|medium|high)\)$`)
			break
		// the currently parsed entry is the last one of all entries
		case (idx + 1) == len(entries):
			regex = regexp.MustCompile(`^\(((?:[^"']|'|\\")+),((?:[^"']|\\'|\\")+),(none|low|medium|high)\)"\}$`)
			break
		default:
			regex = regexp.MustCompile(`^\(((?:[^"']|\\'|\\")+),((?:[^"']|\\'|\\")+),(none|low|medium|high)\)$`)
		}
		submatches := regex.FindStringSubmatch(entry)
		values := submatches[1:]
		if len(values) < 3 {
			return errors.New("not enough items in object")
		}
		doc.Type = strings.Trim(values[0], `\"`)
		doc.Location = values[1]
		doc.Verbosity = NoneHighRange(values[2])
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
	documentationString = fmt.Sprintf("ARRAY[%s]", documentationString)
	return documentationString, nil
}
