package semver_solver

import "github.com/blang/semver"

type Constraint struct {
	Range  semver.Range
	Origin *Artifact
}
