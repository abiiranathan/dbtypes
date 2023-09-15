package dbtypes_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/abiiranathan/dbtypes"
)

func TestDateMarshal(t *testing.T) {
	date := dbtypes.Date(time.Date(2015, 10, 21, 0, 0, 0, 0, time.UTC))
	dateJSON, err := json.Marshal(date)
	if err != nil {
		t.Fatalf("Failed to marshal Date: %v", err)
	}
	if string(dateJSON) != "\"2015-10-21\"" {
		t.Errorf("Unexpected Date JSON: %s", dateJSON)
	}
}

func TestDateUnMarshal(t *testing.T) {
	dateJSON := []byte("\"2015-10-21\"")
	var date dbtypes.Date
	err := json.Unmarshal(dateJSON, &date)
	if err != nil {
		t.Fatalf("Failed to unmarshal Date: %v", err)
	}
	if date.String() != "2015-10-21" {
		t.Errorf("Unexpected Date JSON: %s", dateJSON)
	}
}

func TestParseDateFromString(t *testing.T) {
	date, err := dbtypes.ParseDateFromString("2015-10-21")
	if err != nil {
		t.Fatalf("Failed to parse Date: %v", err)
	}
	if date.String() != "2015-10-21" {
		t.Errorf("Unexpected Date JSON: %s", date)
	}
}

func TestNewDate(t *testing.T) {
	date := dbtypes.NewDate(2015, time.October, 21)
	if date.String() != "2015-10-21" {
		t.Errorf("Unexpected Date JSON: %s", date)
	}
}
