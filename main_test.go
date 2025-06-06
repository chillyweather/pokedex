package main

import (
	"testing"
)

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
			input:    "  Exit  ",
			expected: []string{"exit"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "   ",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Expected length %d, got %d for input '%s'", len(c.expected), len(actual), c.input)
			continue
		}
		for i, word := range actual {
			if word != c.expected[i] {
				t.Errorf("Expected '%s', got '%s' at position %d for input '%s'", c.expected[i], word, i, c.input)
			}
		}
	}
}
