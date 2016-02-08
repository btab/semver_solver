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

	// TODO: improve this interface to look more like the one above (adder)
	constraints := ss.ConstraintSet{}
	constraints["A"] = []ss.Constraint{
		{RangeString: ">1.1.0", Range: MustParseRange(">1.1.0")},
	}

	solver := ss.Solver{Source: source}

	solver.Solve(constraints)
}

func MustParseRange(s string) semver.Range {
	r, err := semver.ParseRange(s)

	if err != nil {
		panic(`semver: ParseRange(` + s + `): ` + err.Error())
	}

	return r
}
