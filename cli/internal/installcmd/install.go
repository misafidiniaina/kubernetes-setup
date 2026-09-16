package installcmd

import (
	"fmt"
	"io"

	"github.com/devops/kubernetes-setup/cli/internal/dependencies"
	"github.com/devops/kubernetes-setup/cli/internal/packages"
)

type ManagerFactory func() (packages.Manager, error)

func Execute(runner dependencies.CommandRunner, factory ManagerFactory, apply bool, output io.Writer) int {
	result := dependencies.Check(runner)
	missing := result.MissingRequired()
	if len(missing) == 0 {
		fmt.Fprintln(output, "All required dependencies are already installed.")
		return 0
	}

	missingNames := make([]string, 0, len(missing))
	for _, dependency := range missing {
		missingNames = append(missingNames, dependency.Name)
	}
	manager, err := factory()
	if err != nil {
		fmt.Fprintf(output, "Unable to prepare installation: %v\n", err)
		return 1
	}
	packageNames := packages.PackageNames(manager.Name(), missingNames)
	command := packages.FormatCommand(manager, packageNames)
	if !apply {
		fmt.Fprintf(output, "Installation plan: %s\n", command)
		fmt.Fprintln(output, "Re-run with --apply to install the missing dependencies.")
		return 0
	}
	if err := manager.Install(packageNames); err != nil {
		fmt.Fprintf(output, "Installation failed: %v\n", err)
		return 1
	}
	fmt.Fprintln(output, "Dependencies installed successfully.")
	return 0
}
