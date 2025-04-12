// Copyright 2025 Francisco Oliveto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package insights

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// setup sets up a test HTTP server along with a insights.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup(t *testing.T) (client *Client, mux *http.ServeMux) {
	t.Helper()
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	apiMux := http.NewServeMux()
	apiMux.Handle("/v3/", http.StripPrefix("/v3", mux))

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiMux)

	// client is the deps.dev client being tested;
	// it is configured to use test server.
	client = NewClient()
	client.BaseURL, _ = url.Parse(server.URL + "/v3/")

	t.Cleanup(server.Close)

	return client, mux
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testQueryParameter(t *testing.T, r *http.Request, key string, want string) {
	t.Helper()
	if got := r.FormValue(key); got != want {
		t.Errorf("FormValue(%q) returned %q, want %q", key, got, want)
	}
}

func TestGetPackage(t *testing.T) {
	// TODO: should this test run in parallel?
	client, mux := setup(t)

	mux.HandleFunc("/systems/go/packages/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", "application/json; charset=utf-8")
		fmt.Fprint(w, `{"packageKey":{"system":"GO","name":"foo"}}`)
	})

	want := &Package{
		PackageKey: PackageKey{System: "GO", Name: "foo"},
	}

	got, err := client.GetPackage(context.Background(), "go", "foo")
	if err != nil {
		t.Errorf("GetPackage failed: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Errorf("GetPackage returned %+v; want %+v", got, want)
	}
}

func TestGetPackageErrorNotFound(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("/systems/bar/packages/baz", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "package not found", http.StatusNotFound)
	})

	_, err := client.GetPackage(context.Background(), "bar", "baz")
	if err == nil {
		t.Errorf("GetPackage expected error")
	}
}

func TestGetVersion(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("/systems/go/packages/rsc.io%2Fgithub/versions/v0.4.1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"versionKey":{"system":"GO","name":"rsc.io/github","version":"v0.4.1"}}`)
	})

	want := &Version{
		VersionKey: VersionKey{
			System:  "GO",
			Name:    "rsc.io/github",
			Version: "v0.4.1",
		},
	}

	got, err := client.GetVersion(context.Background(), "go", "rsc.io/github", "v0.4.1")
	if err != nil {
		t.Errorf("GetVersion failed: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Errorf("GetVersion returned %+v; want %+v", got, want)
	}
}

func TestGetDependencies(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/systems/npm/packages/react/versions/18.2.0:dependencies", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"nodes":[{"versionKey":{"system":"NPM", "name":"react", "version":"18.2.0"}, "bundled":false, "relation":"SELF", "errors":[]}]}`)
	})

	want := &Dependencies{
		Nodes: []Node{
			{
				VersionKey: VersionKey{
					System:  "NPM",
					Name:    "react",
					Version: "18.2.0",
				},
				Bundled:  false,
				Relation: "SELF",
				Errors:   []string{},
			},
		},
	}

	got, err := client.GetDependencies(context.Background(), "npm", "react", "18.2.0")
	if err != nil {
		t.Errorf("GetDependencies failed: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Errorf("c.GetDependencies returned %+v; want %+v", got, want)
	}
}

func TestGetProject(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/projects/github.com%2Frobpike%2Flisp", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"projectKey":{"id":"github.com/robpike/lisp"}, "openIssuesCount":0, "starsCount":978}`)
	})

	want := &Project{
		ProjectKey:      ProjectKey{ID: "github.com/robpike/lisp"},
		OpenIssuesCount: 0,
		StarsCount:      978,
	}

	got, err := client.GetProject(context.Background(), "github.com/robpike/lisp")
	if err != nil {
		t.Errorf("GetProject failed: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Errorf("GetProject returned %+v; want %+v", got, want)
	}
}

func TestGetProjectPackageVersions(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/projects/github.com%2Frobpike%2Flisp:packageversions", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprint(w, `{"versions":[{"versionKey":{"system":"GO", "name":"robpike.io/lisp", "version":"v0.0.0"}, "relationType":"SOURCE_REPO", "relationProvenance":"GO_ORIGIN"}]}`)
		fmt.Fprint(w, `{}`)
	})

	// TODO: add values.
	want := &ProjectPackageVersions{}

	got, err := client.GetProjectPackageVersions(context.Background(), "github.com/robpike/lisp")
	if err != nil {
		t.Errorf("GetProjectPackageVersions failed: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Errorf("GetProjectPackageVersions returned %+v; want %+v", got, want)
	}
}

func TestGetAdvisory(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/advisories/GHSA-2qrg-x229-3v8q", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"advisoryKey":{"id":"GHSA-2qrg-x229-3v8q"}}`)
	})

	want := &Advisory{AdvisoryKey: AdvisoryKey{ID: "GHSA-2qrg-x229-3v8q"}}

	got, err := client.GetAdvisory(context.Background(), "GHSA-2qrg-x229-3v8q")
	if err != nil {
		t.Errorf("GetAdvisory failed: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Errorf("GetAdvisory returned %+v; want %+v", got, want)
	}
}

func TestQuery(t *testing.T) {
	client, mux := setup(t)

	mux.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", "application/json; charset=utf-8")
		testQueryParameter(t, r, "versionKey.system", "npm")
		testQueryParameter(t, r, "versionKey.name", "react")
		testQueryParameter(t, r, "versionKey.version", "18.2.0")
		fmt.Fprint(w, `{"results":[{"version":{"versionKey":{"system":"NPM", "name":"react", "version":"18.2.0"}}}]}`)
	})

	want := &QueryResult{
		Results: []Result{{
			Version{VersionKey: VersionKey{
				System:  "NPM",
				Name:    "react",
				Version: "18.2.0",
			}},
		}},
	}

	opts := &QueryOptions{System: "npm", Name: "react", Version: "18.2.0"}
	got, err := client.Query(context.Background(), opts)
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Query returned %+v; want %+v", got, want)
	}
}
