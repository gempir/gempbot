package slice

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func Diff(s1, s2 []string) (added, deleted []string) {
	len1, len2 := len(s1), len(s2)

	added = make([]string, 0, len2)
	deleted = make([]string, 0, len1)

	for i1, i2 := 0, 0; i1 < len1 || i2 < len2; {
		if i1 == len1 {
			added = append(added, s2[i2:]...)
			return
		} else if i2 == len2 {
			deleted = append(deleted, s1[i1:]...)
			return
		}

		if s1[i1] < s2[i2] {
			deleted = append(deleted, s1[i1])
			i1++
		} else if s1[i1] > s2[i2] {
			added = append(added, s2[i2])
			i2++
		} else {
			i1++
			i2++
		}
	}
	return
}

func MapToSlice(m map[string]string) []string {
	a := make([]string, 0, len(m))
	for _, v := range m {
		a = append(a, v)
	}

	return a
}
