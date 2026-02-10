package services

import (
    "bytes"
    "encoding/json"
    "fmt"
    "os"

    "github.com/cli/go-gh/v2/pkg/api"
    "github.com/messagedigest-net/gh-advanced-security/model"
)

// UploadSarif uploads a SARIF file to the specified repository using the
// Code Scanning SARIF upload endpoint.
func UploadSarif(owner, repo, filePath, commitSHA, ref, checkoutURI, toolName string, validate bool, jsonOutput bool) error {
    // Read SARIF file
    data, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }

    payload := model.UploadSarif{
        CommitSHA:   commitSHA,
        Ref:         ref,
        Sarif:       string(data),
        CheckoutUri: checkoutURI,
        ToolName:    toolName,
        Validate:    validate,
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    path := fmt.Sprintf("repos/%s/%s/code-scanning/sarifs", owner, repo)

    resp, err := client.Request("POST", path, bytes.NewReader(body))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    success := resp.StatusCode >= 200 && resp.StatusCode < 300
    if !success {
        return api.HandleHTTPError(resp)
    }

    if jsonOutput {
        var info model.SarifUploadInformation
        if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
            return err
        }
        return jsonLister(info)
    }

    fmt.Println("SARIF uploaded. Processing may be asynchronous.")
    return nil
}
