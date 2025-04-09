package insight

import (
	"context"
	"fmt"
	"net/url"
)

type Dependency struct {
	// The name of the package.
	Name string

	// The requirement of the package.
	Requirement string
}

type DependencyGroup struct {
	// The target framework that this dependency group is for.
	TargetFramework string

	// The requirements belonging to this dependency group.
	Dependencies []Dependency
}

type NuGet struct {
	// The requirements grouped by target framework.
	DependencyGroups []DependencyGroup
}

type NPMDependencies struct {
	// The "dependencies" field of a package.json, represented as a list of
	// name, requirement pairs.
	Dependencies []Dependency

	// The "devDependencies" field of a package.json. The format is the
	// same as "dependencies".
	DevDependencies []Dependency

	// The "optionalDependencies" field of a package.json. The format is
	// the same as "dependencies".
	OptionalDependencies []Dependency

	// The "peerDependencies" field of a package.json. The format is the
	// same as "dependencies".
	PeerDependencies []Dependency

	// The "bundleDependencies" field of a package.json: a list of package
	// names. In the package.json this may also just be the boolean value
	// "true", in which case this field will contain the names of all the
	// dependencies from the "dependencies" field.
	BundleDependencies []string
}

type Bundle struct {
	// The path inside the tarball where this dependency was found.
	Path string

	// The name of the bundled package, as declared inside the bundled
	// package.json.
	Name string

	// The version of this package, as declared inside the bundled
	// package.json.
	Version string

	// The dependency-related fields from the bundled package.json.
	Dependencies NPMDependencies
}

type NPM struct {
	// The dependency-related fields declared in the requested package version's
	// package.json.
	Dependencies NPMDependencies

	// Contents of any additional package.json files found inside the
	// "node_modules" folder of the version's tarball, including nested
	// "node_modules".
	Bundled []Bundle
}

type MavenDependency struct {
	// The name of the package.
	Name string

	// The version requirement of the dependency.
	Version string

	// The classifier of the dependency, which distinguishes artifacts that
	// differ in content.
	Classifier string

	// The type of the dependency, defaults to jar.
	Type string

	// The scope of the dependency, specifies how to limit the transitivity
	// of a dependency.
	Scope string

	// Whether the dependency is optional or not.
	Optional string

	// The dependencies to be excluded, in the form of a list of package
	// names.
	// Exclusions may contain wildcards in both groupID and artifactID.
	Exclusions []string
}

type Property struct {
	// The name of the property.
	Name string

	// The value of the property.
	Value string
}

type Repository struct {
	// The ID of the repository.
	ID string

	// The URL of the repository.
	URL string

	// Whether the description of the repository follows a common layout.
	Layout string

	// Whether the repository is enabled for release downloads.
	ReleasesEnabled string

	// Whether the repository is enabled for snapshot downloads.
	SnapshotsEnabled string
}

type JDK struct {
	// The JDK requirement to activate the profile.
	JDK string
}

type OS struct {
	// The name of the operating system.
	Name string

	// The family of the operating system.
	Family string

	// The CPU architecture of the operating system.
	Arch string

	// The version of the operating system.
	Version string
}

type File struct {
	// The name of the file that its existence activates the profile.
	Exists string

	// The name of the file, activate the profile if the file is missing.
	Missing string
}

type Activation struct {
	// Whether the profile is active by default.
	ActiveByDefault string

	// The JDK requirement of the activation.
	JDK JDK

	// The operating system requirement of the activation.
	OS OS

	// The property requirement of the activation.
	Property struct {
		// The property requirement to activate the profile.
		// This can be a system property or CLI user property.
		Property Property
	}

	// The file requirement of the activation.
	File File
}

type Profile struct {
	// The ID of the profile.
	ID string

	// The activation requirement of the profile.
	Activation Activation

	// The dependencies specified in the profile.
	Dependencies []MavenDependency

	// The dependency management specified in the profile.
	DependencyManagement []MavenDependency

	// The properties specified in the profile.
	Properties []Property

	// The repositories specified in the profile.
	Repositories []Repository
}

type Maven struct {
	// The direct parent of a package version.
	Parent VersionKey

	// The list of dependencies.
	Dependencies []MavenDependency

	// The list of dependency management.
	// The format is the same as dependencies.
	DependencyManagement []MavenDependency

	// The list of properties, used to resolve placeholders.
	Properties []Property

	// The list of repositories.
	Repositories []Repository

	// The list of profiles.
	Profiles []Profile
}

// Requirements contains a system-specific representation of the requirements
// specified by a package version. Only one of its fields will be set.
type Requirements struct {
	// The NuGet-specific representation of the version's requirements.
	//
	// Note that the term "dependency" is used here to mean "a single unresolved
	// requirement" to be consistent with how the term is used in the NuGet
	// ecosystem. This is different to how it is used elsewhere in the deps.dev
	// API.
	NuGet NuGet

	// The npm-specific representation of the version's requirements.
	//
	// Note that the term "dependency" is used here to mean "a single unresolved
	// requirement" to be consistent with how the term is used in the npm
	// ecosystem. This is different to how it is used elsewhere in the deps.dev
	// API.
	NPM NPM

	// The Maven-specific representation of the version's requirements.
	//
	// Note that the term "dependency" is used here to mean "a single unresolved
	// requirement" to be consistent with how the term is used in the Maven
	// ecosystem. This is different to how it is used elsewhere in the deps.dev
	// API.
	//
	// This data is as it is declared in a version POM file. The data in parent
	// POMs are not merged.
	// Any string field may contain references to properties, and the properties
	// are not interpolated.
	Maven Maven
}

// GetRequirements returns the requirements for a given version in a system-specific format.
//
// deps.dev API doc: https://docs.deps.dev/api/v3/#getrequirements
func (c *Client) GetRequirements(ctx context.Context, system, name, version string) (*Requirements, error) {
	path := fmt.Sprintf("/systems/%s/packages/%s/versions/%s:requirements", url.PathEscape(system), url.PathEscape(name), url.PathEscape(version))
	r := new(Requirements)
	if err := c.get(ctx, path, r); err != nil {
		return nil, err
	}
	return r, nil
}
