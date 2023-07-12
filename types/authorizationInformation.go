package types

import (
	"database/sql/driver"
	"fmt"
)

type AuthorizationDataLocation string

const (
	LOCATION_HEADER AuthorizationDataLocation = "header"
	LOCATION_QUERY  AuthorizationDataLocation = "query"
)

func (loc AuthorizationDataLocation) ToString() string {
	return fmt.Sprintf("%s", loc)
}

type AuthorizationInformation struct {
	Location AuthorizationDataLocation
	Key      string
	Token    string
}

// Value implements the driver.Value interface to allow writing this struct
// into the database using the default pq driver
func (a AuthorizationInformation) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s,%s)", a.Location.ToString(), a.Key, a.Token), nil
}

func (a *AuthorizationInformation) Scan(src interface{}) error {
	rowString := string(src.([]byte))
	_, err := fmt.Sscanf(rowString, "(%s,%s,%s)", &a.Location, &a.Key, &a.Token)
	return err
}
