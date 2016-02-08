package main

import (
	"log"
	ss "semver_solver"
)

func main() {
	source := ss.MockArtifactSource{}

	AssertNoError(source.AddArtifact("A", "1.0.0"))
	AssertNoError(source.AddArtifact("A", "1.2.0"))
	AssertNoError(source.AddArtifact("A", "1.3.0"))
	AssertNoError(source.AddArtifact("A", "2.1.0"))
	AssertNoError(source.AddArtifact("A", "2.1.2"))

	AssertNoError(source.AddArtifact("B", "1.1.0"))

	constraints := ss.ConstraintSet{}
	AssertNoError(constraints.AddConstraint("A", "<2.0.0"))
	AssertNoError(source.AddArtifactWithDeps("B", "1.2.0", constraints))

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
