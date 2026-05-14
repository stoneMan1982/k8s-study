package generic

import "testing"

func TestSumIntsOrFloats(t *testing.T) {
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}
	if SumIntsOrFloats(ints) != 46 {
		t.Error("SumIntsOrFloats failed with ints")
	}
	if SumIntsOrFloats(floats) != 62.97 {
		t.Error("SumIntsOrFloats failed with floats")
	}
}

func TestMax(t *testing.T) {
	if Max(3, 5) != 5 {
		t.Error("Max failed with ints")
	}
	if Max(2.7, 1.3) != 2.7 {
		t.Error("Max failed with floats")
	}
	if Max("apple", "banana") != "banana" {
		t.Error("Max failed with strings")
	}
}
