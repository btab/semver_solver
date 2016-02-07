package semver_solver

type ArtifactSource interface {
	AllVersionsOf(name string) []Artifact
}
