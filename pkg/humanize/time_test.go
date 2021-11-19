func TestSecondsToString(t *testing.T) {
	var tests = []struct {
		s int
		want string
	}{
		{0, "0s"},
		{55, "55s"},
		{60, "1m"},
		{120, "2m"},
		{121, "2m 1s"},
		{4000, "66m 40s"},
	}

	for _, tt := range tests{
		testname:=fmt.Sprintf("%d", tt.s)
		t.Run(testname, func(t *testing.T) {
			ans:=SecondsToString(tt.s)
			if ans != tt.want {
				t.Errorf("Got %s, want %s", ans, tt.want)
			}
		})
	}
}
