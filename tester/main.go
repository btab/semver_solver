// ls *.go tester/main.go tester/scenarios/* | entr -c go run tester/main.go

package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	ss "semver_solver"
	"strings"

	"github.com/blang/semver"
)

func main() {
	scenarioPaths, err := filepath.Glob("tester/scenarios/*")
	AssertNoError(err)

	for _, path := range scenarioPaths {
		ParseScenario(path).Run()
	}

	source := ss.MockArtifactSource{}

	AssertNoError(source.AddArtifact("A", "1.0.0"))
	AssertNoError(source.AddArtifact("A", "1.2.0"))
	AssertNoError(source.AddArtifact("A", "1.3.0"))
	AssertNoError(source.AddArtifact("A", "2.1.0"))
	AssertNoError(source.AddArtifact("A", "2.1.2"))

	AssertNoError(source.AddArtifact("B", "1.1.0"))

	constraints := ss.ConstraintSet{}
	AssertNoError(constraints.AddConstraint("A", "<2.0.0"))
	AssertNoError(source.AddArtifactWithDeps("B", "1.2.0", constraints))

	constraints = ss.ConstraintSet{}
	AssertNoError(constraints.AddConstraint("A", ">1.1.0"))
	AssertNoError(constraints.AddConstraint("B", "=1.2.0"))

	solver := ss.Solver{Source: source}
	artifacts, err := solver.Solve(constraints)

	log.Println(artifacts)
	log.Println(err)
}

func AssertNoError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Scenario struct {
	name         string
	source       ss.MockArtifactSource
	constraints  ss.ConstraintSet
	expectations map[string]semver.Version
}

var (
	availableMatcher   = regexp.MustCompile("(\\S+?)@(.+?)$")
	constraintsMatcher = regexp.MustCompile("(\\S+?)[=>!]{2}(.+?)$")
	expectMatcher      = regexp.MustCompile("(\\S+?)@(.+?)$")
)

func ParseScenario(path string) *Scenario {
	scenario := &Scenario{
		name:         filepath.Base(path),
		source:       ss.MockArtifactSource{},
		constraints:  ss.ConstraintSet{},
		expectations: map[string]semver.Version{},
		// TODO: add support for expected errors
	}

	file, err := os.Open(path)
	AssertNoError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	section := ""

	die := func(scenario *Scenario, line, msg string) {
		log.Fatalf("unable to parse line in scenario '%s': %s (%s)",
			scenario.name, line, msg)
	}

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "\t") {
			section = strings.TrimSpace(line)
			continue
		}

		line = strings.TrimSpace(line)

		switch section {

		case "Available": // TODO: add DEP support
			parts := availableMatcher.FindStringSubmatch(line)

			if len(parts) != 3 {
				die(scenario, line, "expected 2 sub-parts")
			}

			if err := scenario.source.AddArtifact(parts[1], parts[2]); err != nil {
				die(scenario, line, err.Error())
			}

		case "Constraints":
			parts := constraintsMatcher.FindStringSubmatch(line)

			if len(parts) != 3 {
				die(scenario, line, "expected 2 sub-parts")
			}

			if err := scenario.constraints.AddConstraint(parts[1], parts[2]); err != nil {
				die(scenario, line, err.Error())
			}

		case "Expect":
			parts := expectMatcher.FindStringSubmatch(line)

			if len(parts) != 3 {
				die(scenario, line, "expected 2 sub-parts")
			}

			version, err := semver.Parse(parts[2])
			if err != nil {
				die(scenario, line, err.Error())
			}

			scenario.expectations[parts[1]] = version

		default:
			die(scenario, line, "in unknown section: "+section)
		}
	}

	AssertNoError(scanner.Err())

	return scenario
}

func (s *Scenario) Run() {
	solver := ss.Solver{s.source}

	artifacts, err := solver.Solve(s.constraints)

	if err != nil && len(s.expectations) > 0 {
		log.Printf("scenario %s unexpected error: %s\n", err.Error())
		return
	}

	for _, artifact := range artifacts {
		expectedVersion, ok := s.expectations[artifact.Name]

		if !ok || expectedVersion.NE(artifact.Version) {
			log.Printf("scenario %s generated unexpected artifact: %s\n",
				s.name, artifact.String())
			return
		}

		delete(s.expectations, artifact.Name)
	}

	for name, version := range s.expectations {
		log.Printf("scenario %s failed to generate expected artifact: %s@%s\n",
			s.name, name, version.String())
		return
	}
}
