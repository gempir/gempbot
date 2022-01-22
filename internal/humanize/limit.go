package humanize

func CharLimiter(str string, num int) string {
	result := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		result = str[0:num] + "..."
	}
	return result
}
