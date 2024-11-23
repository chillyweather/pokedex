package main

import (
	"testing"
)

func TestNormalizedInput(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{
			input: "Hello world",
			want: []string{
				"hello",
				"world",
			},
		},
		{
			input: "Wabby Dabby Doo",
			want: []string{
				"wabby",
				"dabby",
				"doo",
			},
		},
	}

	for _, cs := range cases {
		got := normalizeInput(cs.input)

		if len(got) != len(cs.want) {
			t.Errorf("The length are not equal: %v vs %v", got, cs.want)
			continue
		}

		for i := range got {
			actualWord := got[i]
			expectedWord := cs.want[i]

			if actualWord != expectedWord {
				t.Errorf("Expected %v, but got %v", expectedWord, actualWord)
			}
		}
	}
}
