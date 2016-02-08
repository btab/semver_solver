package semver_solver

import (
	"errors"
	"fmt"
	"strings"
)

type Solver struct {
	Source ArtifactSource
}

func (s *Solver) Solve(cs ConstraintSet) ([]Artifact, error) {
	ws := newWorkingSet(s.Source)

	var allFailures []string

	for len(cs) > 0 && len(allFailures) == 0 {
		var artifactsPicked []*Artifact

		for name, constraints := range cs {
			var failures []string

			for _, constraint := range constraints {
				artifact := ws.apply(name, constraint)

				if artifact == nil {
					failures = append(failures, constraint.String())
				} else {
					artifactsPicked = append(artifactsPicked, artifact)
				}
			}

			if len(failures) == 0 {
				continue
			}

			failure := fmt.Sprintf("unable to satisfy constraints for %s: %v", name, failures)
			allFailures = append(allFailures, failure)
		}

		cs = ConstraintSet{}
		for _, artifact := range artifactsPicked {
			for name, constraints := range artifact.dependsOn {
				for _, constraint := range constraints {
					cs.AddConstraintWithOrigin(name, constraint.RangeString, artifact)
				}
			}
		}
	}

	if len(allFailures) > 0 {
		return nil, errors.New(strings.Join(allFailures, "\n"))
	} else {
		return ws.picks(), nil
	}
}
