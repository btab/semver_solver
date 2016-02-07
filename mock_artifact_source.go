package semver_solver

import (
	"log"

	"github.com/blang/semver"
)

type MockArtifactSource map[string]map[string]Artifact

func (source MockArtifactSource) AddArtifact(name string, version semver.Version) {
	source.AddArtifactWithDeps(name, version, ConstraintSet{})
}

func (source MockArtifactSource) AddArtifactWithDeps(name string, version semver.Version, deps ConstraintSet) {
	verString := version.String()

	versions, ok := source[name]

	if !ok {
		versions = map[string]Artifact{}
		source[name] = versions
	}

	if _, ok = versions[verString]; ok {
		log.Fatalf("already added %s @ %s to mock artifact source", name, verString)
	}

	versions[verString] = Artifact{
		name:      name,
		version:   version,
		dependsOn: deps,
	}
}

func (source MockArtifactSource) AllVersionsOf(name string) []Artifact {
	var result []Artifact

	for _, artifact := range source[name] {
		result = append(result, artifact)
	}

	return result
}
