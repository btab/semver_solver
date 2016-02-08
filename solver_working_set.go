package semver_solver

import "sort"

type workingSet struct {
	source          ArtifactSource
	artifactsByName map[string][]Artifact
	constraints     ConstraintSet
}

func newWorkingSet(source ArtifactSource) *workingSet {
	return &workingSet{
		source:          source,
		artifactsByName: map[string][]Artifact{},
		constraints:     ConstraintSet{},
	}
}

func (ws *workingSet) apply(name string, constraint Constraint) *Artifact {
	compiledRange := constraint.Range
	for _, c := range ws.constraints[name] {
		compiledRange = compiledRange.AND(c.Range)
	}
	ws.constraints[name] = append(ws.constraints[name], constraint)

	artifacts := ws.ensureCache(name)
	for i, artifact := range artifacts {
		if compiledRange(artifact.Version) {
			ws.artifactsByName[name] = artifacts[i:]
			return &artifact
		}
	}

	return nil
}

func (ws *workingSet) ensureCache(name string) []Artifact {
	artifacts, ok := ws.artifactsByName[name]

	if ok {
		return artifacts
	}

	artifacts = ws.source.AllVersionsOf(name)

	// copy for isolation
	localCopy := make([]Artifact, len(artifacts))
	copy(localCopy, artifacts)

	// reverse sort to make all the 'default to latest' optimizations work
	sort.Sort(sort.Reverse(SortableArtifacts(localCopy)))

	ws.artifactsByName[name] = localCopy
	return localCopy
}

func (ws *workingSet) picks() []Artifact {
	var result []Artifact

	for _, artifacts := range ws.artifactsByName {
		result = append(result, artifacts[0])
	}

	return result
}
