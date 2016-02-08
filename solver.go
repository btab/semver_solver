package semver_solver

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// This would become the 'simple solver' if we ever added other solvers (SAT etc)

type Solver struct {
	Source ArtifactSource
}

func (s *Solver) Solve(cs ConstraintSet) ([]Artifact, error) {
	ws := WorkingSet{
		source:          s.Source,
		artifactsByName: map[string][]Artifact{},
		constraints:     ConstraintSet{},
	}

	var allFailures []string

	for len(cs) > 0 && len(allFailures) == 0 {
		var artifactsPicked []*Artifact

		for name, constraints := range cs {
			var failures []string

			for _, constraint := range constraints {
				artifact := ws.Apply(name, constraint)

				if artifact == nil {
					failures = append(failures, constraint.String())
				} else {
					artifactsPicked = append(artifactsPicked, artifact)
				}
			}

			if len(failures) == 0 {
				continue
			}

			failure := fmt.Sprintf("unable to satisfy constraints for %s: %v", name, failures)
			allFailures = append(allFailures, failure)
		}

		cs = ConstraintSet{}
		for _, artifact := range artifactsPicked {
			for name, constraints := range artifact.dependsOn {
				for _, constraint := range constraints {
					cs.AddConstraintWithOrigin(name, constraint.RangeString, artifact)
				}
			}
		}
	}

	if len(allFailures) > 0 {
		return nil, errors.New(strings.Join(allFailures, "\n"))
	} else {
		return ws.Picks(), nil
	}
}

// Support types

type WorkingSet struct {
	source          ArtifactSource
	artifactsByName map[string][]Artifact
	constraints     ConstraintSet
}

func (ws *WorkingSet) EnsureCache(name string) []Artifact {
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

func (ws *WorkingSet) Apply(name string, constraint Constraint) *Artifact {
	compiledRange := constraint.Range
	for _, c := range ws.constraints[name] {
		compiledRange = compiledRange.AND(c.Range)
	}
	ws.constraints[name] = append(ws.constraints[name], constraint)

	artifacts := ws.EnsureCache(name)
	for i, artifact := range artifacts {
		if compiledRange(artifact.version) {
			ws.artifactsByName[name] = artifacts[i:]
			return &artifact
		}
	}

	return nil
}

func (ws *WorkingSet) Picks() []Artifact {
	var result []Artifact

	for _, artifacts := range ws.artifactsByName {
		result = append(result, artifacts[0])
	}

	return result
}
