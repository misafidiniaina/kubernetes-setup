package installcmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/devops/kubernetes-setup/cli/internal/dependencies"
	"github.com/devops/kubernetes-setup/cli/internal/packages"
)

type fakeDependencyRunner struct {
	installed map[string]bool
}

func (runner fakeDependencyRunner) LookPath(name string) (string, error) {
	if runner.installed[name] {
		return "/usr/bin/" + name, nil
	}
	return "", errors.New("command not found")
}

func (runner fakeDependencyRunner) Version(string) (string, error) {
	return "tool 1.0.0", nil
}

type fakePackageManager struct {
	installed []string
	err       error
}

func (manager *fakePackageManager) Name() string {
	return "apt-get"
}

func (manager *fakePackageManager) Install(packageNames []string) error {
	manager.installed = append(manager.installed, packageNames...)
	return manager.err
}

func TestExecutePrintsPlanWithoutInstalling(t *testing.T) {
	manager := &fakePackageManager{}
	var output strings.Builder

	exitCode := Execute(
		fakeDependencyRunner{installed: map[string]bool{"ssh": true, "ssh-keygen": true}},
		func() (packages.Manager, error) { return manager, nil },
		false,
		&output,
	)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if len(manager.installed) != 0 {
		t.Fatalf("expected no installation, got %v", manager.installed)
	}
	if !strings.Contains(output.String(), "sudo apt-get install -y ansible") {
		t.Fatalf("expected installation plan, got %q", output.String())
	}
}

func TestExecuteInstallsMissingPackages(t *testing.T) {
	manager := &fakePackageManager{}
	var output strings.Builder

	exitCode := Execute(
		fakeDependencyRunner{installed: map[string]bool{"ssh": true, "ssh-keygen": true}},
		func() (packages.Manager, error) { return manager, nil },
		true,
		&output,
	)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if len(manager.installed) != 1 || manager.installed[0] != "ansible" {
		t.Fatalf("expected ansible installation, got %v", manager.installed)
	}
	if !strings.Contains(output.String(), "Dependencies installed successfully.") {
		t.Fatalf("expected success message, got %q", output.String())
	}
}

func TestExecuteReportsInstallerFailure(t *testing.T) {
	manager := &fakePackageManager{err: errors.New("permission denied")}
	var output strings.Builder

	exitCode := Execute(
		fakeDependencyRunner{installed: map[string]bool{"ssh": true, "ssh-keygen": true}},
		func() (packages.Manager, error) { return manager, nil },
		true,
		&output,
	)

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(output.String(), "permission denied") {
		t.Fatalf("expected installer error, got %q", output.String())
	}
}

var _ dependencies.CommandRunner = fakeDependencyRunner{}
