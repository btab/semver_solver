package semver_solver

import (
	"github.com/blang/semver"
)

type ConstraintSet map[string][]Constraint

func (cs ConstraintSet) AddConstraint(artifactName, cText string) error {
	return cs.AddConstraintWithOrigin(artifactName, cText, nil)
}

func (cs ConstraintSet) AddConstraintWithOrigin(artifactName, cText string, origin *Artifact) error {
	r, err := semver.ParseRange(cText)

	if err != nil {
		return err
	}

	c := Constraint{
		RangeString: cText,
		Range:       r,
		Origin:      origin,
	}

	cs[artifactName] = append(cs[artifactName], c)

	return nil
}
