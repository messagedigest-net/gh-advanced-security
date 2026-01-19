package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/jsonpretty"
)

var client *api.RESTClient

func init() {
	initRestClient()
}

func initRestClient() {
	var err error
	client, err = api.DefaultRESTClient()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// jsonLister remains using interface{} as json.Marshal accepts any type
func jsonLister(object interface{}) error {
	jsonObject, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		return err
	}

	// Ensure terminal is initialized (from terminal.go)
	t := GetTerminal()
	jsonpretty.Format(t.Out(), bytes.NewReader(jsonObject), "\t", t.IsColorEnabled())
	return nil
}

// getPages is now Generic [T any].
// T represents the shape of the data you expect (e.g., []model.Repository).
// We pass *T so we can unmarshal directly into it.
func getPages[T any](path string, target *T) (next string, err error) {
	resp, err := client.Request("GET", path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		return "", api.HandleHTTPError(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		return "", nil
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Type-safe unmarshal
	err = json.Unmarshal(b, target)
	if err != nil {
		return "", err
	}

	// Robust Link Header Parsing
	linkHeader, ok := resp.Header["Link"]
	if ok {
		links := strings.Split(linkHeader[0], ",")
		for _, v := range links {
			parts := strings.Split(v, ";")
			if len(parts) > 1 {
				// Trim spaces to handle " rel="next"" vs "rel="next""
				rel := strings.TrimSpace(parts[1])
				if rel == "rel=\"next\"" {
					next = strings.TrimSpace(parts[0])
					next, _ = strings.CutPrefix(next, "<")
					next, _ = strings.CutSuffix(next, ">")
					break
				}
			}
		}
	}

	return next, nil
}
