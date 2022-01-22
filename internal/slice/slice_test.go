package slice

import "testing"

func TestCanFindStringInSlice(t *testing.T) {
	a := []string{"a", "b", "c", "d"}
	if !Contains(a, "a") {
		t.Error("Expected to find 'a' in slice")
	}
	if Contains(a, "e") {
		t.Error("Expected to not find 'e' in slice")
	}
}

func TestCanFindSliceDiff(t *testing.T) {
	a := []string{"a", "b", "c", "d"}
	b := []string{"b", "c", "d", "e"}

	added, deleted := Diff(a, b)
	if len(added) != 1 || added[0] != "e" {
		t.Error("Expected to find 'e' in added slice")
	}
	if len(deleted) != 1 || deleted[0] != "a" {
		t.Error("Expected to find 'a' in deleted slice")
	}
}

func TestCanConvertMapToSlice(t *testing.T) {
	m := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
		"d": "d",
	}
	a := MapToSlice(m)
	if len(a) != 4 {
		t.Error("Expected to find 4 elements in slice")
	}
	if a[3] != "d" {
		t.Error("Expected to find 'd' in slice")
	}
}
