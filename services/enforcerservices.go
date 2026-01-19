package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/messagedigest-net/gh-advanced-security/model"
)

var enforcerSvcs *EnforcerServices

type EnforcerServices struct{}

func GetEnforcerServices() *EnforcerServices {
	if enforcerSvcs == nil {
		enforcerSvcs = &EnforcerServices{}
	}
	return enforcerSvcs
}

// EnablePushProtection checks and enables feature for a single repo
func (e *EnforcerServices) EnablePushProtection(owner, repo string) error {
	path := fmt.Sprintf("repos/%s/%s", owner, repo)

	// Construct payload: Secret Scanning MUST be enabled to enable Push Protection
	payload := model.RepoUpdateRequest{
		SecurityAndAnalysis: &model.SecurityAndAnalysisReq{
			SecretScanning: &model.StatusReq{Status: "enabled"},
			PushProtection: &model.StatusReq{Status: "enabled"},
		},
	}

	// We use client.Request directly for PATCH as Go-GH might not have a helper for this specific body
	// Note: You'll need to update services/core.go to expose a Patch helper or access client
	// For now, assuming you add a Patch method to core.go similar to Get

	// Simulated Patch call (since we need to add Patch to core.go):
	// return client.Patch(path, payload, nil)
	// See "Missing Piece" below for the core.go update
	return patch(path, payload)
}

// BulkEnablePushProtection enables it for ALL repos in an org
// Uses a Worker Pool to limit concurrency to 5 requests at a time (Safety against Rate Limits)
func (e *EnforcerServices) BulkEnablePushProtection(org string) error {
	repoSvc := GetRepositoryServices()

	fmt.Printf("Fetching all repositories for %s... (this may take a moment)\n", org)

	repos, err := repoSvc.FetchAllForOrg(org)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d repositories. Starting enforcement...\n", len(repos))

	// Worker Pool Setup
	jobs := make(chan model.Repository, len(repos))
	results := make(chan string, len(repos))
	var wg sync.WaitGroup

	// 5 Workers
	concurrency := 5
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for repo := range jobs {
				// Skip if already enabled to save API tokens
				if repo.SecurityAndAnalysis.SecretScanningPushProtection.Status == "enabled" {
					results <- fmt.Sprintf("Add: %s (Already Enabled)", repo.Name)
					continue
				}

				err := e.EnablePushProtection(repo.Owner.Login, repo.Name)
				if err != nil {
					results <- fmt.Sprintf("Fail: %s (%s)", repo.Name, err.Error())
				} else {
					results <- fmt.Sprintf("Success: %s", repo.Name)
				}
				// Rate limit niceness
				time.Sleep(200 * time.Millisecond)
			}
		}()
	}

	// Send jobs
	for _, repo := range repos {
		jobs <- repo
	}
	close(jobs)

	// Wait for workers in a separate goroutine to close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Print results as they come in
	successCount := 0
	failCount := 0
	for res := range results {
		if strings.HasPrefix(res, "Success") {
			successCount++
			fmt.Println("\033[32m" + res + "\033[0m") // Green
		} else if strings.HasPrefix(res, "Fail") {
			failCount++
			fmt.Println("\033[31m" + res + "\033[0m") // Red
		} else {
			// Already enabled
			fmt.Println(res)
		}
	}

	fmt.Printf("\nDone. Enabled: %d, Failed: %d\n", successCount, failCount)
	return nil
}

// EnableSecretScanning enables ONLY Secret Scanning (without Push Protection)
func (e *EnforcerServices) EnableSecretScanning(owner, repo string) error {
	path := fmt.Sprintf("repos/%s/%s", owner, repo)

	payload := model.RepoUpdateRequest{
		SecurityAndAnalysis: &model.SecurityAndAnalysisReq{
			SecretScanning: &model.StatusReq{Status: "enabled"},
			// We intentionally omit PushProtection here
		},
	}

	return patch(path, payload)
}

// BulkEnableSecretScanning enables it for ALL repos in an org
func (e *EnforcerServices) BulkEnableSecretScanning(org string) error {
	repoSvc := GetRepositoryServices()

	fmt.Printf("Fetching repositories for %s...\n", org)
	repos, err := repoSvc.FetchAllForOrg(org)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d repositories. Starting enforcement...\n", len(repos))

	jobs := make(chan model.Repository, len(repos))
	results := make(chan string, len(repos))
	var wg sync.WaitGroup

	// Worker Pool (Concurrency: 5)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for repo := range jobs {
				// Skip if already enabled
				if repo.SecurityAndAnalysis.SecretScanning.Status == "enabled" {
					results <- fmt.Sprintf("Add: %s (Already Enabled)", repo.Name)
					continue
				}

				err := e.EnableSecretScanning(repo.Owner.Login, repo.Name)
				if err != nil {
					results <- fmt.Sprintf("Fail: %s (%s)", repo.Name, err.Error())
				} else {
					results <- fmt.Sprintf("Success: %s", repo.Name)
				}
				time.Sleep(200 * time.Millisecond) // Rate limit safety
			}
		}()
	}

	for _, repo := range repos {
		jobs <- repo
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	// Output Loop
	successCount := 0
	failCount := 0
	for res := range results {
		if strings.HasPrefix(res, "Success") {
			successCount++
			fmt.Println("\033[32m" + res + "\033[0m")
		} else if strings.HasPrefix(res, "Fail") {
			failCount++
			fmt.Println("\033[31m" + res + "\033[0m")
		} else {
			fmt.Println(res)
		}
	}

	fmt.Printf("\nDone. Enabled: %d, Failed: %d\n", successCount, failCount)
	return nil
}
