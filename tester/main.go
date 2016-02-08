package main

import (
	"log"
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

	constraints := ss.ConstraintSet{}
	AssertNoError(constraints.AddConstraint("A", "<2.0.0"))
	source.AddArtifactWithDeps("B", semver.MustParse("1.2.0"), constraints)

	constraints = ss.ConstraintSet{}
	AssertNoError(constraints.AddConstraint("A", ">1.1.0"))
	AssertNoError(constraints.AddConstraint("B", "=1.2.0"))

	solver := ss.Solver{Source: source}
	artifacts, err := solver.Solve(constraints)

	log.Println(artifacts)
	log.Println(err)
}

func AssertNoError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
