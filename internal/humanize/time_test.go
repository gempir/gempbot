package humanize

import (
	"fmt"
	"testing"
)

func TestStringToSeconds(t *testing.T) {
	var tests = []struct {
		str  string
		want int
	}{
		{"60s", 60},
		{"1m", 60},
		{"60m", 3600}, //Big example
		{"22s", 22},   //Basic example
		{"", 0},       //Empty string
		{"10d", 0},    //Invalid input
	}

	for _, tt := range tests {
		testname := tt.str
		t.Run(testname, func(t *testing.T) {
			ans, err := StringToSeconds(tt.str)
			if ans != tt.want {
				fmt.Printf("%s", err)
				t.Errorf("Got %d, want %d", ans, tt.want)
			}
		})
	}
}

func TestSecondsToString(t *testing.T) {
	var tests = []struct {
		s    int
		want string
	}{
		{0, "0s"},
		{55, "55s"},
		{60, "1m"},
		{120, "2m"},
		{121, "2m 1s"},
		{4000, "66m 40s"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.s)
		t.Run(testname, func(t *testing.T) {
			ans := SecondsToString(tt.s)
			if ans != tt.want {
				t.Errorf("Got %s, want %s", ans, tt.want)
			}
		})
	}
}
