package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/messagedigest-net/gh-advanced-security/services"
	"github.com/spf13/cobra"
)

var showRepoCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"repo"},
	Short:   "Show repository details",
	Long:    `Show detailed security configuration for a specific repository.`,
	Example: `gh advanced-security show repo owner/repo`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetRepositoryServices()

		// Interactive target selection
		target, flags := services.GetTarget(cmd, args, "Which repository do you want to show? (owner/repo)")

		// Validate input format
		if !strings.Contains(target, "/") {
			fmt.Println("Error: Please use format 'owner/repo'")
			os.Exit(1)
		}

		// Call the Service
		// 'json' is the persistent flag from root.go
		err := svc.Show(target, flags.JSON)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// Hook into the parent 'show' command
	showCmd.AddCommand(showRepoCmd)
}
