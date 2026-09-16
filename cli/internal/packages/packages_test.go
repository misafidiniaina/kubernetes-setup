package packages

import (
	"errors"
	"os"
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
	got := PackageNames("apt-get", []string{"ansible", "ansible-playbook", "ssh", "ssh-keygen"})
	want := []string{"ansible", "openssh-client"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected packages %v, got %v", want, got)
	}
}

func TestPackageNamesUseManagerSpecificSSHPackage(t *testing.T) {
	got := PackageNames("dnf", []string{"ssh", "ssh-keygen"})
	want := []string{"openssh-clients"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected packages %v, got %v", want, got)
	}
}

func TestReadOSRelease(t *testing.T) {
	file := t.TempDir() + "/os-release"
	content := "NAME=\"Kali GNU/Linux\"\nID=kali\nID_LIKE=debian\nVERSION_ID=\"2026.1\"\n"
	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatalf("write os-release fixture: %v", err)
	}

	info, err := ReadOSRelease(file)
	if err != nil {
		t.Fatalf("read os-release: %v", err)
	}
	if info.DisplayName() != "Kali GNU/Linux 2026.1" {
		t.Fatalf("unexpected OS name: %q", info.DisplayName())
	}
	if !reflect.DeepEqual(info.IDLike, []string{"debian"}) {
		t.Fatalf("unexpected OS family: %v", info.IDLike)
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
