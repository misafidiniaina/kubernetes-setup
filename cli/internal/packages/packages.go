package packages

import (
	"fmt"
	"os/exec"
	"strings"
)

type CommandRunner interface {
	Run(name string, args ...string) error
}

type Manager interface {
	Name() string
	Install(packageNames []string) error
}

type SystemCommandRunner struct{}

func (SystemCommandRunner) Run(name string, args ...string) error {
	return exec.Command(name, args...).Run()
}

type AptManager struct {
	Runner CommandRunner
}

func (manager AptManager) Name() string {
	return "apt-get"
}

func (manager AptManager) Install(packageNames []string) error {
	if len(packageNames) == 0 {
		return nil
	}
	args := append([]string{"apt-get", "install", "-y"}, packageNames...)
	return manager.Runner.Run("sudo", args...)
}

func Detect(runner CommandLookup, commandRunner CommandRunner) (Manager, error) {
	if _, err := runner.LookPath("apt-get"); err == nil {
		return AptManager{Runner: commandRunner}, nil
	}
	return nil, fmt.Errorf("no supported package manager found (supported: apt-get)")
}

type CommandLookup interface {
	LookPath(name string) (string, error)
}

func PackageNames(dependencies []string) []string {
	knownPackages := map[string]string{
		"ansible":          "ansible",
		"ansible-playbook": "ansible",
		"ssh":              "openssh-client",
		"ssh-keygen":       "openssh-client",
	}
	seen := make(map[string]struct{}, len(dependencies))
	result := make([]string, 0, len(dependencies))
	for _, dependency := range dependencies {
		packageName, ok := knownPackages[dependency]
		if !ok {
			continue
		}
		if _, ok := seen[packageName]; ok {
			continue
		}
		seen[packageName] = struct{}{}
		result = append(result, packageName)
	}
	return result
}

func FormatCommand(manager Manager, packageNames []string) string {
	if len(packageNames) == 0 {
		return ""
	}
	if manager.Name() == "apt-get" {
		return "sudo apt-get install -y " + strings.Join(packageNames, " ")
	}
	return manager.Name() + " install " + strings.Join(packageNames, " ")
}
