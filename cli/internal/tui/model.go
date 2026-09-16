package tui

import (
	"bytes"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/devops/kubernetes-setup/cli/internal/dependencies"
	"github.com/devops/kubernetes-setup/cli/internal/installcmd"
	"github.com/devops/kubernetes-setup/cli/internal/packages"
)

type screen int

const (
	checkScreen screen = iota
	installScreen
)

type checkCompletedMsg struct {
	result dependencies.CheckResult
}

type installCompletedMsg struct {
	exitCode int
	output   string
}

type Model struct {
	runner         dependencies.CommandRunner
	managerFactory installcmd.ManagerFactory
	result         dependencies.CheckResult
	output         string
	screen         screen
	loading        bool
	width          int
	height         int
}

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	readyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	missingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	dimStyle     = lipgloss.NewStyle().Faint(true)
)

func New(runner dependencies.CommandRunner, managerFactory installcmd.ManagerFactory) Model {
	return Model{
		runner:         runner,
		managerFactory: managerFactory,
		loading:        true,
	}
}

func DefaultManagerFactory(runner dependencies.CommandRunner) installcmd.ManagerFactory {
	return func() (packages.Manager, error) {
		return packages.Detect(runner, packages.SystemCommandRunner{})
	}
}

func (model Model) Init() tea.Cmd {
	return model.checkCmd()
}

func (model Model) checkCmd() tea.Cmd {
	return func() tea.Msg {
		return checkCompletedMsg{result: dependencies.Check(model.runner)}
	}
}

func (model Model) installCmd() tea.Cmd {
	return func() tea.Msg {
		var output bytes.Buffer
		exitCode := installcmd.Execute(model.runner, model.managerFactory, false, &output)
		return installCompletedMsg{exitCode: exitCode, output: output.String()}
	}
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = message.Width
		model.height = message.Height
	case checkCompletedMsg:
		model.result = message.result
		model.loading = false
		model.screen = checkScreen
	case installCompletedMsg:
		model.output = message.output
		model.loading = false
		model.screen = installScreen
	case tea.KeyMsg:
		switch message.String() {
		case "q", "ctrl+c":
			return model, tea.Quit
		case "r":
			model.loading = true
			model.output = ""
			model.screen = checkScreen
			return model, model.checkCmd()
		case "i":
			model.loading = true
			return model, model.installCmd()
		case "esc":
			model.output = ""
			model.screen = checkScreen
		}
	}
	return model, nil
}

func (model Model) View() string {
	if model.loading {
		return titleStyle.Render("kube-bootstrap") + "\n\nChecking local dependencies...\n\n" + footer()
	}
	if model.screen == installScreen {
		return titleStyle.Render("kube-bootstrap / install") + "\n\n" + model.output + "\n" + footer()
	}

	var builder strings.Builder
	builder.WriteString(titleStyle.Render("kube-bootstrap / dependencies"))
	builder.WriteString("\n\n")
	for _, dependency := range model.result.Dependencies {
		status := missingStyle.Render("missing")
		marker := "!"
		if dependency.Installed {
			status = readyStyle.Render("installed")
			marker = "✓"
		}
		version := dependency.Version
		if version == "" {
			version = "version unavailable"
		}
		builder.WriteString(fmt.Sprintf("%s %-17s %-12s %s\n", marker, dependency.Name, status, dimStyle.Render(version)))
	}
	builder.WriteString("\n")
	if model.result.Ready {
		builder.WriteString(readyStyle.Render("Environment ready."))
	} else {
		builder.WriteString(missingStyle.Render("Environment incomplete."))
	}
	builder.WriteString("\n\n")
	builder.WriteString(footer())
	return builder.String()
}

func footer() string {
	return dimStyle.Render("r refresh  i installation plan  q quit")
}
