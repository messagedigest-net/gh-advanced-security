# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Changed
- Standardized Go source file names in `model/` and `services/` to use idiomatic snake_case (e.g., `delete_analysis.go`, `dependabotalert.go`).
- Renamed test helper files to standard conventions ending in `_test.go` (e.g., `reposervices_helpers_test.go`, `terminal_helpers_test.go`).
- Extensive refactoring and cleanup in core commands (`cmd/`) and services.
- Updated line endings and general formatting across models and services.

### Removed
- Removed the compiled `advanced-security` binary from tracking.
