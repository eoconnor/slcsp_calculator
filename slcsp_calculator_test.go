package main

import (
	"testing"
)

type GetRateAreaTestCases struct {
	inputState       string
	inputRateAreaNum string
	expectedOutput   string
}

type ContainsTestCases struct {
	inputSlice     []string
	inputStr       string
	expectedOutput bool
}

func TestGetRateArea(t *testing.T) {
	cases := []GetRateAreaTestCases{
		{"AL", "11", "al 11"},
		{"Ca", "7", "ca 7"},
		{"", "", " "},
	}

	for _, c := range cases {
		output := getRateArea(c.inputState, c.inputRateAreaNum)
		if output != c.expectedOutput {
			t.Errorf("getRateArea(%q, %q) == %q, want %q", c.inputState, c.inputRateAreaNum, output, c.expectedOutput)
		}
	}
}

func TestContains(t *testing.T) {
	cases := []ContainsTestCases{
		{[]string{"al 1", "mn 7", "wy 13"}, "mn 7", true},
		{[]string{"al 1", "mn 7", "wy 13"}, "MN 7", false},
		{[]string{"al 1", "mn 7", "wy 13"}, "mn7", false},
		{[]string{}, "mn 7", false},
	}

	for _, c := range cases {
		output := contains(c.inputSlice, c.inputStr)
		if output != c.expectedOutput {
			t.Errorf("contains(%q, %q) == %t, want %t", c.inputSlice, c.inputStr, output, c.expectedOutput)
		}
	}
}
