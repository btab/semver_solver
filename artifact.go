package semver_solver

import "github.com/blang/semver"

type Artifact struct {
	name      string
	version   semver.Version
	dependsOn ConstraintSet
}

func (a *Artifact) String() string {
	return a.name + "@" + a.version.String()
}
