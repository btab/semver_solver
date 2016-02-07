package semver_solver

type SortableArtifacts []Artifact

func (a SortableArtifacts) Len() int {
	return len(a)
}

func (a SortableArtifacts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a SortableArtifacts) Less(i, j int) bool {
	return a[i].version.LT(a[j].version)
}
