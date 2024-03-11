package types

import (
	"encoding/json"
	"errors"
)

type Tuple struct {
	left  string `db:"x"`
	right string `db:"y"`
}

func (t Tuple) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string{t.left, t.right})
}

func (t *Tuple) UnmarshalJSON(src []byte) error {
	var values []string
	err := json.Unmarshal(src, &values)
	if err != nil {
		return err
	}
	if len(values) != 2 {
		return errors.New("invalid tuple format")
	}
	t.left = values[0]
	t.right = values[1]
	return nil
}
