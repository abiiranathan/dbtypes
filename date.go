package dbtypes

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"
)

type Date time.Time

// The Standard date layout.
const DateLayout = "2006-01-02"

// Scan implements the sql.Scanner interface
func (date *Date) Scan(value any) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Date(nullTime.Time)
	return
}

// Value implements the driver.Valuer interface.
func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}

// Custom function used by the gorm ORM if used.
func (date Date) GormDataType() string {
	return "date"
}

// GobEncode encodes the date with gob encoding.
func (date Date) GobEncode() ([]byte, error) {
	return time.Time(date).GobEncode()
}

// GobDecode decodes bytes in b into a date object.
func (date *Date) GobDecode(b []byte) error {
	return (*time.Time)(date).GobDecode(b)
}

// Marshals Date type with the standard date layout.
// If date is a zero value, it will return null bytes.
func (date Date) MarshalJSON() ([]byte, error) {
	datetime := time.Time(date)

	if datetime.IsZero() {
		return []byte("null"), nil
	}

	dateBytes, err := datetime.MarshalJSON()
	if err != nil {
		return []byte(""), err
	}

	// Transform the date to format of layout in format yyyy-mm-dd
	dateString := string(dateBytes[1 : len(dateBytes)-1])
	dateString = fmt.Sprintf("\"%s-%02s-%02s\"", dateString[0:4], dateString[5:7], dateString[8:10])
	return []byte(dateString), nil
}

// Custom Json decoder
// Called to convert json strings to go types
func (date *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("date should be a string, got %v", data)
	}

	// by convention, unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	// If s is an empty string, assume a zero-value to allow for optional date.
	if strings.TrimSpace(s) == "" {
		s = "0001-01-01"
	}

	// Make sure that the user has provided the standard date format
	_, err := time.Parse(DateLayout, s)
	if err != nil {
		return fmt.Errorf("date should be of the format: yyyy-mm-dd")
	}

	// Convert date string to the standard format to RFC 3339 format
	s = fmt.Sprintf("\"%sT00:00:00Z\"", s)
	return (*time.Time)(date).UnmarshalJSON([]byte(s))
}

// Implement a FormScanner interface to be parsed from a
// multipart/form or www-x-urlencoded form.
// If value is an empty string, no parsing is performed.
// You should validate the date after parsing the form/json.
// See https://github.com/abiiranathan/rex.git
func (date *Date) FormScan(value any) error {
	dateStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid date. Expected value as a string")
	}

	// Skip empty Date.
	if dateStr == "" {
		return nil
	}

	parsedDate, err := ParseDate(dateStr)
	if err != nil {
		return err
	}

	*date = parsedDate
	return nil
}

func (date Date) Year() int {
	return time.Time(date).Year()
}

func (date Date) Month() int {
	return int(time.Time(date).Month())
}

func (date Date) Day() int {
	return time.Time(date).Day()
}

// Format the date using standard golang time.Format layout.
func (date Date) Format(layout string) string {
	if date.IsZero() {
		return ""
	}
	return time.Time(date).Format(layout)
}

// String returns a string version of the date using layout 2006-01-02.
func (date Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", date.Year(), date.Month(), date.Day())
}

// Create a new date object.
func NewDate(year int, month time.Month, day int) Date {
	return Date(time.Date(year, month, day, 0, 0, 0, 0, time.Local))
}

// Parse date from a string. Returns a new Date object or error
// if dateStr is empty of has invalid format.
func ParseDate(dateStr string) (Date, error) {
	var date Date

	if dateStr == "" {
		return date, fmt.Errorf("date string is empty")
	}

	b, err := json.Marshal(dateStr)
	if err != nil {
		return date, err
	}

	err = date.UnmarshalJSON(b)
	if err != nil {
		return date, err
	}
	return date, nil
}

// Today returns today's date.
func Today() Date {
	return NewDate(time.Now().Date())
}

// IsZero returns true if the Date is zero.
func (date Date) IsZero() bool {
	return time.Time(date).IsZero()
}

// Equal returns true is the date is equal to other date. The difference in hours is ignored.
func (date Date) Equal(other Date) bool {
	t1 := time.Time(date).Truncate(24 * time.Hour).UTC()
	t2 := time.Time(other).Truncate(24 * time.Hour).UTC()
	return t1.Equal(t2)
}

// Before returns true is the date is before other date. The difference in hours is ignored.
func (date Date) Before(other Date) bool {
	t1 := time.Time(date).Truncate(24 * time.Hour).UTC()
	t2 := time.Time(other).Truncate(24 * time.Hour).UTC()
	return t1.Before(t2)
}

// After returns true is the date is after other date. The difference in hours is ignored.
func (date Date) After(other Date) bool {
	t1 := time.Time(date).Truncate(24 * time.Hour).UTC()
	t2 := time.Time(other).Truncate(24 * time.Hour).UTC()
	return t1.After(t2)
}

// AddDate add years, months and days to the date and returns the new date.
func (date Date) AddDate(years int, months int, days int) Date {
	return Date(time.Time(date).AddDate(years, months, days))
}

func (date Date) AddDays(days int) Date {
	return date.AddDate(0, 0, days)
}

func (date Date) AddMonths(months int) Date {
	return date.AddDate(0, months, 0)
}

func (date Date) AddYears(years int) Date {
	return date.AddDate(years, 0, 0)
}

// Returns the number of days in the month of the date.
func (date Date) DaysInMonth() int {
	nextMonth := time.Time(date).AddDate(0, 1, 0)
	lastDayOfMonth := time.Date(nextMonth.Year(),
		nextMonth.Month(), 0, 0, 0, 0, 0, nextMonth.Location())
	return lastDayOfMonth.Day()
}

func (date Date) DaysInYear() int {
	year := time.Time(date).Year()
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return 366
	}
	return 365
}

// ToTime converts Date to time.Time
func (d Date) ToTime() time.Time {
	return time.Time(d)
}

// DaysBetween returns the absolute number of days between the date and the other date.
func (d Date) DaysBetween(other Date) int {
	start := d.ToTime().Truncate(24 * time.Hour)
	end := other.ToTime().Truncate(24 * time.Hour)
	return int(math.Abs(end.Sub(start).Hours() / 24))
}
