// Copyright 2025 Francisco Oliveto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package insight

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// PackageKey identifies a package by name.
type PackageKey struct {
	// The package management system containing the package.
	System string

	// The name of the package.
	Name string
}

// VersionKey identifies a package by version.
type VersionKey struct {
	// The package management system containing the package.
	System string

	// The name of the package.
	Name string

	// The version of the package.
	Version string
}

// AdvisoryKey identifies a security advisory.
type AdvisoryKey struct {
	// The OSV identifier for the security advisory.
	ID string
}

// ProjectKey identifies a project.
type ProjectKey struct {
	// A project identifier of the form `github.com/user/repo`,
	// `gitlab.com/user/repo`, or `bitbucket.org/user/repo`.
	ID string
}

// Link represents a link declared by or derived from package version metadata,
// to an external web resource such as a homepage or source code repository.
type Link struct {
	// A label describing the resource that the link points to.
	Label string

	// The URL of the link.
	URL string
}

// SLSAProvenance contains provenance information extracted from a SLSA
// provenance statement.
type SLSAProvenance struct {
	// The source code repository used to build the version.
	SourceRepository string

	// The commit of the source code repository the version was built from.
	Commit string

	// The URL of the provenance statement if there is one.
	URL string

	// The Sigstore bundle containing this attestation was verified using the
	// [sigstore-go](https://github.com/sigstore/sigstore-go) library.
	Verified bool
}

// Attestation represents a generic attestation. Fields are populated based
// on 'type'.
type Attestation struct {
	// The type of attestation.
	// One of https://slsa.dev/provenance/v0.2, https://slsa.dev/provenance/v1,
	// https://docs.pypi.org/attestations/publish/v1.
	Type string

	// The URL of the attestation if there is one.
	URL string

	// The attestation has been cryptographically verified by deps.dev.
	// For attestations distributed in a Sigstore bundle, this field indicates
	// the bundle was verified using the
	// [sigstore-go](https://github.com/sigstore/sigstore-go) library.
	Verified bool

	// Only set if type is https://slsa.dev/provenance/v0.2,
	// https://slsa.dev/provenance/v1,
	// https://docs.pypi.org/attestations/publish/v1.
	// The source code repository used to build the version.
	SourceRepository string

	// The commit of the source code repository the version was built from.
	Commit string
}

func (c *Client) get(path string, v any) error {
	url, _ := url.JoinPath(c.BasePath, path)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(data))
	}
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

// Package holds information about a package, including a list of its available
// versions, with the default version marked if known.
type Package struct {
	// The name of the package.
	PackageKey PackageKey

	// The available versions of the package.
	Versions []Version
}

