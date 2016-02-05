package main

import (
	"fmt"

	"github.com/blang/semver"
)

func main() {
	v, _ := semver.Make("0.0.1-alpha.preview+123.github")
	fmt.Printf("Major: %d\n", v.Major)
	fmt.Printf("Minor: %d\n", v.Minor)
	fmt.Printf("Patch: %d\n", v.Patch)
	fmt.Printf("Pre: %s\n", v.Pre)
	fmt.Printf("Build: %s\n", v.Build)
}
