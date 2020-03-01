//+build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const dsName string = "sheets-datasource"

func buildBackend(variant string, enableDebug bool, env map[string]string) error {
	varStr := ""
	if variant != "" {
		varStr = fmt.Sprintf("_%s", variant)
	}
	args := []string{
		"build", "-o", fmt.Sprintf("dist/%s%s", dsName, varStr), "-tags", "netgo",
	}
	if enableDebug {
		args = append(args, "-gcflags=all=-N -l")
	} else {
		args = append(args, []string{"-ldflags", "-w"}...)
	}
	args = append(args, "./pkg")
	// TODO: Change to sh.RunWithV once available.
	if err := sh.RunWith(env, "go", args...); err != nil {
		return err
	}

	return nil
}

type Build mg.Namespace

// Backend builds the back-end plugin.
func (Build) Backend() error {
	return buildBackend("", false, map[string]string{})
}

// BackendDebug builds the back-end plugin in debug mode.
func (Build) BackendDebug() error {
	return buildBackend("", true, map[string]string{})
}

// BackendLinux builds the back-end plugin for Linux.
func (Build) BackendLinux() error {
	env := map[string]string{
		"GOARCH": "amd64",
		"GOOS":   "linux",
	}
	return buildBackend("linux_amd64", false, env)
}

// BackendLinuxDebug builds the back-end plugin for Linux in debug mode.
func (Build) BackendLinuxDebug() error {
	env := map[string]string{
		"GOARCH": "amd64",
		"GOOS":   "linux",
	}
	return buildBackend("linux_amd64", true, env)
}

// Frontend builds the front-end for production.  Note that this build script will also
// clean the `dist` folder
func (Build) Frontend() error {
	mg.Deps(Deps)
	return sh.RunV("./node_modules/.bin/grafana-toolkit", "plugin:build")
}

// BuildAll builds both back-end and front-end components.
func BuildAll() {
	b := Build{}
	// Frontend goes first and cleans the 'dist' folder
	mg.Deps(b.Frontend, b.Backend, b.BackendLinux)
}

// Deps installs dependencies.
func Deps() error {
	return sh.RunV("yarn", "install")
}

// Test runs all tests.
func Test() error {
	mg.Deps(Deps)

	if err := sh.RunV("go", "test", "./pkg/..."); err != nil {
		return nil
	}
	return sh.RunV("yarn", "test")
}

// Lint lints the sources.
func Lint() error {
	return sh.RunV("golangci-lint", "run", "./...")
}

// Format formats the sources.
func Format() error {
	if err := sh.RunV("gofmt", "-w", "."); err != nil {
		return err
	}

	return nil
}

// Dev builds the plugin in dev mode
func Dev() error {
	b := Build{}

	// First build the frontend
	if err := sh.RunV("./node_modules/.bin/grafana-toolkit", "plugin:dev"); err != nil {
		return err
	}

	// Then a debug backend
	mg.Deps(b.BackendDebug)

	return nil
}

// Watch will build the plugin in dev mode and then update when the frontend files change
func Watch() error {
	b := Build{}
	mg.Deps(b.BackendDebug)
	return sh.RunV("./node_modules/.bin/grafana-toolkit", "plugin:dev", "--watch")
}

var Default = BuildAll