// GetPackage returns information about a package, including a list of its available versions.
func (c *Client) GetPackage(system, name string) (*Package, error) {
	path := "/systems/" + url.PathEscape(system) + "/packages/" + url.PathEscape(name)
	var reply Package
	if err := c.get(path, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

// Version holds information about a package version.
type Version struct {
	// The name of the version.
	VersionKey VersionKey

	// The time when this package version was published, if available, as
	// reported by the package management authority.
	PublishedAt string

	// If true, this is the default version of the package: the version that is
	// installed when no version is specified. The precise meaning of this is
	// system-specific, but it is commonly the version with the greatest
	// version number, ignoring pre-release versions.
	IsDefault bool

	// The licenses governing the use of this package version.
	//
	// We identify licenses as
	// [SPDX 2.1](https://spdx.dev/spdx-specification-21-web-version/)
	// expressions. When there is no associated SPDX identifier, we identify a
	// license as "non-standard". When we are unable to obtain license
	// information, this field is empty. When more than one license is listed,
	// their relationship is unspecified.
	//
	// For Cargo, Maven, npm, NuGet, and PyPI, license information is read from
	// the package metadata. For Go, license information is determined using the
	// [licensecheck](https://github.com/google/licensecheck) package.
	//
	// License information is not intended to be legal advice, and you should
	// independently verify the license or terms of any software for your own
	// needs.
	Licenses []string

	// Security advisories known to affect this package version directly. Further
	// information can be requested using the Advisory method.
	//
	// Note that this field does not include advisories that affect dependencies
	// of this package version.
	AdvisoryKeys []AdvisoryKey

	// Links declared by or derived from package version metadata, to external
	// web resources such as a homepage or source code repository. Note that
	// these links are not verified for correctness.
	Links []Link

	// SLSA provenance information for this package version. Extracted from a
	// SLSA provenance attestation. This is only populated for npm package
	// versions. See the 'attestations' field for more attestations (including
	// SLSA provenance) for all systems.
	SLSAProvenances []SLSAProvenance

	// Attestations for this package version.
	Attestations []Attestation

	// URLs for the package management registries this package version is
	// available from.
	// Only set for systems that use a central repository for package
	// distribution: Cargo, Maven, npm, NuGet, and PyPI.
	Registries []string

	// Projects that are related to this package version.
	RelatedProjects []struct {
		// The identifier for the project.
		ProjectKey ProjectKey

		// How the mapping between project and package version was discovered.
		//
		// Can be one of SLSA_ATTESTATION, GO_ORIGIN, PYPI_PUBLISH_ATTESTATION,
		// UNVERIFIED_METADATA.
		RelationProvenance string

		// What the relationship between the project and the package version is.
		//
		// Can be one of SOURCE_REPO, ISSUE_TRACKER.
		RelationType string
	}
}

// GetVersion returns information about a specific package version, including its
// licenses and any security advisories known to affect it.
func (c *Client) GetVersion(system, name, version string) (*Version, error) {
	path := "/systems/" + url.PathEscape(system) + "/packages/" + url.PathEscape(name) + "/versions/" + version
	var reply Version
	if err := c.get(path, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

// Node represents a node in a resolved dependency graph.
type Node struct {
	// The package version represented by this node. Note that the package and
	// version name may differ from the names in the request, if provided, due
	// to canonicalization.
	//
	// In some systems, a graph may contain multiple nodes for the same package
	// version.
	VersionKey VersionKey

	// If true, this is a bundled dependency.
	//
	// For bundled dependencies, the package name in the version key encodes
	// how the dependency is bundled. As an example, a bundled dependency with
	// a name like "a>1.2.3>b>c" is part of the dependency graph of package "a"
	// at version "1.2.3", and has the local name "c". It may or may not be the
	// same as a package with the global name "c".
	Bundled bool

	// Whether this node represents a direct or indirect dependency within this
	// dependency graph. Note that it's possible for a dependency to be both
	// direct and indirect; if so, it is marked as direct.
	//
	// Can be one of SELF, DIRECT, INDIRECT.
	Relation string

	// Errors associated with this node of the graph, such as an unresolved
	// dependency requirement. An error on a node may imply the graph as a
	// whole is incorrect. These error messages have no defined format and are
	// intended for human consumption.
	Errors []string
}

// Edge represents a directed edge in a resolved dependency graph: a
// dependency relation between two nodes.
type Edge struct {
	// The node declaring the dependency, specified as an index into the list of
	// nodes.
	FromNode int

	// The node resolving the dependency, specified as an index into the list of
	// nodes.
	ToNode int

	// The requirement resolved by this edge, as declared by the "from" node.
	// The meaning of this field is system-specific. As an example, in npm, the
	// requirement "^1.0.0" may be resolved by the version "1.2.3".
	Requirement string
}

// Dependencies holds a resolved dependency graph for a package version.
//
// The dependency graph should be similar to one produced by installing the
// package version on a generic 64-bit Linux system, with no other dependencies
// present. The precise meaning of this varies from system to system.
type Dependencies struct {
	// The nodes of the dependency graph. The first node is the root of the graph.
	Nodes []Node

	// The edges of the dependency graph.
	Edges []Edge

	// Any error associated with the dependency graph that is not specific to a
	// node. An error here may imply the graph as a whole is incorrect.
	// This error message has no defined format and is intended for human
	// consumption.
	Error string
}

// GetDependencies returns a resolved dependency graph for the given package version.
func (c *Client) GetDependencies(key VersionKey) (*Dependencies, error) {
	path := "/systems/" + url.PathEscape(key.System) + "/packages/" + url.PathEscape(key.Name) + "/versions/" + url.PathEscape(key.Version) + ":dependencies"
	var reply Dependencies
	if err := c.get(path, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

// Project holds information about a project hosted by GitHub, GitLab, or
// Bitbucket.
type Project struct {
	// The identifier for the project.
	ProjectKey ProjectKey

	// The number of open issues reported by the project host.
	// Only available for GitHub and GitLab.
	OpenIssuesCount int

	// The number of stars reported by the project host.
	// Only available for GitHub and GitLab.
	StarsCount int

	//The number of forks reported by the project host.
	//Only available for GitHub and GitLab.
	ForksCount int

	// The license reported by the project host.
	License string

	// The description reported by the project host
	Description string

	// The homepage reported by the project host.
	Homepage string

	// An [OpenSSF Scorecard](https://github.com/ossf/scorecard) for the project,
	// if one is available.
	Scorecard Scorecard

	// Details of this project's testing by the
	// [OSS-Fuzz service](https://google.github.io/oss-fuzz/).
	// Only set if the project is tested by OSS-Fuzz.
	OSSFuzz OSSFuzzDetails
}

type Scorecard struct {
	// The date at which the scorecard was produced.
	// The time portion of this field is midnight UTC.
	Date string

	// The source code repository and commit the scorecard was produced from.
	Repository struct {
		// The source code repository the scorecard was produced from.
		Name string

		// The source code commit the scorecard was produced from.
		Commit string
	}

	// The version and commit of the Scorecard program used to produce the
	// scorecard.
	Scorecard struct {
		// The version of the Scorecard program used to produce the scorecard.
		Version string

		// The commit of the Scorecard program used to produce the scorecard.
		Commit string
	}

	// The results of the
	// [Scorecard Checks](https://github.com/ossf/scorecard#scorecard-checks)
	// performed on the project.
	Checks []struct {
		// The name of the check.
		Name string

		// Human-readable documentation for the check.
		Documentation struct {
			// A short description of the check.
			ShortDescription string

			// A link to more details about the check.
			URL string
		}

		// A score in the range [0,10]. A higher score is better.
		// A negative score indicates that the check did not run successfully.
		Score int

		// The reason for the score.
		Reason string

		// Further details regarding the check.
		Details []string
	}

	// A weighted average score in the range [0,10]. A higher score is better.
	OverallScore float64

	// Additional metadata associated with the scorecard.
	Metadata []string
}

type OSSFuzzDetails struct {
	// The total number of lines of code in the project.
	LineCount int

	// The number of lines of code covered by fuzzing.
	LineCoverCount int

	// The date the fuzz test that produced the coverage information was run
	// against this project.
	// The time portion of this field is midnight UTC.
	Date string

	// The URL containing the configuration for the project in the
	// OSS-Fuzz repository.
	ConfigURL string
}

// GetProject returns information about projects hosted by GitHub, GitLab,
// or BitBucket.
func (c *Client) GetProject(id string) (*Project, error) {
	path := "/projects/" + url.PathEscape(id)
	var reply Project
	if err := c.get(path, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

type ProjectPackageVersions struct {
	// The versions that were built from the source code contained in this project.
	Versions []struct {
		// The identifier for the version.
		VersionKey VersionKey
		// The SLSA provenance statements that link the version to the project. This
		// is only populated for npm package versions. See the 'attestations' field
		// for more attestations (including SLSA provenance) for all systems.
		SLSAProvenances []SLSAProvenance
		// Attestations that link the version to the project.
		Attestation []Attestation
		// What the relationship between the project and the package version is.
		// Can be one of SOURCE_REPO, ISSUE_TRACKER.
		RelationType string
		// How the mapping between project and package version was discovered.
		// Can be one of SLSA_ATTESTATION, GO_ORIGIN, PYPI_PUBLISH_ATTESTATION,
		// UNVERIFIED_METADATA.
		RelationProvenance string
	}
}

// GetProjectPackageVersions returns known mappings between the requested project
// and package versions.
// At most 1500 package versions are returned. Mappings which were derived from
// attestations are served first.
func (c *Client) GetProjectPackageVersions(key ProjectKey) (*ProjectPackageVersions, error) {
	path := "/projects/" + url.PathEscape(key.ID) + ":packageversions"
	var reply ProjectPackageVersions
	if err := c.get(path, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}
