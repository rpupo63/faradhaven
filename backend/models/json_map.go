package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// JSONMap is a custom type for map[string]interface{} that implements
// sql.Scanner and driver.Valuer interfaces for JSONB column type in GORM.
type JSONMap map[string]interface{}

// Value returns a driver Value.
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	// Marshal the map into a JSON byte slice
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface.
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	// Unmarshal the JSON byte slice into the map
	return json.Unmarshal(bytes, j)
}
