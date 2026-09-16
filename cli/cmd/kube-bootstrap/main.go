package main

import (
	"fmt"
	"io"
	"os"

	"github.com/devops/kubernetes-setup/cli/internal/checkcmd"
	"github.com/devops/kubernetes-setup/cli/internal/dependencies"
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
	fmt.Fprintln(output, "  help     Show this help")
}
