package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/messagedigest-net/gh-advanced-security/services"
	"github.com/spf13/cobra"
)

// 1. Parent Command: 'alerts'
// This acts as a container/menu for the specific alert types
var alertsCmd = &cobra.Command{
	Use:   "alerts",
	Short: "List security alerts",
	Long:  `List Code Scanning or Secret Scanning alerts for a specific repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no specific alert type is chosen, show the interactive menu
		services.ChooseSubCommand(cmd.Commands(), args, "Which type of alerts do you want to list?")
	},
}

// 2. Sub-Command: 'code-scanning'
var codeScanningCmd = &cobra.Command{
	Use:     "code-scanning",
	Aliases: []string{"cs", "code"},
	Short:   "List Code Scanning alerts",
	Example: "gh advanced-security list alerts code-scanning owner/repo",
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetAlertServices()

		// Ensure we have a target repo
		args = services.GetTarget(args, "Which repository? (format: owner/repo)")
		owner, repo := parseRepo(args[0])

		// 'json' is the persistent flag defined in root.go
		err := svc.ListCodeScanning(owner, repo, json, UserPageSize)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// 3. Sub-Command: 'secret-scanning'
var secretScanningCmd = &cobra.Command{
	Use:     "secret-scanning",
	Aliases: []string{"ss", "secret"},
	Short:   "List Secret Scanning alerts",
	Example: "gh advanced-security list alerts secret-scanning owner/repo",
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetAlertServices()

		args = services.GetTarget(args, "Which repository? (format: owner/repo)")
		owner, repo := parseRepo(args[0])

		err := svc.ListSecretScanning(owner, repo, json)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Helper to validate and split "owner/repo"
func parseRepo(input string) (string, string) {
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		fmt.Printf("Invalid repository format '%s'. Please use 'owner/repo'\n", input)
		os.Exit(1)
	}
	return parts[0], parts[1]
}

// In cmd/list-alerts.go (or a new cmd/list-bypasses.go)

var listBypassesCmd = &cobra.Command{
	Use:     "bypasses",
	Short:   "List Push Protection bypasses",
	Example: "gh advanced-security list bypasses owner/repo",
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetAlertServices()

		args = services.GetTarget(args, "Which repository? (owner/repo)")
		owner, repo := parseRepo(args[0])

		err := svc.ListPushProtectionBypasses(owner, repo, json)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var dependabotCmd = &cobra.Command{
	Use:     "dependabot",
	Aliases: []string{"dep", "dependencies"},
	Short:   "List Dependabot alerts",
	Example: "gh advanced-security list alerts dependabot owner/repo",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Get the Service (requires services/dependencyservices.go)
		svc := services.GetDependencyServices()

		// 2. Target Resolution
		args = services.GetTarget(args, "Which repository? (format: owner/repo)")
		owner, repo := parseRepo(args[0])

		// 3. Execution
		// 'json' is the persistent flag from root.go
		err := svc.ListDependabotAlerts(owner, repo, json)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// Build the hierarchy: list -> alerts -> [code-scanning, secret-scanning]
	listCmd.AddCommand(alertsCmd)
	alertsCmd.AddCommand(codeScanningCmd)
	alertsCmd.AddCommand(secretScanningCmd)
	alertsCmd.AddCommand(dependabotCmd)
	listCmd.AddCommand(listBypassesCmd)
}
