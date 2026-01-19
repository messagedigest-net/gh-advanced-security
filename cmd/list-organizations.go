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
var organizationsCmd = &cobra.Command{
	Use:     "organizations",
	Aliases: []string{"orgs"},
	Short:   "List organizations for current user",
	Long: `List all the organizations for the current user and the following security configurations:
	- Dependecy Graph
	- Dependabot Alerts
	- Dependabot Security Updates
	- Enable Advanced Security for new Repos
	- Secret Scanning
	- Secret Scanning Push Protection
	- Secret Scanning Push Protection Custom Link
	- Secret Scanning Push Protection Custom Link Enabled`,
	Run: func(cmd *cobra.Command, args []string) {
		var service services.Lister
		var err error

		service = services.GetOrganizationServices()

		err = service.List(json)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	listCmd.AddCommand(organizationsCmd)
}
