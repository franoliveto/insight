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

func decode(data string, v any) error {
	r := strings.NewReader(data)
	err := json.NewDecoder(r).Decode(v)
	if err != nil {
		return err
	}
	return nil
}

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

	c := NewClient()
	c.BasePath = ts.URL
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

	c := NewClient()
	c.BasePath = ts.URL
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

	c := NewClient()
	c.BasePath = ts.URL
	got, err := c.GetVersion("go", "rsc.io/github", "v0.4.1")
	if err != nil {
		t.Errorf("c.GetVersion error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("c.GetVersion() == %v; want %v", got, want)
	}
}

func TestGetDependencies(t *testing.T) {
	body := `{"nodes":[{"versionKey":{"system":"NPM", "name":"react", "version":"18.2.0"}, "bundled":false, "relation":"SELF", "errors":[]}, {"versionKey":{"system":"NPM", "name":"js-tokens", "version":"4.0.0"}, "bundled":false, "relation":"INDIRECT", "errors":[]}, {"versionKey":{"system":"NPM", "name":"loose-envify", "version":"1.4.0"}, "bundled":false, "relation":"DIRECT", "errors":[]}], "edges":[{"fromNode":0, "toNode":2, "requirement":"^1.1.0"}, {"fromNode":2, "toNode":1, "requirement":"^3.0.0 || ^4.0.0"}], "error":""}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}))
	defer ts.Close()

	want := new(Dependencies)
	if err := decode(body, want); err != nil {
		t.Errorf("%v", err)
	}

	c := NewClient()
	c.BasePath = ts.URL
	got, err := c.GetDependencies(VersionKey{System: "npm", Name: "react", Version: "18.2.0"})
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("c.GetDependencies() == %v; want %v", got, want)
	}
}

func TestGetProject(t *testing.T) {
	body := `{"projectKey":{"id":"github.com/robpike/lisp"}, "openIssuesCount":0, "starsCount":978, "forksCount":50, "license":"BSD-3-Clause", "description":"Toy Lisp 1.5 interpreter", "homepage":"", "scorecard":{"date":"2025-03-10T00:00:00Z", "repository":{"name":"github.com/robpike/lisp", "commit":"e311180f2a0d4ddb7469b9c303e0b152f2e8c4f6"}, "scorecard":{"version":"v5.1.1-14-g9c2d4b23", "commit":"9c2d4b23e7e488f8747b2ea9305f83a299e49c90"}, "checks":[{"name":"Maintained", "documentation":{"shortDescription":"Determines if the project is \"actively maintained\".", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#maintained"}, "score":0, "reason":"0 commit(s) and 0 issue activity found in the last 90 days -- score normalized to 0", "details":[]}, {"name":"Dangerous-Workflow", "documentation":{"shortDescription":"Determines if the project's GitHub Action workflows avoid dangerous patterns.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#dangerous-workflow"}, "score":-1, "reason":"no workflows found", "details":[]}, {"name":"Pinned-Dependencies", "documentation":{"shortDescription":"Determines if the project has declared and pinned the dependencies of its build process.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#pinned-dependencies"}, "score":-1, "reason":"no dependencies found", "details":[]}, {"name":"Code-Review", "documentation":{"shortDescription":"Determines if the project requires human code review before pull requests (aka merge requests) are merged.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#code-review"}, "score":6, "reason":"Found 6/10 approved changesets -- score normalized to 6", "details":[]}, {"name":"Packaging", "documentation":{"shortDescription":"Determines if the project is published as a package that others can easily download, install, easily update, and uninstall.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#packaging"}, "score":-1, "reason":"packaging workflow not detected", "details":["Warn: no GitHub/GitLab publishing workflow detected."]}, {"name":"Token-Permissions", "documentation":{"shortDescription":"Determines if the project's workflows follow the principle of least privilege.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#token-permissions"}, "score":-1, "reason":"No tokens found", "details":[]}, {"name":"Binary-Artifacts", "documentation":{"shortDescription":"Determines if the project has generated executable (binary) artifacts in the source repository.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#binary-artifacts"}, "score":10, "reason":"no binaries found in the repo", "details":[]}, {"name":"CII-Best-Practices", "documentation":{"shortDescription":"Determines if the project has an OpenSSF (formerly CII) Best Practices Badge.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#cii-best-practices"}, "score":0, "reason":"no effort to earn an OpenSSF best practices badge detected", "details":[]}, {"name":"Vulnerabilities", "documentation":{"shortDescription":"Determines if the project has open, known unfixed vulnerabilities.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#vulnerabilities"}, "score":10, "reason":"0 existing vulnerabilities detected", "details":[]}, {"name":"Security-Policy", "documentation":{"shortDescription":"Determines if the project has published a security policy.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#security-policy"}, "score":0, "reason":"security policy file not detected", "details":["Warn: no security policy file detected", "Warn: no security file to analyze", "Warn: no security file to analyze", "Warn: no security file to analyze"]}, {"name":"License", "documentation":{"shortDescription":"Determines if the project has defined a license.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#license"}, "score":10, "reason":"license file detected", "details":["Info: project has a license file: LICENSE:0", "Info: FSF or OSI recognized license: BSD 3-Clause \"New\" or \"Revised\" License: LICENSE:0"]}, {"name":"Fuzzing", "documentation":{"shortDescription":"Determines if the project uses fuzzing.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#fuzzing"}, "score":0, "reason":"project is not fuzzed", "details":["Warn: no fuzzer integrations found"]}, {"name":"Branch-Protection", "documentation":{"shortDescription":"Determines if the default and release branches are protected with GitHub's branch protection settings.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#branch-protection"}, "score":0, "reason":"branch protection not enabled on development/release branches", "details":["Warn: branch protection not enabled for branch 'master'"]}, {"name":"Signed-Releases", "documentation":{"shortDescription":"Determines if the project cryptographically signs release artifacts.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#signed-releases"}, "score":-1, "reason":"no releases found", "details":[]}, {"name":"SAST", "documentation":{"shortDescription":"Determines if the project uses static code analysis.", "url":"https://github.com/ossf/scorecard/blob/9c2d4b23e7e488f8747b2ea9305f83a299e49c90/docs/checks.md#sast"}, "score":0, "reason":"SAST tool is not run on all commits -- score normalized to 0", "details":["Warn: 0 commits out of 12 are checked with a SAST tool"]}], "overallScore":3.8, "metadata":[]}}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}))
	defer ts.Close()

	want := new(Project)
	if err := decode(body, want); err != nil {
		t.Errorf("%v", err)
	}

	c := NewClient()
	c.BasePath = ts.URL
	got, err := c.GetProject("github.com/robpike/lisp")
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("c.GetProject() == %v; want %v", got, want)
	}
}

func TestGetProjectPackageVersions(t *testing.T) {
	body := `{"versions":[{"versionKey":{"system":"GO", "name":"robpike.io/lisp", "version":"v0.0.0-20241117212301-e311180f2a0d"}, "relationType":"SOURCE_REPO", "relationProvenance":"GO_ORIGIN", "slsaProvenances":[], "attestations":[]}]}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}))
	defer ts.Close()

	c := NewClient()

	want := new(ProjectPackageVersions)
	if err := decode(body, want); err != nil {
		t.Errorf("%v", err)
	}

	got, err := c.GetProjectPackageVersions(ProjectKey{ID: "github.com/robpike/lisp"})
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("c.GetProjectPackageVersions() == %v; want %v", got, want)
	}
}

func TestGetAdvisory(t *testing.T) {
	body := `{"advisoryKey":{"id":"GHSA-2qrg-x229-3v8q"}, "url":"https://osv.dev/vulnerability/GHSA-2qrg-x229-3v8q", "title":"Deserialization of Untrusted Data in Log4j", "aliases":["CVE-2019-17571"], "cvss3Score":9.8, "cvss3Vector":"CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H"}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}))
	defer ts.Close()

	c := NewClient()

	want := new(Advisory)
	if err := decode(body, want); err != nil {
		t.Errorf("%v", err)
	}

	got, err := c.GetAdvisory(AdvisoryKey{ID: "GHSA-2qrg-x229-3v8q"})
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("c.GetAdvisory() == %v; want %v", got, want)
	}
}
