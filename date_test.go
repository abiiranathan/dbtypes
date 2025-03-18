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
	date, err := dbtypes.ParseDate("2015-10-21")
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

func TestDate_IsZero(t *testing.T) {
	tests := []struct {
		name string
		date dbtypes.Date
		want bool
	}{
		{
			name: "zero date",
			date: dbtypes.Date{},
			want: true,
		},
		{
			name: "non-zero date",
			date: dbtypes.NewDate(2015, 12, 1),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date := tt.date
			if got := date.IsZero(); got != tt.want {
				t.Errorf("Date.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Equal(t *testing.T) {
	type args struct {
		other dbtypes.Date
	}
	tests := []struct {
		name string
		date dbtypes.Date
		args args
		want bool
	}{
		{
			name: "equal dates",
			date: dbtypes.NewDate(2015, 12, 1),
			args: args{
				other: dbtypes.NewDate(2015, 12, 1),
			},
			want: true,
		},
		{
			name: "unequal dates",
			date: dbtypes.NewDate(2015, 12, 1),
			args: args{
				other: dbtypes.NewDate(2015, 12, 2),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date := tt.date
			if got := date.Equal(tt.args.other); got != tt.want {
				t.Errorf("Date.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_DaysInMonth(t *testing.T) {
	tests := []struct {
		name string
		date dbtypes.Date
		want int
	}{
		{
			name: "31 days",
			date: dbtypes.NewDate(2015, 12, 1),
			want: 31,
		},
		{
			name: "30 days",
			date: dbtypes.NewDate(2015, 11, 1),
			want: 30,
		},
		{
			name: "29 days",
			date: dbtypes.NewDate(2016, 2, 1),
			want: 29,
		},
		{
			name: "28 days",
			date: dbtypes.NewDate(2015, 2, 1),
			want: 28,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date := tt.date
			if got := date.DaysInMonth(); got != tt.want {
				t.Errorf("Date.DaysInMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaysInYear(t *testing.T) {
	tests := []struct {
		name string
		date dbtypes.Date
		want int
	}{
		{
			name: "Non-leap year",
			date: dbtypes.NewDate(2021, time.January, 1),
			want: 365,
		},
		{
			name: "Leap year",
			date: dbtypes.NewDate(2020, time.January, 1),
			want: 366,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date.DaysInYear(); got != tt.want {
				t.Errorf("DaysInYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaysBetween(t *testing.T) {
	tests := []struct {
		name  string
		date1 dbtypes.Date
		date2 dbtypes.Date
		want  int
	}{
		{
			name:  "Same day",
			date1: dbtypes.NewDate(2021, time.January, 1),
			date2: dbtypes.NewDate(2021, time.January, 1),
			want:  0,
		},
		{
			name:  "One day apart",
			date1: dbtypes.NewDate(2021, time.January, 1),
			date2: dbtypes.NewDate(2021, time.January, 2),
			want:  1,
		},
		{
			name:  "One year apart",
			date1: dbtypes.NewDate(2021, time.January, 1),
			date2: dbtypes.NewDate(2022, time.January, 1),
			want:  365,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date1.DaysBetween(tt.date2); got != tt.want {
				t.Errorf("DaysBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Before(t *testing.T) {
	tests := []struct {
		name  string
		date1 dbtypes.Date
		date2 dbtypes.Date
		want  bool
	}{
		{
			name:  "date1 before date2",
			date1: dbtypes.NewDate(2023, 10, 1),
			date2: dbtypes.NewDate(2023, 10, 2),
			want:  true,
		},
		{
			name:  "date1 after date2",
			date1: dbtypes.NewDate(2023, 10, 2),
			date2: dbtypes.NewDate(2023, 10, 1),
			want:  false,
		},
		{
			name:  "same date",
			date1: dbtypes.NewDate(2023, 10, 1),
			date2: dbtypes.NewDate(2023, 10, 1),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date1.Before(tt.date2); got != tt.want {
				t.Errorf("Date.Before() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_After(t *testing.T) {
	tests := []struct {
		name  string
		date1 dbtypes.Date
		date2 dbtypes.Date
		want  bool
	}{
		{
			name:  "date1 after date2",
			date1: dbtypes.NewDate(2023, 10, 2),
			date2: dbtypes.NewDate(2023, 10, 1),
			want:  true,
		},
		{
			name:  "date1 before date2",
			date1: dbtypes.NewDate(2023, 10, 1),
			date2: dbtypes.NewDate(2023, 10, 2),
			want:  false,
		},
		{
			name:  "same date",
			date1: dbtypes.NewDate(2023, 10, 1),
			date2: dbtypes.NewDate(2023, 10, 1),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date1.After(tt.date2); got != tt.want {
				t.Errorf("Date.After() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Equal_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		date1 dbtypes.Date
		date2 dbtypes.Date
		want  bool
	}{
		{
			name:  "zero dates",
			date1: dbtypes.Date{},
			date2: dbtypes.Date{},
			want:  true,
		},
		{
			name:  "different time zones",
			date1: dbtypes.Date(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)),
			date2: dbtypes.Date(time.Date(2023, 10, 1, 0, 0, 0, 0, time.FixedZone("EST", -5*3600))),
			want:  true,
		},
		{
			name:  "subsecond difference",
			date1: dbtypes.Date(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)),
			date2: dbtypes.Date(time.Date(2023, 10, 1, 0, 0, 0, 1, time.UTC)),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date1.Equal(tt.date2); got != tt.want {
				t.Errorf("Date.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDateFromString_InvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "malformed date",
			input:   "2023-13-01",
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "01-10-2023",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dbtypes.ParseDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDateFromString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDate_DaysInMonth_LeapYear(t *testing.T) {
	tests := []struct {
		name string
		date dbtypes.Date
		want int
	}{
		{
			name: "February in leap year",
			date: dbtypes.NewDate(2020, 2, 1),
			want: 29,
		},
		{
			name: "February in non-leap year",
			date: dbtypes.NewDate(2021, 2, 1),
			want: 28,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date.DaysInMonth(); got != tt.want {
				t.Errorf("Date.DaysInMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaysBetween_Negative(t *testing.T) {
	tests := []struct {
		name  string
		date1 dbtypes.Date
		date2 dbtypes.Date
		want  int
	}{
		{
			name:  "date2 before date1",
			date1: dbtypes.NewDate(2023, 10, 2),
			date2: dbtypes.NewDate(2023, 10, 1),
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date1.DaysBetween(tt.date2); got != tt.want {
				t.Errorf("DaysBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}
