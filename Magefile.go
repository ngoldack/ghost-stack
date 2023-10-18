//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	BUILD_DIR       = "build"
	CLI_BINARY_NAME = "ghost-cli"
	CLI_BINARY      = BUILD_DIR + "/" + CLI_BINARY_NAME
)

// A build target for compiling your project.
func Build() error {
	fmt.Println("Building the project...")
	cmd := exec.Command("go", "build", "-v", "-o", CLI_BINARY, "./cmd/cli")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// A build target for running tests.
func Test() error {
	fmt.Println("Running tests...")
	cmd := exec.Command("go", "test", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// A build target for cleaning up build artifacts.
func Clean() error {
	fmt.Println("Cleaning up...")
	return os.RemoveAll("output_binary")
}
