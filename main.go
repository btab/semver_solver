package main

import (
	"errors"
	"log"

	"github.com/blang/semver"
)

var (
	NoConstraints = []Constraint{}
)

func main() {
	availableArtifacts := []Artifact{
		{"A", semver.MustParse("1.0.0"), NoConstraints},
		{"A", semver.MustParse("1.1.0"), NoConstraints},
		{"A", semver.MustParse("1.2.0"), NoConstraints},
		{"A", semver.MustParse("1.4.0"), NoConstraints},
		{"A", semver.MustParse("2.1.0"), NoConstraints},
		{"A", semver.MustParse("2.1.2"), NoConstraints},
		{"B", semver.MustParse("1.1.0"), NoConstraints},
		{"B", semver.MustParse("1.2.0"), NoConstraints},
	}

	constraints := []Constraint{
		{"A", MustParseRange(">1.1.0")},
	}

	log.Println(availableArtifacts)
	log.Println(constraints)

	solution, err := Solve(constraints)

	log.Println(solution)
	log.Println(err)
}

type Artifact struct {
	name      string
	version   semver.Version
	dependsOn []Constraint
}

type Constraint struct {
	name string
	rng  semver.Range
}

func Solve(constraints []Constraint) ([]Artifact, error) {
	return nil, errors.New("no solution")
}

func MustParseRange(s string) semver.Range {
	r, err := semver.ParseRange(s)

	if err != nil {
		panic(`semver: ParseRange(` + s + `): ` + err.Error())
	}

	return r
}
