package semver_solver

import (
	"fmt"
	"log"
	"sort"

	"github.com/blang/semver"
)

// This would become the 'simple solver' if we ever added other solvers (SAT etc)

type Solver struct {
	Source ArtifactSource
}

func (s *Solver) Solve(initCS ConstraintSet) {
	ws := WorkingSet{
		source:          s.Source,
		artifactsByName: make(map[string][]Artifact),
	}

	var allErrors []error

	for name, constraints := range initCS {
		var failures []string

		for _, constraint := range constraints {
			if ws.ConsumeUntil(name, constraint.Range) == false {
				failures = append(failures, constraint.String())
			}
		}

		if len(failures) == 0 {
			continue
		}

		err := fmt.Errorf("unable to satisfy constraints for %s: %v", name, failures)
		allErrors = append(allErrors, err)
	}

	if len(allErrors) > 0 {
		log.Println(allErrors)
		return
	}

	cs := ConstraintSet{}
	for name := range initCS {
		cs[name] = initCS[name]
	}

	// picks := map[string]Artifact{}
	// picks are by definition the heads of the working set

	log.Println(_solve(ws, cs))
	log.Println(ws)
}

func _solve(ws WorkingSet, cs ConstraintSet) error {
	for name, constraints := range cs {
		for _, constraint := range constraints {
			ws.ConsumeUntil(name, constraint.Range)
		}
	}

	// for name, svRange := range workingCS {
	// chomp through non-matching items in working set if at root
	// otherwise return error if non-root and head does not match constraint

	// filtered := Filter(workingAS[name], svRange)
	// what if this further filters something we've already filtered?
	//   -- simple case is that it filters out an existing pick (backtrack)
	//   -- but it could filter out other things
	//   -- what if, instead of a complete filter, the working set for each artifact
	//   ---- existed to eliminate its head as a candidate?
	//   -- so as soon as we pick something we pop it off the head?

	// if len(filtered) == 0 {
	// 	// TODO: instrument for friendly debug (we can scan the picks for a list of active constraints)
	// 	return errors.New("constraints filtered out all possible picks for " + name)
	// }

	// if picked, ok := picks[name]; ok {
	// 	if filtered[0] != picked {
	// 		return errors.New("constraints filtered out existing pick for " + name)
	// 	}
	// 	return nil
	// }

	// picks[name] = filtered[0]
	// new constaints

	// }

	// any new constraints? do another getAllNewArtifactInfo(sourceAS, workingAS, names)

	return nil
}

// Support types

type WorkingSet struct {
	source          ArtifactSource
	artifactsByName map[string][]Artifact
}

func (ws *WorkingSet) EnsureCache(name string) ([]Artifact, bool) {
	artifacts, ok := ws.artifactsByName[name]

	if ok {
		return artifacts, true
	}

	artifacts = ws.source.AllVersionsOf(name)

	// copy for isolation
	localCopy := make([]Artifact, len(artifacts))
	copy(localCopy, artifacts)

	// reverse sort to make all the 'default to latest' optimizations work
	sort.Sort(sort.Reverse(SortableArtifacts(localCopy)))

	ws.artifactsByName[name] = localCopy
	return localCopy, false
}

// TODO: rename to Enforce or something like that
func (ws *WorkingSet) ConsumeUntil(name string, svRange semver.Range) (ok bool) {
	artifacts, wasInCache := ws.EnsureCache(name)

	// TODO: need to consume only if the cache is freshly populated
	//       otherwise we have to error on non-head match, chomp the head and start again

	if wasInCache {
		return svRange(artifacts[0].version)
	}

	for i, artifact := range artifacts {
		if svRange(artifact.version) {
			ws.artifactsByName[name] = artifacts[i:]
			return true
		}
	}

	return false
}

func (ws *WorkingSet) Get(name string) []Artifact {
	artifacts, _ := ws.EnsureCache(name)
	return artifacts
}
