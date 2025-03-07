// Copyright 2025 Francisco Oliveto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package insight

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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
