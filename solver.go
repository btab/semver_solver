package semver_solver

import (
	"log"
	"sort"
)

// This would become the 'simple solver' if we ever added other solvers (SAT etc)

type Solver struct {
	Source ArtifactSource
}

func (s Solver) Solve(initCS ConstraintSet) {
	ws := WorkingSet{
		source:          s.Source,
		artifactsByName: make(map[string][]Artifact),
	}

	cs := ConstraintSet{}
	for name := range initCS {
		cs[name] = initCS[name]
	}

	// picks := map[string]Artifact{}
	// picks are by definition the heads of the working set

	log.Println(s._solve(ws, cs))
	log.Println(ws)
}

func (s Solver) _solve(ws WorkingSet, workingCS ConstraintSet) error {
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

func (ws WorkingSet) EnsureCache(name string) []Artifact {
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

// this is the simple version which may be replaced by Head, Tail and Chomp operations
func (ws WorkingSet) Get(name string) []Artifact {
	artifacts := ws.EnsureCache(name)
	return artifacts
}
