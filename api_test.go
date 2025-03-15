// Copyright 2025 Francisco Oliveto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package insight

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestGetPackage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"packageKey":{"system":"GO","name":"foo"},"versions":[{"versionKey":{"system":"GO","name":"foo","version":"v0.1.0"},"publishedAt":"2019-07-25T19:01:57Z","isDefault":false},{"versionKey":{"system":"GO","name":"foo","version":"v0.2.0"},"publishedAt":"2019-07-25T19:02:00Z","isDefault":false}]}`)
	}))
	defer ts.Close()

	want := &Package{
		PackageKey: PackageKey{System: "GO", Name: "foo"},
		Versions: []Version{
			{
				VersionKey:  VersionKey{System: "GO", Name: "foo", Version: "v0.1.0"},
				PublishedAt: "2019-07-25T19:01:57Z",
				IsDefault:   false,
			},
			{
				VersionKey:  VersionKey{System: "GO", Name: "foo", Version: "v0.2.0"},
				PublishedAt: "2019-07-25T19:02:00Z",
				IsDefault:   false,
			},
		},
	}

	c := NewClient(ts.URL, nil)
	got, err := c.GetPackage("go", "foo")
	if err != nil {
		t.Errorf("c.GetPackage error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("c.GetPackage() == %v; want %v", got, want)
	}
}

func TestGetPackageError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "package not found", http.StatusNotFound)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, nil)
	_, err := c.GetPackage("bar", "baz")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestGetVersion(t *testing.T) {
	body := `{"versionKey":{"system":"GO","name":"rsc.io/github","version":"v0.4.1"},"publishedAt":"2024-06-21T16:57:04Z","isDefault":false,"licenses":["BSD-3-Clause"],"advisoryKeys":[],"links":[{"label":"SOURCE_REPO","url":"https://github.com/rsc/github"}],"slsaProvenances":[],"attestations":[],"registries":[],"relatedProjects":[{"projectKey":{"id":"github.com/rsc/github"},"relationProvenance":"GO_ORIGIN","relationType":"SOURCE_REPO"}]}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}))
	defer ts.Close()

	want := new(Version)
	r := strings.NewReader(body)
	if err := json.NewDecoder(r).Decode(want); err != nil {
		t.Errorf("%v", err)
	}

	c := NewClient(ts.URL, nil)
	got, err := c.GetVersion("go", "rsc.io/github", "v0.4.1")
	if err != nil {
		t.Errorf("c.GetVersion error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("c.GetVersion() == %v; want %v", got, want)
	}
}
