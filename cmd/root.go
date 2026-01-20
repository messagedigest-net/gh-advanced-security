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
  Long:  `GitHub CLI extension to manage GitHub Advanced Security features...`,
  Run: func(cmd *cobra.Command, args []string) {
    services.ChooseSubCommand(cmd.Commands(), args, "What do you want to do?")
  },
}

// Global prompter variable can stay if used widely,
// but the flag variables (user, all, json, etc) are GONE.
var prompt *prompter.Prompter

func init() {
  prompt = services.GetPrompt()

  // Delegate Flag Definition to the Service
  services.DefineGlobalFlags(rootCmd)
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
