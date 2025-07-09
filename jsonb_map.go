package nullutil

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONBMap map[string]interface{}

func (j *JSONBMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONBMap", value)
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONBMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}
