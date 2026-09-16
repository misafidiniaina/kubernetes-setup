package packages

import (
	"errors"
	"reflect"
	"testing"
)

type fakeLookup struct {
	available bool
}

func (lookup fakeLookup) LookPath(string) (string, error) {
	if lookup.available {
		return "/usr/bin/apt-get", nil
	}
	return "", errors.New("not found")
}

type fakeCommandRunner struct {
	name string
	args []string
}

func (runner *fakeCommandRunner) Run(name string, args ...string) error {
	runner.name = name
	runner.args = args
	return nil
}

func TestPackageNamesDeduplicatesSharedPackages(t *testing.T) {
	got := PackageNames([]string{"ansible", "ansible-playbook", "ssh", "ssh-keygen"})
	want := []string{"ansible", "openssh-client"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected packages %v, got %v", want, got)
	}
}

func TestDetectReturnsAptManager(t *testing.T) {
	manager, err := Detect(fakeLookup{available: true}, &fakeCommandRunner{})
	if err != nil {
		t.Fatalf("expected manager, got error: %v", err)
	}
	if manager.Name() != "apt-get" {
		t.Fatalf("expected apt-get manager, got %q", manager.Name())
	}
}

func TestAptManagerInstallsWithSudo(t *testing.T) {
	runner := &fakeCommandRunner{}
	manager := AptManager{Runner: runner}

	if err := manager.Install([]string{"ansible", "openssh-client"}); err != nil {
		t.Fatalf("unexpected install error: %v", err)
	}
	if runner.name != "sudo" {
		t.Fatalf("expected sudo command, got %q", runner.name)
	}
	want := []string{"apt-get", "install", "-y", "ansible", "openssh-client"}
	if !reflect.DeepEqual(runner.args, want) {
		t.Fatalf("expected args %v, got %v", want, runner.args)
	}
}
