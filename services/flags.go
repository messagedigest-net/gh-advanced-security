package services

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GlobalFlags holds the values for the persistent flags
type GlobalFlags struct {
	JSON     bool
	User     bool
	All      bool
	PageSize int
}

var flags GlobalFlags

// DefineGlobalFlags registers the flags on the root command.
// Call this from cmd/root.go init().
func DefineGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&flags.JSON, "json", "j", false, "Output in JSON")
	cmd.PersistentFlags().BoolVarP(&flags.User, "user", "u", false, "Show user data instead of organization")
	cmd.PersistentFlags().BoolVarP(&flags.All, "all", "a", false, "Get all data for paged API responses (no pause)")
	cmd.PersistentFlags().IntVarP(&flags.PageSize, "page", "p", 0, "Number of lines to show per page (default: terminal height)")

	viper.BindPFlag("json", cmd.PersistentFlags().Lookup("json"))
	viper.BindPFlag("user", cmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("all", cmd.PersistentFlags().Lookup("all"))
	viper.BindPFlag("page", cmd.PersistentFlags().Lookup("page"))
}

// ParseGlobalFlags extracts the values from the command context.
// Call this from GetTarget or any command that needs the context.
func ParseGlobalFlags(cmd *cobra.Command, prompt []string) (*GlobalFlags, error) {
	if len(prompt) > 1 {
		err := cmd.ParseFlags(prompt)
		if err != nil {
			return nil, err
		}
	}
	return &flags, nil
}

func GetGlobalFlags() *GlobalFlags {
	return &flags
}
