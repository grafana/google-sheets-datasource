//+build mage

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/google/gops/goprocess"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const dsName string = "sheets-datasource"

func getExecutableName(os string, arch string) string {
	exeName := fmt.Sprintf("%s_%s_%s", dsName, os, arch)
	if "windows" == os {
		exeName = fmt.Sprintf("%s.exe", exeName)
	}
	return exeName
}

func findRunningProcess(exe string) *goprocess.P {
	for _, process := range goprocess.FindAll() {
		if strings.HasSuffix(process.Path, exe) {
			return &process
		}
	}
	return nil
}

func killProcess(process *goprocess.P) error {
	log.Printf("Killing: %s (%d)", process.Path, process.PID)
	return syscall.Kill(process.PID, 9)
}

func buildBackend(os string, arch string, enableDebug bool) error {
	exeName := getExecutableName(os, arch)

	args := []string{
		"build", "-o", fmt.Sprintf("dist/%s", exeName), "-tags", "netgo",
	}
	if enableDebug {
		args = append(args, "-gcflags=all=-N -l")
	} else {
		args = append(args, []string{"-ldflags", "-w"}...)
	}
	args = append(args, "./pkg")

	env := map[string]string{
		"GOARCH": arch,
		"GOOS":   os,
	}

	// TODO: Change to sh.RunWithV once available.
	return sh.RunWith(env, "go", args...)
}

// Build is a namespace.
type Build mg.Namespace

// Linux builds the back-end plugin for Linux.
func (Build) Linux() error {
	return buildBackend("linux", "amd64", false)
}

// Windows builds the back-end plugin for Windows.
func (Build) Windows() error {
	return buildBackend("windows", "amd64", false)
}

// Darwin builds the back-end plugin for OSX.
func (Build) Darwin() error {
	return buildBackend("darwin", "amd64", false)
}

// Debug builds the debug version for the current platform
func (Build) Debug() error {
	return buildBackend(runtime.GOOS, runtime.GOARCH, true)
}

// Backend build a production build for all platforms
func (Build) Backend() {
	b := Build{}
	mg.Deps(b.Linux, b.Windows, b.Darwin)
}

// // Frontend builds the front-end for production.
// func (Build) Frontend() error {
// 	mg.Deps(Deps)
// 	return sh.RunV("./node_modules/.bin/grafana-toolkit", "plugin:build")
// }

// // FrontendDev builds the front-end for development.
// func (Build) FrontendDev() error {
// 	mg.Deps(Deps)
// 	return sh.RunV("./node_modules/.bin/grafana-toolkit", "plugin:dev")
// }

// BuildAll builds both back-end and front-end components.
func BuildAll() {
	b := Build{}
	mg.Deps(b.Backend)
}

// // Deps installs dependencies.
// func Deps() error {
// 	return nil //sh.RunV("yarn", "install")
// }

// Test run backend tests
func Test() error {
	//mg.Deps(Deps)

	if err := sh.RunV("go", "test", "./pkg/..."); err != nil {
		return nil
	}
	return nil // sh.RunV("yarn", "test")
}

// Coverage runs backend tests and make a coverage report
func Coverage() error {
	os.MkdirAll(filepath.Join(".", "coverage"), os.ModePerm)

	if err := sh.RunV("go", "test", "./pkg/...", "-v", "-cover", "-coverprofile=coverage/backend.out"); err != nil {
		return nil
	}

	if err := sh.RunV("go", "tool", "cover", "-html=coverage/backend.out", "-o", "coverage/backend.html"); err != nil {
		return nil
	}

	return nil
}

// Lint audits the source style
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

// // Dev builds the plugin in dev mode.
// func Dev() error {
// 	b := Build{}
// 	mg.Deps(b.BackendLinuxDebug, b.FrontendDev) // TODO: only the current architecture
// 	return nil
// }

// // Watch will build the plugin in dev mode and then update when the frontend files change.
// func Watch() error {
// 	b := Build{}
// 	mg.Deps(b.BackendLinuxDebug)

// 	// The --watch will never return
// 	return sh.RunV("./node_modules/.bin/grafana-toolkit", "plugin:dev", "--watch")
// }

// Clean cleans build artifacts, by deleting the dist directory.
func Clean() error {
	return os.RemoveAll("dist")
}

// Debugger makes a new debug build and attaches dvl
func Debugger() error {
	// 1. kill any running instance
	exeName := getExecutableName(runtime.GOOS, runtime.GOARCH)

	// Kill any running processs
	process := findRunningProcess(exeName)
	if process != nil {
		err := killProcess(process)
		if err != nil {
			return err
		}
	}

	// Debug build
	b := Build{}
	mg.Deps(b.Debug)

	if runtime.GOOS == "linux" {
		// 	ptrace_scope=`cat /proc/sys/kernel/yama/ptrace_scope`
		// 	if [ "$ptrace_scope" != 0 ]; then
		// 	  echo "WARNING: ptrace_scope set to value other than 0, this might prevent debugger from connecting, try writing \"0\" to /proc/sys/kernel/yama/ptrace_scope.
		//   Read more at https://www.kernel.org/doc/Documentation/security/Yama.txt"
		// 	  read -p "Set ptrace_scope to 0? y/N (default N)" set_ptrace_input
		// 	  if [ "$set_ptrace_input" == "y" ] || [ "$set_ptrace_input" == "Y" ]; then
		// 		echo 0 | sudo tee /proc/sys/kernel/yama/ptrace_scope
		// 	  fi
		// 	fi
	}

	// Wait for grafana to start plugin
	for i := 0; i < 20; i++ {
		process := findRunningProcess(exeName)
		if process != nil {
			log.Printf("Running PID: %d", process.PID)

			// dlv attach ${PLUGIN_PID} --headless --listen=:${PORT} --api-version 2 --log
			if err := sh.RunV("dvl",
				"attach",
				strconv.Itoa(process.PID),
				"--headless",
				"--listen=:3222",
				"--api-version", "2",
				"--log"); err != nil {
				return err
			}
			// And then kill dvl
			return sh.RunV("pkill", "dlv")
		}

		log.Printf("waiting for grafana to start: %s...", exeName)
		time.Sleep(250 * time.Millisecond)
	}
	return fmt.Errorf("could not find process: %s, perhaps grafana is not running?", exeName)
}

// Default configures the default target.
var Default = BuildAll
