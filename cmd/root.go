package cmd

import (
  "fmt"
  "os"

  "github.com/cli/go-gh/v2/pkg/prompter"
  "github.com/messagedigest-net/gh-advanced-security/services"
  "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "advanced-security",
  Short: "GitHub CLI extension to manage GitHub Advanced Security features",
  Long: `GitHub CLI extension to manage GitHub Advanced Security features:
                DependaBot
                Code Scanning
                Secret Scanning
                Security Events`,
  Run: func(cmd *cobra.Command, args []string) {
    services.ChooseSubCommand(cmd.Commands(), args, "What do you want to do?")
  },
}

var (
  prompt       *prompter.Prompter
  json         bool
  user         bool
  all          bool
  UserPageSize int
)

func init() {
  prompt = services.GetPrompt()
  rootCmd.PersistentFlags().BoolVarP(&json, "json", "j", false, "Output in JSON")
  rootCmd.PersistentFlags().BoolVarP(&user, "user", "u", false, "Show user data instead of organization. e.g.: gh advanced-security list repos -u username")
  rootCmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "Get all data for paged API responses.")
  rootCmd.PersistentFlags().IntVarP(&UserPageSize, "page", "p", 0, "Number of lines to show per page (default: terminal height, max: 100)")

}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
