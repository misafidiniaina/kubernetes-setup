package features

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/devops/kubernetes-setup/cli/internal/checkcmd"
	"github.com/devops/kubernetes-setup/cli/internal/dependencies"
	"github.com/devops/kubernetes-setup/cli/internal/installcmd"
	"github.com/devops/kubernetes-setup/cli/internal/packages"
)

type scenarioState struct {
	installed   map[string]bool
	manager     *fakePackageManager
	output      strings.Builder
	exitCode    int
	commandName string
}

type fakeCommandRunner struct {
	installed map[string]bool
}

func (runner fakeCommandRunner) LookPath(name string) (string, error) {
	if runner.installed[name] {
		return "/usr/bin/" + name, nil
	}
	return "", errors.New("command not found")
}

func (runner fakeCommandRunner) Version(string) (string, error) {
	return "tool 1.0.0", nil
}

type fakePackageManager struct {
	installed []string
}

func (manager *fakePackageManager) Name() string {
	return "apt-get"
}

func (manager *fakePackageManager) Install(packageNames []string) error {
	manager.installed = append(manager.installed, packageNames...)
	return nil
}

func TestFeatureScenarios(t *testing.T) {
	suite := godog.TestSuite{
		Name:                 "dependency features",
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: InitializeTestSuite,
		Options: &godog.Options{
			Format: "progress",
			Paths:  []string{"dependencies.feature"},
		},
	}
	if suite.Run() != 0 {
		t.Fatal("feature scenarios failed")
	}
}

func InitializeTestSuite(*godog.TestSuiteContext) {}

func InitializeScenario(ctx *godog.ScenarioContext) {
	state := &scenarioState{}
	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		state.installed = make(map[string]bool)
		state.manager = &fakePackageManager{}
		state.output.Reset()
		state.exitCode = 0
		state.commandName = ""
		return ctx, nil
	})

	ctx.Given(`^ansible, ansible-playbook, ssh, and ssh-keygen are installed$`, func() error {
		for _, name := range requiredDependencies() {
			state.installed[name] = true
		}
		return nil
	})
	ctx.Given(`^the required deployment tools are installed$`, func() error {
		for _, name := range requiredDependencies() {
			state.installed[name] = true
		}
		return nil
	})
	ctx.Given(`^make, kubectl, ansible-lint, and yamllint are missing$`, func() error {
		return nil
	})
	ctx.Given(`^ansible-playbook is missing$`, func() error {
		state.installed["ansible"] = true
		state.installed["ssh"] = true
		state.installed["ssh-keygen"] = true
		return nil
	})
	ctx.Given(`^ansible and ansible-playbook are missing$`, func() error {
		state.installed["ssh"] = true
		state.installed["ssh-keygen"] = true
		return nil
	})
	ctx.Given(`^ssh and ssh-keygen are installed$`, func() error {
		state.installed["ssh"] = true
		state.installed["ssh-keygen"] = true
		return nil
	})
	ctx.When(`^the user runs "kube-bootstrap (check|install|install --apply)"$`, func(command string) error {
		state.commandName = command
		runner := fakeCommandRunner{installed: state.installed}
		switch command {
		case "check":
			state.exitCode = checkcmd.Execute(runner, &state.output)
		case "install":
			state.exitCode = installcmd.Execute(runner, func() (packages.Manager, error) {
				return state.manager, nil
			}, false, &state.output)
		case "install --apply":
			state.exitCode = installcmd.Execute(runner, func() (packages.Manager, error) {
				return state.manager, nil
			}, true, &state.output)
		default:
			return fmt.Errorf("unsupported command %q", command)
		}
		return nil
	})
	ctx.Then(`^the environment is reported as ready$`, func() error {
		return assertOutput(state, 0, "Environment ready.")
	})
	ctx.Then(`^the environment is still reported as ready$`, func() error {
		return assertOutput(state, 0, "Environment ready.")
	})
	ctx.Then(`^the optional tools are reported as missing$`, func() error {
		for _, name := range []string{"make", "kubectl", "ansible-lint", "yamllint"} {
			if !strings.Contains(state.output.String(), name+" ") {
				return fmt.Errorf("optional dependency %q is absent from output %q", name, state.output.String())
			}
		}
		return nil
	})
	ctx.Then(`^the environment is reported as incomplete$`, func() error {
		return assertOutput(state, 1, "Environment incomplete.")
	})
	ctx.Then(`^dependency installation is recommended$`, func() error {
		return assertOutput(state, state.exitCode, "kube-bootstrap install")
	})
	ctx.Then(`^an apt installation plan for ansible is displayed$`, func() error {
		return assertOutput(state, 0, "Installation plan: sudo apt-get install -y ansible")
	})
	ctx.Then(`^no package is installed$`, func() error {
		if len(state.manager.installed) != 0 {
			return fmt.Errorf("expected no installed packages, got %v", state.manager.installed)
		}
		return nil
	})
	ctx.Then(`^ansible is installed through the detected package manager$`, func() error {
		if len(state.manager.installed) != 1 || state.manager.installed[0] != "ansible" {
			return fmt.Errorf("expected ansible to be installed, got %v", state.manager.installed)
		}
		return nil
	})
}

func requiredDependencies() []string {
	return []string{"ansible", "ansible-playbook", "ssh", "ssh-keygen"}
}

func assertOutput(state *scenarioState, expectedExitCode int, expectedText string) error {
	if state.exitCode != expectedExitCode {
		return fmt.Errorf("expected exit code %d, got %d; output: %q", expectedExitCode, state.exitCode, state.output.String())
	}
	if !strings.Contains(state.output.String(), expectedText) {
		return fmt.Errorf("expected %q in output %q", expectedText, state.output.String())
	}
	return nil
}

var _ dependencies.CommandRunner = fakeCommandRunner{}
