package dbtypes

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"encoding/json"
	"fmt"
)

// JSON implements the database/sql/driver Scanner and Valuer interfaces,
// as well as gob.GobEncoder and gob.GobDecoder interfaces.
type JSON map[string]interface{}

func init() {
	// Register JSON type for gob encoding/decoding
	gob.Register(JSON{})
	gob.Register(&Date{})
}

// Scan scans a value into JSON, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

// Value returns the JSON value, implements driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

// Custom function used by the gorm ORM if used.
func (j JSON) GormDataType() string {
	return "jsonb"
}

// GobEncode encodes the JSON value using gob encoding.
func (j JSON) GobEncode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(j)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON: %v", err)
	}
	return buffer.Bytes(), nil
}

// GobDecode decodes the gob-encoded data into a JSON value.
func (j *JSON) GobDecode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&j)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}
	return nil
}
