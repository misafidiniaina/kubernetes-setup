package main

import (
	"strings"
	"testing"
)

type stubRunner struct{}

func (stubRunner) LookPath(name string) (string, error) {
	return "/usr/bin/" + name, nil
}

func (stubRunner) Version(name string) (string, error) {
	return name + " 1.0.0", nil
}

func TestRunDispatchesCheckCommand(t *testing.T) {
	var stdout, stderr strings.Builder

	exitCode := run([]string{"check"}, &stdout, &stderr, stubRunner{})

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if !strings.Contains(stdout.String(), "Environment ready.") {
		t.Fatalf("expected ready message, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
}

func TestRunRejectsMissingCommand(t *testing.T) {
	var stdout, stderr strings.Builder

	exitCode := run(nil, &stdout, &stderr, stubRunner{})

	if exitCode != 2 {
		t.Fatalf("expected usage exit code 2, got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "Usage: kube-bootstrap <command>") {
		t.Fatalf("expected usage message, got %q", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout output, got %q", stdout.String())
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	var stdout, stderr strings.Builder

	exitCode := run([]string{"unknown"}, &stdout, &stderr, stubRunner{})

	if exitCode != 2 {
		t.Fatalf("expected usage exit code 2, got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "unknown command: unknown") {
		t.Fatalf("expected unknown command message, got %q", stderr.String())
	}
}

func TestRunRejectsUnexpectedArguments(t *testing.T) {
	var stdout, stderr strings.Builder

	exitCode := run([]string{"check", "extra"}, &stdout, &stderr, stubRunner{})

	if exitCode != 2 {
		t.Fatalf("expected usage exit code 2, got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "check does not accept arguments") {
		t.Fatalf("expected argument error, got %q", stderr.String())
	}
}
