package main

import (
	ss "semver_solver"

	"github.com/blang/semver"
)

func main() {
	source := ss.MockArtifactSource{}

	source.AddArtifact("A", semver.MustParse("1.0.0"))
	source.AddArtifact("A", semver.MustParse("1.2.0"))
	source.AddArtifact("A", semver.MustParse("1.3.0"))
	source.AddArtifact("A", semver.MustParse("2.1.0"))
	source.AddArtifact("A", semver.MustParse("2.1.2"))
	source.AddArtifact("B", semver.MustParse("1.1.0"))
	source.AddArtifact("B", semver.MustParse("1.2.0"))

	constraints := ss.ConstraintSet{}
	constraints["A"] = MustParseRange(">1.1.0")

	// TODO: need a function that can look at a set of constraints and establish that they are disjoint
	//         i.e. A <1 and A >2 could never be solved in principle
	//         this will be handy as every transitive pick can introduce new contraints and we should
	//           reject picks whose added dependecies would be disjoint with existing constraints
	//         likely a pairwise check - note can't get inside existing constraint as it's built of functions

	// log.Println(as)
	// log.Println(constraints)

	solver := ss.Solver{Source: source}

	solver.Solve(constraints)
}

// func Filter(artifacts []Artifact, svRange semver.Range) []Artifact {
// 	var filtered []Artifact

// 	for _, a := range artifacts {
// 		if svRange(a.version) {
// 			filtered = append(filtered, a)
// 		}
// 	}

// 	return filtered
// }

func MustParseRange(s string) semver.Range {
	r, err := semver.ParseRange(s)

	if err != nil {
		panic(`semver: ParseRange(` + s + `): ` + err.Error())
	}

	return r
}
