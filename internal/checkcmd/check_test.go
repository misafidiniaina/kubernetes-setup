package checkcmd

import (
	"strings"
	"testing"

	"github.com/devops/kubernetes-setup/internal/dependencies"
)

func TestRunReportsReadyEnvironment(t *testing.T) {
	var output strings.Builder
	result := dependencies.CheckResult{
		Ready: true,
		Dependencies: []dependencies.Dependency{
			{Name: "ansible-playbook", Version: "ansible-playbook [core 2.16.5]", Installed: true, Required: true},
			{Name: "kubectl", Installed: false, Required: false},
		},
	}

	exitCode := Run(result, &output)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(output.String(), "Environment ready") {
		t.Fatalf("expected ready message, got %q", output.String())
	}
	if !strings.Contains(output.String(), "kubectl") {
		t.Fatalf("expected optional dependency in output, got %q", output.String())
	}
}

func TestRunReportsMissingRequiredDependency(t *testing.T) {
	var output strings.Builder
	result := dependencies.CheckResult{
		Ready: false,
		Dependencies: []dependencies.Dependency{
			{Name: "ansible-playbook", Installed: false, Required: true},
		},
	}

	exitCode := Run(result, &output)

	if exitCode == 0 {
		t.Fatal("expected a non-zero exit code")
	}
	if !strings.Contains(output.String(), "ansible-playbook") {
		t.Fatalf("expected missing dependency in output, got %q", output.String())
	}
	if !strings.Contains(output.String(), "kube-bootstrap install") {
		t.Fatalf("expected installation hint, got %q", output.String())
	}
}
