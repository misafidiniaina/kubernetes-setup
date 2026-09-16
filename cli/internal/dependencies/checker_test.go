package dependencies

import (
	"errors"
	"testing"
)

type fakeCommandRunner struct {
	paths    map[string]string
	versions map[string]string
	errors   map[string]error
}

func (runner fakeCommandRunner) LookPath(name string) (string, error) {
	if err := runner.errors[name]; err != nil {
		return "", err
	}
	path, ok := runner.paths[name]
	if !ok {
		return "", errors.New("command not found")
	}
	return path, nil
}

func (runner fakeCommandRunner) Version(name string) (string, error) {
	if err := runner.errors[name]; err != nil {
		return "", err
	}
	version, ok := runner.versions[name]
	if !ok {
		return "", errors.New("version unavailable")
	}
	return version, nil
}

func completeRunner() fakeCommandRunner {
	return fakeCommandRunner{
		paths: map[string]string{
			"ansible":          "/usr/bin/ansible",
			"ansible-playbook": "/usr/bin/ansible-playbook",
			"ssh":              "/usr/bin/ssh",
			"ssh-keygen":       "/usr/bin/ssh-keygen",
			"kubectl":          "/usr/bin/kubectl",
			"make":             "/usr/bin/make",
			"ansible-lint":     "/usr/bin/ansible-lint",
			"yamllint":         "/usr/bin/yamllint",
		},
		versions: map[string]string{
			"ansible":          "ansible [core 2.16.5]",
			"ansible-playbook": "ansible-playbook [core 2.16.5]",
			"ssh":              "OpenSSH_9.6",
			"ssh-keygen":       "OpenSSH_9.6",
			"kubectl":          "Client Version: v1.30.0",
			"make":             "GNU Make 4.3",
			"ansible-lint":     "ansible-lint 24.2.0",
			"yamllint":         "yamllint 1.35.1",
		},
	}
}

func TestCheckDependenciesReportsReadyWhenRequiredToolsAreInstalled(t *testing.T) {
	result := Check(completeRunner())

	if !result.Ready {
		t.Fatalf("expected environment to be ready, got: %+v", result)
	}
	if len(result.Dependencies) != 8 {
		t.Fatalf("expected 8 dependencies, got %d", len(result.Dependencies))
	}
}

func TestCheckDependenciesReportsMissingRequiredTool(t *testing.T) {
	runner := completeRunner()
	delete(runner.paths, "ansible-playbook")

	result := Check(runner)

	if result.Ready {
		t.Fatal("expected environment to be unavailable")
	}
	dependency := result.Dependency("ansible-playbook")
	if dependency == nil || dependency.Installed {
		t.Fatalf("expected ansible-playbook to be reported as missing, got: %+v", dependency)
	}
	if dependency.Required != true {
		t.Fatal("ansible-playbook should be required")
	}
}

func TestCheckDependenciesAllowsMissingOptionalTool(t *testing.T) {
	runner := completeRunner()
	delete(runner.paths, "kubectl")
	delete(runner.paths, "make")
	delete(runner.paths, "ansible-lint")
	delete(runner.paths, "yamllint")

	result := Check(runner)

	if !result.Ready {
		t.Fatal("kubectl should not block the environment")
	}
	dependency := result.Dependency("kubectl")
	if dependency == nil || dependency.Installed {
		t.Fatalf("expected kubectl to be reported as missing, got: %+v", dependency)
	}
	if dependency.Required {
		t.Fatal("kubectl should be optional")
	}

	for _, name := range []string{"make", "ansible-lint", "yamllint"} {
		dependency := result.Dependency(name)
		if dependency == nil || dependency.Installed || dependency.Required {
			t.Fatalf("expected %s to be an optional missing dependency, got: %+v", name, dependency)
		}
	}
}

func TestCheckDependenciesKeepsToolVersion(t *testing.T) {
	runner := completeRunner()
	runner.versions["ansible-playbook"] = "ansible-playbook [core 2.16.5]\nCopyright (C) 2024"

	result := Check(runner)
	dependency := result.Dependency("ansible-playbook")

	if dependency.Version != "ansible-playbook [core 2.16.5]" {
		t.Fatalf("unexpected version: %q", dependency.Version)
	}
}
