package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

const (
	availableContent = `<html>
<input type="button" value="Add to Cart" id="00020180109102414710EdNr3gMB067B" class="btn-add m-t-20 btn-addCart" disabled="disabled">
</html>`

	unavailableContent = `<html>
<input type="button" value="Sold Out" id="00020180109102414710EdNr3gMB067B" class="btn-add m-t-20 btn-addCart" disabled="disabled">
</html>`
)

const (
	available   = "available"
	unavailable = "unavailable"
)

var routes = map[string]string{
	available:   availableContent,
	unavailable: unavailableContent,
}

func init() {
	http.HandleFunc("/", handler)
	go func() { http.ListenAndServe(":8080", nil) }()
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path[1:] {
	case available:
		w.Write([]byte(routes[available]))
	case unavailable:
		w.Write([]byte(routes[unavailable]))
	default:
		panic(fmt.Sprintf("bad route '%s'", r.URL.Path))
	}
}

func TestLoadData(t *testing.T) {
	m, err := loadData("test.json")
	if err != nil {
		t.Fatalf("failed to unmarshal data.json: %v", err)
	}

	if len(m.Matches) < 1 {
		t.Fatal("array of matches is empty")
	}

	if len(m.Matches[0].URI) < 1 {
		t.Fatal("URI of first match is empty")
	}

	if len(m.Matches[0].Regex) < 1 {
		t.Fatal("regex of first match has no elements")
	}

	if len(m.Matches[0].Regex[0]) < 1 {
		t.Fatal("first regex of first match is empty")
	}
}

func TestGetPageContent(t *testing.T) {
	m, err := loadData("test.json")
	if err != nil {
		t.Fatalf("failed to unmarshal data.json: %v", err)
	}

	for r := range routes {
		c, err := getPageContent((m.Matches[0].URI) + r)
		if err != nil {
			t.Fatalf("failed to get page context: %v", err)
		}

		if len(c) < 1 {
			t.Fatal("page content is empty")
		}

		if strings.Index(c, "html") == -1 {
			t.Fatal("page content doesn't appear to be HTML")
		}
	}
}

func TestMatches(t *testing.T) {
	m, err := loadData("test.json")
	if err != nil {
		t.Fatalf("failed to unmarshal data.json: %v", err)
	}

	for r := range routes {
		c, err := getPageContent((m.Matches[0].URI) + r)
		if err != nil {
			t.Fatalf("failed to get page context: %v", err)
		}

		match := findMatch(c, m.Matches[0].Regex[0])

		switch r {
		case available:
			if !match {
				t.Fatal("didn't find match when there should be a match")
			}
		case unavailable:
			if match {
				t.Fatal("found match when there shouldn't be a match")
			}
		}
	}
}
