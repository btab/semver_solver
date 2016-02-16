// ls *.go tester/main.go tester/scenarios/* | entr -c go run tester/main.go

package main

import (
	"bufio"
	"fmt"
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
		println()
	}
}

func AssertNoError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Scenario struct {
	name         string
	source       ss.MockArtifactSource
	constraints  []*ss.Constraint
	expectations map[string]semver.Version
}

var (
	availableMatcher   = regexp.MustCompile("(\\S+?)@(.+?)( -> .*?)?$")
	constraintsMatcher = regexp.MustCompile("(\\S+?)([=<>!]{1,2}.+?)$")
	expectMatcher      = regexp.MustCompile("(\\S+?)@(.+?)$")
)

func ParseScenario(path string) *Scenario {
	scenario := &Scenario{
		name:         filepath.Base(path),
		source:       ss.MockArtifactSource{},
		constraints:  nil,
		expectations: map[string]semver.Version{},
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

		if strings.HasPrefix(line, "#") {
			continue
		}

		if !strings.HasPrefix(line, "\t") {
			section = strings.TrimSpace(line)
			continue
		}

		line = strings.TrimSpace(line)

		switch section {

		case "Available":
			parts := availableMatcher.FindStringSubmatch(line)

			if len(parts) != 4 {
				die(scenario, line, "wrong number of sub-parts")
			}

			var deps []*ss.Constraint

			if len(parts[3]) > 0 {
				depsString := strings.TrimPrefix(parts[3], " -> ")

				for _, depString := range strings.Split(depsString, ",") {
					depParts := constraintsMatcher.FindStringSubmatch(depString)

					if len(depParts) != 3 {
						msg := fmt.Sprintf("wrong number of sub-parts in '%s'", depString)
						die(scenario, line, msg)
					}

					dep, err := ss.NewConstraint(depParts[1], depParts[2])
					if err != nil {
						die(scenario, line, err.Error())
					}
					deps = append(deps, dep)
				}
			}

			if err := scenario.source.AddArtifactWithDeps(parts[1], parts[2], deps); err != nil {
				die(scenario, line, err.Error())
			}

		case "Constraints":
			parts := constraintsMatcher.FindStringSubmatch(line)

			if len(parts) != 3 {
				die(scenario, line, "wrong number of sub-parts")
			}

			c, err := ss.NewConstraint(parts[1], parts[2])
			if err != nil {
				die(scenario, line, err.Error())
			}
			scenario.constraints = append(scenario.constraints, c)

		case "Expect":
			parts := expectMatcher.FindStringSubmatch(line)

			if len(parts) != 3 {
				die(scenario, line, "wrong number of sub-parts")
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
	log.Printf("scenario %s running...\n", s.name)

	solver := ss.Solver{Source: s.source}

	var cs []*ss.Constraint
	for _, c := range s.constraints {
		cs = append(cs, c)
	}

	artifacts, err := solver.Solve(cs)

	log.Printf("scenario %s picks: %v\n", s.name, artifacts)
	log.Printf("scenario %s error: %v\n", s.name, err)

	if err != nil && len(s.expectations) > 0 {
		log.Printf("scenario %s unexpected error: %s",
			s.name, err.Error())
		return
	}

	for _, artifact := range artifacts {
		expectedVersion, ok := s.expectations[artifact.Name]

		if !ok || expectedVersion.NE(artifact.Version) {
			log.Printf("scenario %s generated unexpected artifact: %s",
				s.name, artifact.String())
			return
		}

		delete(s.expectations, artifact.Name)
	}

	for name, version := range s.expectations {
		log.Printf("scenario %s failed to generate expected artifact: %s@%s",
			s.name, name, version.String())
		return
	}
}
