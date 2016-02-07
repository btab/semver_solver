package main

import (
	"errors"
	"log"
	"sort"

	"github.com/blang/semver"
)

var (
	NoConstraints = ConstraintSet{}
)

func main() {
	// simulate repo
	as := ArtifactSet{
		"A": {
			{"A", semver.MustParse("1.0.0"), NoConstraints},
			{"A", semver.MustParse("1.1.0"), NoConstraints},
			{"A", semver.MustParse("1.2.0"), NoConstraints},
			{"A", semver.MustParse("1.4.0"), NoConstraints},
			{"A", semver.MustParse("2.1.0"), NoConstraints},
			{"A", semver.MustParse("2.1.2"), NoConstraints},
		},
		"B": {
			{"B", semver.MustParse("1.1.0"), NoConstraints},
			{"B", semver.MustParse("1.2.0"), NoConstraints},
		},
	}

	constraints := ConstraintSet{
		"A": MustParseRange(">1.1.0"),
	}

	// TODO: need a function that can look at a set of constraints and establish that they are disjoint
	//         i.e. A <1 and A >2 could never be solved in principle
	//         this will be handy as every transitive pick can introduce new contraints and we should
	//           reject picks whose added dependecies would be disjoint with existing constraints
	//         likely a pairwise check - note can't get inside existing constraint as it's built of functions

	log.Println(as)
	log.Println(constraints)

	solution, err := Solve(as, constraints)

	log.Println(solution)
	log.Println(err)
}

type Artifact struct {
	name      string
	version   semver.Version
	dependsOn ConstraintSet
}

// Support sorting
type SortableArtifacts []Artifact

func (a SortableArtifacts) Len() int {
	return len(a)
}

// Swap swaps two versions inside the collection by its indices
func (a SortableArtifacts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less checks if version at index i is less than version at index j
func (a SortableArtifacts) Less(i, j int) bool {
	return a[i].version.LT(a[j].version)
}

// --sorting

// ArtifactSet (dummy implementation of a repo) - this will actually become an interface
type ArtifactSet map[string][]Artifact

func (as ArtifactSet) allVersionsOf(artifactNames []string) ArtifactSet {
	result := ArtifactSet{}

	for _, name := range artifactNames {
		result[name] = as[name] // TODO: will need to reverse sort a copy of the slices after we get them
	}

	return result
}

// -- ArtifactSet

type ConstraintSet map[string]semver.Range

func Filter(artifacts []Artifact, svRange semver.Range) []Artifact {
	var filtered []Artifact

	for _, a := range artifacts {
		if svRange(a.version) {
			filtered = append(filtered, a)
		}
	}

	return filtered
}

func Solve(sourceAS ArtifactSet, sourceCS ConstraintSet) (map[string]Artifact, error) {
	workingAS := ArtifactSet{}

	workingCS := ConstraintSet{}
	var names []string
	for name := range sourceCS {
		workingCS[name] = sourceCS[name]
		names = append(names, name)
	}

	getAllNewArtifactInfo(sourceAS, workingAS, names)

	picks := map[string]Artifact{}

	err := _solve(sourceAS, workingAS, workingCS, picks)

	return picks, err
}

func _solve(sourceAS, workingAS ArtifactSet, workingCS ConstraintSet, picks map[string]Artifact) error {
	for name, svRange := range workingCS {
		filtered := Filter(workingAS[name], svRange)

		if len(filtered) == 0 {
			// TODO: instrument for friendly debug (we can scan the picks for a list of active constraints)
			return errors.New("no solution")
		}

	}

	// any new constraints? do another getAllNewArtifactInfo(sourceAS, workingAS, names)

	return nil
}

func getAllNewArtifactInfo(from, to ArtifactSet, names []string) {
	var namesToGet []string

	for _, name := range names {
		if _, ok := to[name]; !ok {
			namesToGet = append(namesToGet, name)
		}
	}

	if len(namesToGet) == 0 {
		return
	}

	log.Printf("cache miss for artifacts: %v, (requested: %v)\n", namesToGet, names)

	for name, artifacts := range from.allVersionsOf(namesToGet) {
		to[name] = append(to[name], artifacts...)            // copy for isolation
		sort.Sort(sort.Reverse(SortableArtifacts(to[name]))) // and reverse sort to make all the 'default to latest' optimizations work
	}
}

func MustParseRange(s string) semver.Range {
	r, err := semver.ParseRange(s)

	if err != nil {
		panic(`semver: ParseRange(` + s + `): ` + err.Error())
	}

	return r
}
