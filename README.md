SEMVER SOLVER
-------------
A golang native dependency resolver with full semver 2 support. Depends on https://github.com/blang/semver.

USAGE
-----

Install...
```
go get github.com/btab/semver_solver
```

Implement an ```ArtifactSource``` (or use the ```MockArtifactSource```), create a solver and solve for a set of constraints...
```go
package main

import (
	ss "github.com/btab/semver_solver"
)

func main() {
	source := ss.MockArtifactSource{}
	source.AddArtifact("foo", "1.0.0")

	solver := ss.Solver{source}

	constraint, _ := ss.NewConstraint("foo", "<2.0.0")
	constraints := []*ss.Constraint{constraint}

	artifacts, err := solver.Solve(constraints)
}
```

SCENARIO TESTING
----------------

From the root of the project...
```
go run tester/main.go
```

TODO
----

1. find some really nasty-complex real-world examples to smoke test this thing
2. look at apt, yum, bundler (molinillo), npm, cargo, maven, brew
3. generalize to a SAT(-3) solver (minisat / other DLL solver)?
