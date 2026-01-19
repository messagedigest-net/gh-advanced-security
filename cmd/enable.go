package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/messagedigest-net/gh-advanced-security/services"
	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable security features",
	Long:  `Enable security features like Secret Scanning and Push Protection.`,
	Run: func(cmd *cobra.Command, args []string) {
		services.ChooseSubCommand(cmd.Commands(), args, "What do you want to enable?")
	},
}

var pushProtectionCmd = &cobra.Command{
	Use:   "push-protection",
	Short: "Enable Push Protection",
	Long:  `Enable Secret Scanning and Push Protection for a repository or all repositories in an organization.`,
	Example: `
  # Enable for a single repo
  gh advanced-security enable push-protection owner/repo

  # Enable for an entire organization (Interactive)
  gh advanced-security enable push-protection my-org --all`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetEnforcerServices()

		args = services.GetTarget(args, "For which org or repo do you want to enable Push Protection?")

		target := args[0]

		// Check if it's a Repo (has slash) or Org (no slash)
		if strings.Contains(target, "/") {
			// Single Repo Mode
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
			// Organization Mode
			// You might want to add a confirmation prompt here
			fmt.Printf("Enabling Push Protection for ALL repos in %s.\nAre you sure? (y/N): ", target)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				fmt.Println("Aborted.")
				os.Exit(0)
			}

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
  gh advanced-security enable secret-scanning my-org --all`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := services.GetEnforcerServices()

		args = services.GetTarget(args, "For which org or repo do you want to enable Secret Scanning?")

		target := args[0]

		if strings.Contains(target, "/") {
			// Single Repo Mode
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
			// Organization Mode
			fmt.Printf("Enabling Secret Scanning for ALL repos in %s.\nAre you sure? (y/N): ", target)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				fmt.Println("Aborted.")
				os.Exit(0)
			}

			err := svc.BulkEnableSecretScanning(target)
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
}
