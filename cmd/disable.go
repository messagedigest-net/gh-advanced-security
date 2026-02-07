package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/messagedigest-net/gh-advanced-security/services"
	"github.com/spf13/cobra"
)

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable security features",
	Long:  `Disable security features like Secret Scanning, Push Protection and Dependabot.`,
	Run: func(cmd *cobra.Command, args []string) {
		services.ChooseSubCommand(cmd.Commands(), args, "What do you want to disable?")
	},
}

var pushProtectionDisableCmd = &cobra.Command{
	Use:   "push-protection",
	Short: "Disable Push Protection",
	Example: `gh advanced-security disable push-protection owner/repo
gh advanced-security disable push-protection my-org`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetEnforcerServices()
		target, _ := services.GetTarget(cmd, args, "Target (Org or Owner/Repo)?")

		if strings.Contains(target, "/") {
			parts := strings.Split(target, "/")
			owner, repo := parts[0], parts[1]
			fmt.Printf("Disabling Push Protection for %s/%s...\n", owner, repo)
			if err := svc.DisablePushProtection(owner, repo); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Success!")
		} else {
			confirmAction(target, "Push Protection", func() error {
				return svc.BulkDisablePushProtection(target)
			})
		}
	},
}

var secretScanningDisableCmd = &cobra.Command{
	Use:     "secret-scanning",
	Aliases: []string{"ss"},
	Short:   "Disable Secret Scanning",
	Example: `gh advanced-security disable secret-scanning owner/repo`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetEnforcerServices()
		target, _ := services.GetTarget(cmd, args, "Target (Org or Owner/Repo)?")

		if strings.Contains(target, "/") {
			parts := strings.Split(target, "/")
			owner, repo := parts[0], parts[1]
			fmt.Printf("Disabling Secret Scanning for %s/%s...\n", owner, repo)
			if err := svc.DisableSecretScanning(owner, repo); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Success!")
		} else {
			confirmAction(target, "Secret Scanning", func() error {
				return svc.BulkDisableSecretScanning(target)
			})
		}
	},
}

var secretScanningNonProviderPatternsDisableCmd = &cobra.Command{
	Use:     "non-provider-patterns",
	Aliases: []string{"npp"},
	Short:   "Disable Secret Scanning Non-Provider Patterns",
	Example: `gh advanced-security disable non-provider-patterns owner/repo`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetEnforcerServices()
		target, _ := services.GetTarget(cmd, args, "Target (Owner/Repo)?")

		if strings.Contains(target, "/") {
			parts := strings.Split(target, "/")
			owner, repo := parts[0], parts[1]
			fmt.Printf("Disabling Secret Scanning Non-Provider Patterns for %s/%s...\n", owner, repo)
			if err := svc.DisableSecretScanningNonProviderPatterns(owner, repo); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Success!")
		} else {
			fmt.Println("This setting can only be applied on individual repositories.")
			os.Exit(1)
		}
	},
}

var dependabotDisableCmd = &cobra.Command{
	Use:     "dependabot",
	Aliases: []string{"dep"},
	Short:   "Disable Dependabot features",
	Example: `gh advanced-security disable dependabot owner/repo`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetEnforcerServices()
		target, _ := services.GetTarget(cmd, args, "Target (Org or Owner/Repo)?")

		if strings.Contains(target, "/") {
			parts := strings.Split(target, "/")
			owner, repo := parts[0], parts[1]
			fmt.Printf("Disabling Dependabot for %s/%s...\n", owner, repo)

			// Updates primeiro
			svc.DisableDependabotSecurityUpdates(owner, repo)
			// Alerts depois
			if err := svc.DisableDependabotAlerts(owner, repo); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Success! (Alerts and Updates disabled)")
		} else {
			confirmAction(target, "Dependabot (Graph, Alerts, Updates)", func() error {
				return svc.BulkDisableDependabot(target)
			})
		}
	},
}

// Helper para evitar repetição do prompt de confirmação
func confirmAction(target, feature string, action func() error) {
	fmt.Printf("Disabling %s for ALL repositories in '%s'.\n", feature, target)
	fmt.Printf("Are you sure? (y/N): ")
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) != "y" {
		fmt.Println("Aborted.")
		os.Exit(0)
	}
	if err := action(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Success! Changes will be applied asynchronously.")
}

func init() {
	rootCmd.AddCommand(disableCmd)
	disableCmd.AddCommand(pushProtectionDisableCmd)
	disableCmd.AddCommand(secretScanningDisableCmd)
	disableCmd.AddCommand(secretScanningNonProviderPatternsDisableCmd)
	disableCmd.AddCommand(dependabotDisableCmd)
}
