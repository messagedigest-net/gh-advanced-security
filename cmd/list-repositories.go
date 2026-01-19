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

// repositoriesCmd represents the repositories command
var repositoriesCmd = &cobra.Command{
	Use:     "repositories",
	Aliases: []string{"repos"},
	Short:   "List repositories for current user",
	Long:    `List repositories for current user`,
	Run: func(cmd *cobra.Command, args []string) {
		var service services.ListerFor
		var err error

		args = services.GetTarget(args, "For which org (for user, set the [-u] flag) do you want to list the repos?")

		service = services.GetRepositoryServices()

		err = service.ListFor(args[0], user, json, UserPageSize, all)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	listCmd.AddCommand(repositoriesCmd)
}
