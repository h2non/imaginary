package main

import "testing"

func TestToMegaBytes(t *testing.T) {
	tests := []struct {
		value    uint64
		expected float64
	}{
		{1024, 0},
		{1024 * 1024, 1},
		{1024 * 1024 * 10, 10},
		{1024 * 1024 * 100, 100},
		{1024 * 1024 * 250, 250},
	}

	for _, test := range tests {
		val := toMegaBytes(test.value)
		if val != test.expected {
			t.Errorf("Invalid param: %#v != %#v", val, test.expected)
		}
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		value    float64
		expected int
	}{
		{0, 0},
		{1, 1},
		{1.56, 2},
		{1.38, 1},
		{30.12, 30},
	}

	for _, test := range tests {
		val := round(test.value)
		if val != test.expected {
			t.Errorf("Invalid param: %#v != %#v", val, test.expected)
		}
	}
}

func TestToFixed(t *testing.T) {
	tests := []struct {
		value    float64
		expected float64
	}{
		{0, 0},
		{1, 1},
		{123, 123},
		{0.99, 1},
		{1.02, 1},
		{1.82, 1.8},
		{1.56, 1.6},
		{1.38, 1.4},
	}

	for _, test := range tests {
		val := toFixed(test.value, 1)
		if val != test.expected {
			t.Errorf("Invalid param: %#v != %#v", val, test.expected)
		}
	}
}
