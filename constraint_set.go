package semver_solver

import (
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
