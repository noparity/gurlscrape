package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
)

type matches struct {
	Matches []match `json:"matches"`
}

type match struct {
	URI   string   `json:"uri"`
	Regex []string `json:"regex"`
}

func main() {

}

func loadData(filename string) (matches, error) {
	m := matches{}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(b, &m)
	return m, err
}

func getPageContent(uri string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func findMatch(content, regex string) bool {
	r := regexp.MustCompile(regex)
	return r.MatchString(content)
}
