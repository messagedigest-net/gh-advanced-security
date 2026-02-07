package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/messagedigest-net/gh-advanced-security/services"
	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:     "enable",
	Aliases: []string{"en"},
	Short:   "Enable security features",
	Long:    `Enable security features like Secret Scanning, Push Protection and Dependabot.`,
	Run: func(cmd *cobra.Command, args []string) {
		services.ChooseSubCommand(cmd.Commands(), args, "What do you want to enable?")
	},
}

var pushProtectionCmd = &cobra.Command{
	Use:     "push-protection",
	Aliases: []string{"pp"},
	Short:   "Enable Push Protection",
	Long:    `Enable Secret Scanning and Push Protection for a repository or all repositories in an organization.`,
	Example: `
  # Enable for a single repo
  gh advanced-security enable push-protection owner/repo

  # Enable for an entire organization
  gh advanced-security enable push-protection my-org`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetEnforcerServices()

		target, _ := services.GetTarget(cmd, args, "For which org or repo do you want to enable Push Protection?")

		if strings.Contains(target, "/") {
			// === Single Repo Mode ===
			parts := strings.Split(target, "/")
			owner, repo := parts[0], parts[1]

			fmt.Printf("Enabling Push Protection for %s/%s...\n", owner, repo)
			err := svc.EnablePushProtection(owner, repo)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			fmt.Println("Success!")
		} else {
			// Chama o método otimizado (O(1))
			err := svc.BulkEnablePushProtection(target)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

var secretScanningEnableCmd = &cobra.Command{
	Use:     "secret-scanning",
	Aliases: []string{"ss"},
	Short:   "Enable Secret Scanning",
	Long:    `Enable Secret Scanning (without Push Protection) for a repository or organization.`,
	Example: `
  # Enable for a single repo
  gh advanced-security enable secret-scanning owner/repo

  # Enable for an entire organization
  gh advanced-security enable secret-scanning my-org`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetEnforcerServices()

		target, _ := services.GetTarget(cmd, args, "For which org or repo do you want to enable Secret Scanning?")

		if strings.Contains(target, "/") {
			// === Single Repo Mode ===
			parts := strings.Split(target, "/")
			owner, repo := parts[0], parts[1]

			fmt.Printf("Enabling Secret Scanning for %s/%s...\n", owner, repo)
			err := svc.EnableSecretScanning(owner, repo)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			fmt.Println("Success!")
		} else {
			// Chama o método otimizado (O(1))
			err := svc.BulkEnableSecretScanning(target)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

var secretScanningNonProviderPatternsEnableCmd = &cobra.Command{
	Use:     "non-provider-patterns",
	Aliases: []string{"npp"},
	Short:   "Enable Secret Scanning Non-Provider Patterns",
	Long:    `Enable Secret Scanning Non-Provider Patterns for a repository.`,
	Example: `
  # Enable for a single repo
  gh advanced-security enable non-provider-patterns owner/repo`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetEnforcerServices()

		target, _ := services.GetTarget(cmd, args, "For which repo do you want to enable Secret Scanning Non-Provider Patterns?")

		if strings.Contains(target, "/") {
			// === Single Repo Mode ===
			parts := strings.Split(target, "/")
			owner, repo := parts[0], parts[1]

			fmt.Printf("Enabling Secret Scanning Non-Provider Patterns for %s/%s...\n", owner, repo)
			err := svc.EnableSecretScanningNonProviderPatterns(owner, repo)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			fmt.Println("Success!")
		} else {
			//
			fmt.Println("This setting can only be applied on individual repositories.")
			os.Exit(1)
		}
	},
}

var dependabotEnableCmd = &cobra.Command{
	Use:     "dependabot",
	Aliases: []string{"dep"},
	Short:   "Enable Dependabot features",
	Long:    `Enable Dependency Graph, Alerts, and Security Updates.`,
	Example: `
  # Enable for a single repo
  gh advanced-security enable dependabot owner/repo

  # Enable for an entire organization
  gh advanced-security enable dependabot my-org`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetEnforcerServices()

		target, _ := services.GetTarget(cmd, args, "Target (Org or Owner/Repo)?")

		if strings.Contains(target, "/") {
			// === Single Repo Mode ===
			parts := strings.Split(target, "/")
			owner, repo := parts[0], parts[1]

			fmt.Printf("Enabling Dependabot for %s/%s...\n", owner, repo)

			// 1. Dependabot Alerts (PUT)
			if err := svc.EnableDependabotAlerts(owner, repo); err != nil {
				fmt.Printf("Error enabling alerts: %s\n", err)
				os.Exit(1)
			}

			// 2. Security Updates (PATCH)
			if err := svc.EnableDependabotSecurityUpdates(owner, repo); err != nil {
				fmt.Printf("Error enabling updates: %s\n", err)
				os.Exit(1)
			}
			fmt.Println("Success! (Dependency Graph is implied/enabled by Alerts)")

		} else {
			// Chama o método otimizado (O(1))
			err := svc.BulkEnableDependabot(target)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)
	enableCmd.AddCommand(pushProtectionCmd)
	enableCmd.AddCommand(secretScanningEnableCmd)
	enableCmd.AddCommand(secretScanningNonProviderPatternsEnableCmd)
	enableCmd.AddCommand(dependabotEnableCmd)
}
