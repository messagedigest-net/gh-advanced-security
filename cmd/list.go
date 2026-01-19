package cmd

import (
	"github.com/messagedigest-net/gh-advanced-security/services"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List organizations, repositories, alerts...",
	Long: `list organizations
		   list repositories [Org]
		   list codescanning alerts [Org/Repo]`,
	Run: func(cmd *cobra.Command, args []string) {
		services.ChooseSubCommand(cmd.Commands(), args, "What do you want to list?")
	},
}
