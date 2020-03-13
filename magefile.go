// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	err := Versionize()
	if err != nil {
		return err
	}
	cmd := exec.Command("go", "build", "./cmd/versionize")
	return cmd.Run()
}

// Create version info for build process
func Versionize() error {
	_, err := exec.LookPath("versionize")
	if err != nil {
		return err
	}
	cmd := exec.Command("versionize", "-go", "cmd/versionize/metadata.go")
	err = cmd.Run()
  return err
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	cmd := exec.Command("go", "install", "./cmd/versionize")
	return cmd.Run()
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing deps...")
	return nil
	//cmd := exec.Command("go", "get", "github.com/stretchr/piglatin")
	//return cmd.Run()
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("versionize")
}
