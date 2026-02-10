package cmd

import (
    "fmt"
    "os"

    "github.com/messagedigest-net/gh-advanced-security/services"
    "github.com/spf13/cobra"
)

var uploadSarifCmd = &cobra.Command{
    Use:   "upload-sarif",
    Short: "Upload a SARIF file to a repository",
    Long:  `Upload a SARIF file to GitHub Code Scanning for a repository.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Resolve target repo
        target, flags := services.GetTarget(cmd, args, "Which repository? (format: owner/repo)")
        owner, repo := parseRepo(target)

        // Determine SARIF file path
        var sarifPath string
        if len(args) >= 2 {
            sarifPath = args[1]
        } else {
            // Prompt
            prompt := services.GetPrompt()
            p, err := prompt.Input("Path to SARIF file:", "")
            if err != nil || len(p) == 0 {
                fmt.Println("SARIF file path required")
                os.Exit(1)
            }
            sarifPath = p
        }

        // Call upload service
        commit, _ := cmd.Flags().GetString("commit")
        ref, _ := cmd.Flags().GetString("ref")
        checkout, _ := cmd.Flags().GetString("checkout-uri")
        tool, _ := cmd.Flags().GetString("tool")
        validate, _ := cmd.Flags().GetBool("validate")

        if err := services.UploadSarif(owner, repo, sarifPath, commit, ref, checkout, tool, validate, flags.JSON); err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
    },
}

func init() {
    uploadSarifCmd.Flags().StringP("commit", "c", "", "Commit SHA for the SARIF (optional)")
    uploadSarifCmd.Flags().StringP("ref", "r", "", "Git ref (e.g., refs/heads/main)")
    uploadSarifCmd.Flags().String("checkout-uri", "", "Checkout URI for SARIF file")
    uploadSarifCmd.Flags().StringP("tool", "t", "", "Tool name to associate with the SARIF")
    uploadSarifCmd.Flags().BoolP("validate", "V", false, "Validate SARIF without creating an analysis")

    rootCmd.AddCommand(uploadSarifCmd)
}

// parseRepo is available in another cmd file; reused here.
