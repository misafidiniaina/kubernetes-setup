package dependencies

import (
	"os/exec"
	"strings"
)

type CommandRunner interface {
	LookPath(name string) (string, error)
	Version(name string) (string, error)
}

type Dependency struct {
	Name      string
	Path      string
	Version   string
	Required  bool
	Installed bool
}

type CheckResult struct {
	Dependencies []Dependency
	Ready        bool
}

func (result CheckResult) Dependency(name string) *Dependency {
	for index := range result.Dependencies {
		if result.Dependencies[index].Name == name {
			return &result.Dependencies[index]
		}
	}
	return nil
}

func Check(runner CommandRunner) CheckResult {
	definitions := []struct {
		name     string
		required bool
	}{
		{name: "ansible", required: true},
		{name: "ansible-playbook", required: true},
		{name: "ssh", required: true},
		{name: "ssh-keygen", required: true},
		{name: "kubectl", required: false},
		{name: "make", required: false},
		{name: "ansible-lint", required: false},
		{name: "yamllint", required: false},
	}

	result := CheckResult{Ready: true}
	for _, definition := range definitions {
		dependency := Dependency{Name: definition.name, Required: definition.required}
		path, err := runner.LookPath(definition.name)
		if err == nil {
			dependency.Path = path
			dependency.Installed = true
			version, versionErr := runner.Version(definition.name)
			if versionErr == nil {
				dependency.Version = firstVersionLine(version)
			}
		} else if definition.required {
			result.Ready = false
		}
		result.Dependencies = append(result.Dependencies, dependency)
	}
	return result
}

func firstVersionLine(version string) string {
	for _, line := range strings.Split(version, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
	}
	return ""
}

type SystemCommandRunner struct{}

func (SystemCommandRunner) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

func (SystemCommandRunner) Version(name string) (string, error) {
	output, err := exec.Command(name, "--version").CombinedOutput()
	return string(output), err
}
