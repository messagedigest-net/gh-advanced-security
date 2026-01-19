package model

type DependabotAlert struct {
	Number                int                   `json:"number"`
	State                 string                `json:"state"`
	Dependency            Dependency            `json:"dependency"`
	SecurityAdvisory      SecurityAdvisory      `json:"security_advisory"`
	SecurityVulnerability SecurityVulnerability `json:"security_vulnerability"`
	Url                   string                `json:"url"`
	HtmlUrl               string                `json:"html_url"`
	CreatedAt             string                `json:"created_at"`
	UpdatedAt             string                `json:"updated_at"`
	DismissedAt           string                `json:"dismissed_at"`
	DismissedBy           User                  `json:"dismissed_by"`
	DismissedReason       string                `json:"dismissed_reason"`
	DismissedComment      string                `json:"dismissed_comment"`
}

type Dependency struct {
	Package      Package `json:"package"`
	ManifestPath string  `json:"manifest_path"`
	Scope        string  `json:"scope"`
}

type Package struct {
	Ecosystem string `json:"ecosystem"`
	Name      string `json:"name"`
}

type SecurityAdvisory struct {
	GHSAId      string `json:"ghsa_id"`
	CVEId       string `json:"cve_id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type SecurityVulnerability struct {
	Package                Package `json:"package"`
	Severity               string  `json:"severity"`
	VulnerableVersionRange string  `json:"vulnerable_version_range"`
}
