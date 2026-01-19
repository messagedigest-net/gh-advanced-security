package model

type User struct {
	Login             string
	Id                int
	NodeID            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string
	SiteAdmin         bool `json:"site_admin"`
}
