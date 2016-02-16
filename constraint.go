package semver_solver

import "github.com/blang/semver"

type Constraint struct {
	ArtifactName string
	RangeString  string
	Range        semver.Range
}

func NewConstraint(artifactName, rangeString string) (*Constraint, error) {
	r, err := semver.ParseRange(rangeString)

	if err != nil {
		return nil, err
	}

	return &Constraint{
		ArtifactName: artifactName,
		RangeString:  rangeString,
		Range:        r,
	}, nil
}

func (c *Constraint) String() string {
	return c.ArtifactName + c.RangeString
}
