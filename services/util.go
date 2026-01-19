package services

import "github.com/spf13/cobra"

type ByCommandName []*cobra.Command

func (commands ByCommandName) Len() int           { return len(commands) }
func (commands ByCommandName) Less(i, j int) bool { return commands[i].Name() < commands[j].Name() }
func (commands ByCommandName) Swap(i, j int)      { commands[i], commands[j] = commands[j], commands[i] }
