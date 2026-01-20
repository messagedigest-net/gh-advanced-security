/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/messagedigest-net/gh-advanced-security/services"
	"github.com/spf13/cobra"
)

// organizationCmd represents the organization command
var organizationCmd = &cobra.Command{
	Use:     "organization",
	Aliases: []string{"org"},
	Short:   "Show informations about a organization",
	Long: `Shows the following configuration for a given organization:
	- Dependecy Graph
	- Dependabot Alerts
	- Dependabot Security Updates
	- Enable Advanced Security for new Repos
	- Secret Scanning
	- Secret Scanning Push Protection
	- Secret Scanning Push Protection Custom Link
	- Secret Scanning Push Protection Custom Link Enabled`,
	Run: func(cmd *cobra.Command, args []string) {
		var service services.Shower
		var err error

		service = services.GetOrganizationServices()

		target, flags := services.GetTarget(cmd, args, "Which organization do you want to show?")

		err = service.Show(target, flags.JSON)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	showCmd.AddCommand(organizationCmd)
}
