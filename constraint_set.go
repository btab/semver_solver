package semver_solver

import (
	"strings"

	"github.com/blang/semver"
)

type ConstraintSet map[string][]Constraint

func (cs ConstraintSet) AddConstraint(artifactName, rangeString string) error {
	return cs.AddConstraintWithOrigin(artifactName, rangeString, nil)
}

func (cs ConstraintSet) AddConstraintWithOrigin(artifactName, rangeString string, origin *Artifact) error {
	r, err := semver.ParseRange(rangeString)

	if err != nil {
		return err
	}

	c := Constraint{
		RangeString: rangeString,
		Range:       r,
		Origin:      origin,
	}

	cs[artifactName] = append(cs[artifactName], c)

	return nil
}

func (cs ConstraintSet) String() string {
	var parts []string

	for artifactName, constraints := range cs {
		var constraintStrings []string

		for _, constraint := range constraints {
			constraintStrings = append(constraintStrings, artifactName+constraint.String())
		}

		parts = append(parts, strings.Join(constraintStrings, ", "))
	}

	return strings.Join(parts, "; ")
}
