package checkcmd

import (
	"fmt"
	"io"

	"github.com/devops/kubernetes-setup/cli/internal/dependencies"
)

func Execute(runner dependencies.CommandRunner, output io.Writer) int {
	return Run(dependencies.Check(runner), output)
}

func Run(result dependencies.CheckResult, output io.Writer) int {
	fmt.Fprintln(output, "Dependency check")
	for _, dependency := range result.Dependencies {
		status := "missing"
		if dependency.Installed {
			status = "installed"
		}
		version := dependency.Version
		if version == "" {
			version = "version unavailable"
		}
		marker := "!"
		if dependency.Installed {
			marker = "✓"
		}
		fmt.Fprintf(output, "%s %-17s %-9s %s\n", marker, dependency.Name, status, version)
	}

	if result.Ready {
		fmt.Fprintln(output, "\nEnvironment ready.")
		return 0
	}

	fmt.Fprintln(output, "\nEnvironment incomplete.")
	fmt.Fprintln(output, "Run 'kube-bootstrap install' to install missing dependencies.")
	return 1
}
