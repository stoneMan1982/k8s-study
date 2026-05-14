package generic

import "testing"

func TestSumInts(t *testing.T) {
	m := map[string]int64{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	expected := int64(6)
	result := SumInts(m)
	if result != expected {
		t.Errorf("SumInts(%v) = %d; want %d", m, result, expected)
	}
}

func TestSumFloats(t *testing.T) {
	m := map[string]float64{
		"x": 1.5,
		"y": 2.5,
		"z": 3.0,
	}
	expected := 7.0
	result := SumFloats(m)
	if result != expected {
		t.Errorf("SumFloats(%v) = %f; want %f", m, result, expected)
	}
}
