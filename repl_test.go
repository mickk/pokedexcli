package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "This is a test",
			expected: []string{"This", "is", "a", "test"},
		},
		{
			input:    " has this string been split correctly  ",
			expected: []string{"has", "this", "string", "been", "split", "correctly"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Actual slice length does not equal expected slice length")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected: %v does not match Actual %v ", expectedWord, word)
			}
		}
	}
}
