package cmd

import (
	"github.com/messagedigest-net/gh-advanced-security/services"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show organization, repository...",
	Long: `
	show {organization|org} [Org Name]
	show {repositories|repo} [Repo Name]`,
	Run: func(cmd *cobra.Command, args []string) {
		services.ChooseSubCommand(cmd.Commands(), args, "What do you want to show?")
	},
}
