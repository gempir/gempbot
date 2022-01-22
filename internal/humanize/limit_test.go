package humanize

import (
	"testing"
	"fmt"
)

func TestCharLimiterTable(t *testing.T) {
	var tests = []struct {
		str string
		num int
		want string
	}{
		{"Doctor", 2, "Do"+"..."},//Regular case
		{"", 0, ""},//Empty everything
		{"Very long string", 10,"Very lo"+"..."},//Long string input
		{"A", 1, "A"},//Small string
	}

	for _, tt := range tests{
		testname:=fmt.Sprintf("%s,%d", tt.str, tt.num)
		t.Run(testname, func(t *testing.T) {
			ans:=CharLimiter(tt.str, tt.num)
			if ans != tt.want {
				t.Errorf("Got %s, want %s", ans, tt.want)
			}
		})
	}
}