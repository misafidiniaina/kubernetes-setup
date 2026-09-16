package packages

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type CommandRunner interface {
	Run(name string, args ...string) error
}

type Manager interface {
	Name() string
	Install(packageNames []string) error
}

type OSInfo struct {
	ID        string
	Name      string
	VersionID string
	IDLike    []string
}

type SystemCommandRunner struct{}

func (SystemCommandRunner) Run(name string, args ...string) error {
	return exec.Command(name, args...).Run()
}

type AptManager struct {
	Runner CommandRunner
}

func (manager AptManager) Name() string {
	return "apt-get"
}

func (manager AptManager) Install(packageNames []string) error {
	if len(packageNames) == 0 {
		return nil
	}
	args := append([]string{"apt-get", "install", "-y"}, packageNames...)
	return manager.Runner.Run("sudo", args...)
}

func Detect(runner CommandLookup, commandRunner CommandRunner) (Manager, error) {
	osInfo, err := ReadOSRelease("/etc/os-release")
	if err != nil {
		return nil, fmt.Errorf("detect operating system: %w", err)
	}
	return DetectForOS(osInfo, runner, commandRunner)
}

func DetectForOS(osInfo OSInfo, runner CommandLookup, commandRunner CommandRunner) (Manager, error) {
	managerCommands := packageManagerCandidates(osInfo)
	for _, command := range managerCommands {
		if _, err := runner.LookPath(command); err != nil {
			continue
		}
		switch command {
		case "apt-get":
			return AptManager{Runner: commandRunner}, nil
		case "dnf":
			return DnfManager{Runner: commandRunner}, nil
		case "yum":
			return YumManager{Runner: commandRunner}, nil
		case "pacman":
			return PacmanManager{Runner: commandRunner}, nil
		}
	}
	return nil, fmt.Errorf("no supported package manager found for %s (supported: apt-get, dnf, yum, pacman)", osInfo.DisplayName())
}

func packageManagerCandidates(osInfo OSInfo) []string {
	family := append([]string{osInfo.ID}, osInfo.IDLike...)
	for _, id := range family {
		switch strings.ToLower(id) {
		case "debian", "ubuntu", "kali", "linuxmint":
			return []string{"apt-get"}
		case "fedora", "rhel", "centos", "rocky", "almalinux":
			return []string{"dnf", "yum"}
		case "arch", "manjaro":
			return []string{"pacman"}
		}
	}
	return []string{"apt-get", "dnf", "yum", "pacman"}
}

func ReadOSRelease(path string) (OSInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return OSInfo{}, err
	}
	values := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok || key == "" || strings.HasPrefix(key, "#") {
			continue
		}
		values[key] = strings.Trim(strings.TrimSpace(value), `"'`)
	}
	var idLike []string
	if values["ID_LIKE"] != "" {
		idLike = strings.Fields(values["ID_LIKE"])
	}
	return OSInfo{ID: values["ID"], Name: values["NAME"], VersionID: values["VERSION_ID"], IDLike: idLike}, nil
}

func (info OSInfo) DisplayName() string {
	name := info.Name
	if name == "" {
		name = info.ID
	}
	if info.VersionID != "" {
		return name + " " + info.VersionID
	}
	return name
}

type DnfManager struct{ Runner CommandRunner }

func (manager DnfManager) Name() string { return "dnf" }

func (manager DnfManager) Install(packageNames []string) error {
	return manager.Runner.Run("sudo", append([]string{"dnf", "install", "-y"}, packageNames...)...)
}

type YumManager struct{ Runner CommandRunner }

func (manager YumManager) Name() string { return "yum" }

func (manager YumManager) Install(packageNames []string) error {
	return manager.Runner.Run("sudo", append([]string{"yum", "install", "-y"}, packageNames...)...)
}

type PacmanManager struct{ Runner CommandRunner }

func (manager PacmanManager) Name() string { return "pacman" }

func (manager PacmanManager) Install(packageNames []string) error {
	return manager.Runner.Run("sudo", append([]string{"pacman", "-S", "--noconfirm"}, packageNames...)...)
}

type CommandLookup interface {
	LookPath(name string) (string, error)
}

func PackageNames(managerName string, dependencies []string) []string {
	sshPackage := map[string]string{
		"apt-get": "openssh-client",
		"dnf":     "openssh-clients",
		"yum":     "openssh-clients",
		"pacman":  "openssh",
	}[managerName]
	knownPackages := map[string]string{
		"ansible":          "ansible",
		"ansible-playbook": "ansible",
		"ssh":              sshPackage,
		"ssh-keygen":       sshPackage,
	}
	seen := make(map[string]struct{}, len(dependencies))
	result := make([]string, 0, len(dependencies))
	for _, dependency := range dependencies {
		packageName, ok := knownPackages[dependency]
		if !ok || packageName == "" {
			continue
		}
		if _, ok := seen[packageName]; ok {
			continue
		}
		seen[packageName] = struct{}{}
		result = append(result, packageName)
	}
	return result
}

func FormatCommand(manager Manager, packageNames []string) string {
	if len(packageNames) == 0 {
		return ""
	}
	if manager.Name() == "apt-get" {
		return "sudo apt-get install -y " + strings.Join(packageNames, " ")
	}
	return manager.Name() + " install " + strings.Join(packageNames, " ")
}
