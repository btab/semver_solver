package semver_solver

import "github.com/blang/semver"

type Artifact struct {
	Name      string
	Version   semver.Version
	DependsOn []*Constraint
}

func (a *Artifact) String() string {
	return a.Name + "@" + a.Version.String()
}
