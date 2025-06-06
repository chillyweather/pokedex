package main

import (
	"testing"

	"github.com/chillyweather/pokedexcli/internal/cli"
	"github.com/chillyweather/pokedexcli/internal/config"
)

func TestCleanInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single word",
			input:    "help",
			expected: []string{"help"},
		},
		{
			name:     "multiple words",
			input:    "map location area",
			expected: []string{"map", "location", "area"},
		},
		{
			name:     "mixed case",
			input:    "HELP Me",
			expected: []string{"help", "me"},
		},
		{
			name:     "extra spaces",
			input:    "  help   me  ",
			expected: []string{"help", "me"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cli.CleanInput(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("CleanInput(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("CleanInput(%q) = %v, want %v", tt.input, result, tt.expected)
					break
				}
			}
		})
	}
}

func TestGetCommands(t *testing.T) {
	commands := cli.GetCommands()

	expectedCommands := []string{"help", "exit", "map", "mapb"}

	if len(commands) != len(expectedCommands) {
		t.Errorf("Expected %d commands, got %d", len(expectedCommands), len(commands))
	}

	for _, cmdName := range expectedCommands {
		if cmd, exists := commands[cmdName]; !exists {
			t.Errorf("Command %q not found", cmdName)
		} else {
			if cmd.Name != cmdName {
				t.Errorf("Command name mismatch: got %q, want %q", cmd.Name, cmdName)
			}
			if cmd.Description == "" {
				t.Errorf("Command %q has empty description", cmdName)
			}
			if cmd.Callback == nil {
				t.Errorf("Command %q has nil callback", cmdName)
			}
		}
	}
}

func TestConfigNew(t *testing.T) {
	cfg := config.New()

	if cfg == nil {
		t.Fatal("config.New() returned nil")
	}

	expectedNext := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	if cfg.Next != expectedNext {
		t.Errorf("Expected Next URL %q, got %q", expectedNext, cfg.Next)
	}

	if cfg.Previous != "" {
		t.Errorf("Expected empty Previous URL, got %q", cfg.Previous)
	}
}

func TestCommandHelpDoesNotPanic(t *testing.T) {
	commands := cli.GetCommands()
	helpCmd := commands["help"]
	cfg := config.New()

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("commandHelp panicked: %v", r)
		}
	}()

	err := helpCmd.Callback(cfg)
	if err != nil {
		t.Errorf("commandHelp returned error: %v", err)
	}
}
