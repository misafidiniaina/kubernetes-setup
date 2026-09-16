package main

import (
	"fmt"
	"io"
	"os"

	"github.com/devops/kubernetes-setup/cli/internal/checkcmd"
	"github.com/devops/kubernetes-setup/cli/internal/dependencies"
	"github.com/devops/kubernetes-setup/cli/internal/installcmd"
	"github.com/devops/kubernetes-setup/cli/internal/packages"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, dependencies.SystemCommandRunner{}))
}

func run(args []string, stdout, stderr io.Writer, runner dependencies.CommandRunner) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "check":
		if len(args) > 1 {
			fmt.Fprintln(stderr, "check does not accept arguments")
			printUsage(stderr)
			return 2
		}
		return checkcmd.Execute(runner, stdout)
	case "install":
		apply, err := parseApply(args[1:])
		if err != nil {
			fmt.Fprintln(stderr, err)
			printUsage(stderr)
			return 2
		}
		factory := func() (packages.Manager, error) {
			return packages.Detect(runner, packages.SystemCommandRunner{})
		}
		return installcmd.Execute(runner, factory, apply, stdout)
	case "help", "--help", "-h":
		if len(args) > 1 {
			fmt.Fprintln(stderr, "help does not accept arguments")
			printUsage(stderr)
			return 2
		}
		printUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		printUsage(stderr)
		return 2
	}
}

func printUsage(output io.Writer) {
	fmt.Fprintln(output, "Usage: kube-bootstrap <command>")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Commands:")
	fmt.Fprintln(output, "  check    Check local deployment dependencies")
	fmt.Fprintln(output, "  install  Plan or install missing dependencies")
	fmt.Fprintln(output, "  help     Show this help")
}

func parseApply(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	if len(args) == 1 && args[0] == "--apply" {
		return true, nil
	}
	return false, fmt.Errorf("install accepts only --apply")
}
