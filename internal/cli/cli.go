// internal/cli/cli.go
package cli

import "strings"

func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
