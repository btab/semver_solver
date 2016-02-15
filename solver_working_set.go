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

func (ws *workingSet) apply(name string, constraint Constraint) (a *Artifact, newPick bool) {
	compiledRange := constraint.Range
	for _, c := range ws.constraints[name] {
		compiledRange = compiledRange.AND(c.Range)
	}
	ws.constraints[name] = append(ws.constraints[name], constraint)

	artifacts, cacheMiss := ws.ensureCache(name)

	for i, artifact := range artifacts {
		if compiledRange(artifact.Version) {
			if !cacheMiss && i == 0 {
				return &artifact, false
			} else {
				ws.artifactsByName[name] = artifacts[i:]
				return &artifact, true
			}
		}
	}

	return nil, false
}

func (ws *workingSet) ensureCache(name string) (as []Artifact, cacheMiss bool) {
	artifacts, ok := ws.artifactsByName[name]

	if ok {
		return artifacts, false
	}

	artifacts = ws.source.AllVersionsOf(name)

	// copy for isolation
	localCopy := make([]Artifact, len(artifacts))
	copy(localCopy, artifacts)

	// reverse sort to make all the 'default to latest' optimizations work
	sort.Sort(sort.Reverse(SortableArtifacts(localCopy)))

	ws.artifactsByName[name] = localCopy
	return localCopy, true
}

func (ws *workingSet) picks() []Artifact {
	var result []Artifact

	for _, artifacts := range ws.artifactsByName {
		result = append(result, artifacts[0])
	}

	return result
}
