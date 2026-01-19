package services

import (
	"strings"

	"github.com/spf13/cobra"
)

type ByCommandName []*cobra.Command

func (commands ByCommandName) Len() int           { return len(commands) }
func (commands ByCommandName) Less(i, j int) bool { return commands[i].Name() < commands[j].Name() }
func (commands ByCommandName) Swap(i, j int)      { commands[i], commands[j] = commands[j], commands[i] }

// AskForNextPage pauses execution and asks the user if they want more results
func AskForNextPage() bool {
	prompt := GetPrompt()

	// We use Input because we just want a quick confirmation
	// A standard practice is "Press Enter for more, or q to quit"
	resp, err := prompt.Input("Press Enter for next page (or 'q' to quit)", "")
	if err != nil {
		return false
	}

	// Check if user wants to quit
	if strings.ToLower(strings.TrimSpace(resp)) == "q" {
		return false
	}

	return true
}
