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
		args = services.GetTarget(args, "Which repository do you want to show? (owner/repo)")

		// Validate input format
		if !strings.Contains(args[0], "/") {
			fmt.Println("Error: Please use format 'owner/repo'")
			os.Exit(1)
		}

		// Call the Service
		// 'json' is the persistent flag from root.go
		err := svc.Show(args[0], json)
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
