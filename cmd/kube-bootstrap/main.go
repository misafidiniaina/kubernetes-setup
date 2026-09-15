package main

import (
	"fmt"
	"os"

	"github.com/devops/kubernetes-setup/internal/checkcmd"
	"github.com/devops/kubernetes-setup/internal/dependencies"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "check":
		os.Exit(checkcmd.Execute(dependencies.SystemCommandRunner{}, os.Stdout))
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Println("Usage: kube-bootstrap <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  check    Check local deployment dependencies")
	fmt.Println("  help     Show this help")
}
