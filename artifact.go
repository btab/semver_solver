package semver_solver

import "github.com/blang/semver"

type Artifact struct {
	name      string
	version   semver.Version
	dependsOn ConstraintSet
}
