package cmd

import (
	"fmt"
	"os"

	"github.com/messagedigest-net/gh-advanced-security/services"
	"github.com/spf13/cobra"
)

var dependencyGraphCmd = &cobra.Command{
	Use:     "dependency-graph",
	Aliases: []string{"dep", "dg"},
	Short:   "Manage Dependency Graph and Dependabot",
	Long:    `Interact with GitHub Dependency Graph: export SBOMs and view Dependabot alerts.`,
	Run: func(cmd *cobra.Command, args []string) {
		services.ChooseSubCommand(cmd.Commands(), args, "Choose a dependency action:")
	},
}

var sbomCmd = &cobra.Command{
	Use:     "sbom",
	Short:   "Export SBOM (CycloneDX)",
	Long:    `Export the Software Bill of Materials (SBOM) for a repository in CycloneDX format.`,
	Example: `gh advanced-security dependency-graph sbom owner/repo > sbom.json`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetDependencyServices()

		target, _ := services.GetTarget(cmd, args, "Which repository? (owner/repo)")
		owner, repo := parseRepo(target) // Reusing helper from list-alerts.go

		err := svc.ExportSBOM(owner, repo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var dependabotAlertsCmd = &cobra.Command{
	Use:     "alerts",
	Aliases: []string{"list-alerts"},
	Short:   "List Dependabot Alerts",
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetDependencyServices()

		target, flags := services.GetTarget(cmd, args, "Which repository? (owner/repo)")
		owner, repo := parseRepo(target)

		err := svc.ListDependabotAlerts(owner, repo, flags.JSON, flags.PageSize, flags.All)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(dependencyGraphCmd)
	dependencyGraphCmd.AddCommand(sbomCmd)
	dependencyGraphCmd.AddCommand(dependabotAlertsCmd)
}
