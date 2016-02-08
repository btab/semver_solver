package semver_solver

import "github.com/blang/semver"

type Constraint struct {
	RangeString string
	Range       semver.Range
	Origin      *Artifact
}

func (c *Constraint) String() string {
	originString := ""
	if c.Origin != nil {
		originString = " from " + c.Origin.String()
	}

	return c.RangeString + originString
}
