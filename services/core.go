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

func jsonLister(object interface{}) error {
	jsonObject, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		return err
	}

	jsonpretty.Format(terminal.Out(), bytes.NewReader(jsonObject), "\t", terminal.IsColorEnabled())
	//fmt.Println(string(jsonObject))
	return nil
}

func getPages(path string, target interface{}) (next string, err error) {
	resp, err := client.Request("GET", path, nil)
	if err != nil {
		return next, err
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		err = api.HandleHTTPError(resp)
		return next, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return next, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return next, err
	}

	err = json.Unmarshal(b, &target)
	if err != nil {
		return next, err
	}

	linkHeader, ok := resp.Header["Link"]

	if ok {
		links := strings.Split(linkHeader[0], ",")
		for _, v := range links {
			parts := strings.Split(v, ";")
			if strings.Compare(parts[1], " rel=\"next\"") == 0 {
				next = strings.TrimSpace(parts[0])
				next, _ = strings.CutPrefix(next, "<")
				next, _ = strings.CutSuffix(next, ">")

				break
			}
		}
	}

	return next, err
}
