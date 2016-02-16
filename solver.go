package semver_solver

import (
	"fmt"
	"sort"
	"strings"
)

type Solver struct {
	Source ArtifactSource
}

type cell struct {
	constraint *Constraint
	parent     *cell
	picks      []Artifact
	children   []*cell
	garbage    bool
}

type conflictSnapshot struct {
	cellsByArtifactName map[string]*cell
	conflictingCell     *cell
}

func (cs *conflictSnapshot) String() string {
	var activated []string

	for _, cell := range cs.cellsByArtifactName {
		activated = append(activated, cell.picks[0].String())
	}

	cell := cs.conflictingCell

	return "constraint " + cell.constraint.String() +
		" from " + cell.parent.picks[0].String() +
		" conflicted with picked artifacts [" + strings.Join(activated, ", ") + "]"
}

func (s *Solver) Solve(constraints []*Constraint) ([]*Artifact, error) {
	var pendingCells []*cell
	for _, c := range constraints {
		cell := &cell{
			constraint: c,
		}
		pendingCells = append(pendingCells, cell)
	}

	cellsByArtifactName := map[string]*cell{}
	var firstConflict *conflictSnapshot

	for len(pendingCells) > 0 {
		var newPendingCells []*cell

		for _, pCell := range pendingCells {
			if pCell.garbage {
				continue
			}

			constraint := pCell.constraint
			name := constraint.ArtifactName

			existingCell := cellsByArtifactName[name]

			// Artifact never seen before, activate its cell
			//  - record its activation with the global list
			//  - add its dependencies to the tail of the pending cells

			if existingCell == nil {
				cellsByArtifactName[name] = pCell

				pCell.picks = retrieveAllVersions(s.Source, name)
				matchIndex := indexOfFirstMatch(pCell.picks, constraint)
				if matchIndex == -1 {
					return nil, fmt.Errorf("no artifacts match %s", constraint.String())
				}

				pick(pCell, matchIndex, &newPendingCells)

				continue
			}

			// New constraint is compatible with existing pick

			if constraint.Range(existingCell.picks[0].Version) {
				continue
			}

			// New constraint is incompatible with existing pick
			//  - log if this is the first such conflict (in case we can't find a solution)
			//  - backtrack up the tree, until an alternative path is found

			if firstConflict == nil {
				firstConflict = &conflictSnapshot{
					cellsByArtifactName: map[string]*cell{},
					conflictingCell:     pCell,
				}

				for name, cell := range cellsByArtifactName {
					firstConflict.cellsByArtifactName[name] = cell
				}
			}

			matchIndex := indexOfFirstMatch(existingCell.picks[1:], existingCell.constraint)
			if matchIndex == -1 {
				return nil, fmt.Errorf("no solutions found: %v", firstConflict)
			}

			pruneChildren(existingCell)

			pick(existingCell, matchIndex+1, &newPendingCells)
		}

		pendingCells = newPendingCells
	}

	var artifacts []*Artifact
	for _, cell := range cellsByArtifactName {
		artifacts = append(artifacts, &cell.picks[0])
	}

	return artifacts, nil
}

func indexOfFirstMatch(artifacts []Artifact, constraint *Constraint) int {
	for i, a := range artifacts {
		if constraint.Range(a.Version) {
			return i
		}
	}

	return -1
}

func retrieveAllVersions(s ArtifactSource, name string) []Artifact {
	allArtifacts := s.AllVersionsOf(name)

	// copy for isolation
	localCopy := make([]Artifact, len(allArtifacts))
	copy(localCopy, allArtifacts)

	// reverse sort to make all the 'default to latest' optimizations work
	sort.Sort(sort.Reverse(SortableArtifacts(localCopy)))

	return localCopy
}

func pruneChildren(fromCell *cell) {
	for _, cell := range fromCell.children {
		cell.garbage = true
		pruneChildren(cell)
	}

	fromCell.children = nil
}

func pick(c *cell, index int, newCells *[]*cell) {
	c.picks = c.picks[index:]

	cells := cellsFromDeps(c)
	c.children = append(c.children, cells...)
	*newCells = append(*newCells, cells...)
}

func cellsFromDeps(c *cell) []*cell {
	var cells []*cell

	for _, dep := range c.picks[0].DependsOn {
		newCell := &cell{
			constraint: dep,
			parent:     c,
		}
		cells = append(cells, newCell)
	}

	return cells
}
