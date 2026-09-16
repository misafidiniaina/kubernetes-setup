package tui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/devops/kubernetes-setup/cli/internal/packages"
)

type fakeRunner struct {
	installed map[string]bool
}

func (runner fakeRunner) LookPath(name string) (string, error) {
	if runner.installed[name] {
		return "/usr/bin/" + name, nil
	}
	return "", errors.New("command not found")
}

func (runner fakeRunner) Version(name string) (string, error) {
	return name + " 1.0.0", nil
}

type fakeManager struct {
	installed []string
}

func (manager *fakeManager) Name() string {
	return "apt-get"
}

func (manager *fakeManager) Install(packageNames []string) error {
	manager.installed = append(manager.installed, packageNames...)
	return nil
}

func completeRunner() fakeRunner {
	return fakeRunner{installed: map[string]bool{
		"ansible":          true,
		"ansible-playbook": true,
		"ssh":              true,
		"ssh-keygen":       true,
		"kubectl":          true,
		"make":             true,
		"ansible-lint":     true,
		"yamllint":         true,
	}}
}

func TestModelLoadsAndRendersDependencyStatus(t *testing.T) {
	model := New(completeRunner(), func() (packages.Manager, error) {
		return &fakeManager{}, nil
	})

	updated, _ := model.Update(model.Init()())
	view := updated.(Model).View()

	if !strings.Contains(view, "kube-bootstrap / dependencies") {
		t.Fatalf("expected dependency title, got %q", view)
	}
	if !strings.Contains(view, "Environment ready.") {
		t.Fatalf("expected ready status, got %q", view)
	}
	if !strings.Contains(view, "ansible") {
		t.Fatalf("expected ansible dependency, got %q", view)
	}
}

func TestModelShowsInstallationPlan(t *testing.T) {
	manager := &fakeManager{}
	runner := fakeRunner{installed: map[string]bool{
		"ssh":        true,
		"ssh-keygen": true,
	}}
	model := New(runner, func() (packages.Manager, error) {
		return manager, nil
	})

	updated, _ := model.Update(model.Init()())
	model = updated.(Model)
	updated, command := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updated.(Model)
	if command == nil {
		t.Fatal("expected installation command")
	}
	updated, _ = model.Update(command())
	view := updated.(Model).View()

	if !strings.Contains(view, "sudo apt-get install -y ansible") {
		t.Fatalf("expected installation plan, got %q", view)
	}
	if len(manager.installed) != 0 {
		t.Fatalf("expected plan mode not to install packages, got %v", manager.installed)
	}
}

func TestModelQuitReturnsQuitCommand(t *testing.T) {
	model := New(completeRunner(), func() (packages.Manager, error) {
		return &fakeManager{}, nil
	})
	_, command := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if command == nil {
		t.Fatal("expected quit command")
	}
}
