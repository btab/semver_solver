package semver_solver

import (
	"log"

	"github.com/blang/semver"
)

type MockArtifactSource map[string]map[string]Artifact

func (source MockArtifactSource) AddArtifact(name, versionString string) error {
	return source.AddArtifactWithDeps(name, versionString, nil)
}

func (source MockArtifactSource) AddArtifactWithDeps(name, versionString string, deps []*Constraint) error {
	versions, ok := source[name]

	if !ok {
		versions = map[string]Artifact{}
		source[name] = versions
	}

	if _, ok = versions[versionString]; ok {
		log.Fatalf("already added %s @ %s to mock artifact source", name, versionString)
	}

	version, err := semver.Parse(versionString)

	if err != nil {
		return err
	}

	versions[versionString] = Artifact{
		Name:      name,
		Version:   version,
		DependsOn: deps,
	}

	return nil
}

func (source MockArtifactSource) AllVersionsOf(name string) []Artifact {
	var result []Artifact

	for _, artifact := range source[name] {
		result = append(result, artifact)
	}

	return result
}
