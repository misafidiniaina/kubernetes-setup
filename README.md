# Kubernetes Cluster Bootstrap with Ansible

Quickly deploy a Kubernetes cluster on existing Linux machines without configuring each server manually. The inventory describes the nodes accessible over SSH; Ansible prepares the systems, installs containerd and Kubernetes, initializes the control plane, joins the nodes, and installs the cluster network.

## What the Project Automates

- Ubuntu/Debian preparation: swap, kernel modules, sysctl, NTP, and firewall;
- installation of containerd, kubelet, kubeadm, and kubectl;
- initialization of the first control plane with kubeadm;
- addition of extra control planes and workers;
- installation of Calico, Flannel, or Weave;
- optional installation of NGINX Ingress, Metrics Server, and cert-manager;
- final validation of the nodes and cluster state;
- staged execution through Ansible tags.

The project targets reproducible deployments on existing infrastructure. It does not provision machines and does not replace a complete production architecture with an external load balancer, backups, and centralized secret management.

## Prerequisites

- Ansible 2.9 or later;
- Ubuntu/Debian nodes with at least 2 CPUs and 2 GB of RAM;
- SSH access from the machine that runs Ansible;
- a user with `sudo` privileges;
- network connectivity between the nodes;
- a single Kubernetes version defined in `group_vars/all.yml`.

## Quick Deployment

```bash
cp inventory/hosts.ini.example inventory/hosts.ini
# Adjust the addresses, user, and SSH key
vim inventory/hosts.ini
vim group_vars/all.yml

make ping
make syntax-check
make setup
make validate
```

## Verification CLI

The first Go module checks the control machine's dependencies before running Ansible:

```bash
go run ./cli/cmd/kube-bootstrap check
```

The required commands are `ansible`, `ansible-playbook`, `ssh`, and `ssh-keygen`. `kubectl` is optional and is used for local cluster validation. A missing required tool blocks deployment and points to the `kube-bootstrap install` module.

The command exits with status `0` when all required tools are installed, `1` when a required tool is missing, and `2` for invalid usage. Optional tools are reported but do not block deployment.

To run the checks directly with Go:

```bash
go -C cli test ./...
go -C cli run ./cmd/kube-bootstrap check
```

The installer detects the Linux distribution from `/etc/os-release` and selects a supported package manager: `apt-get` for Debian-based systems, `dnf` or `yum` for Red Hat-based systems, and `pacman` for Arch-based systems. It prints an installation plan by default; pass `--apply` to execute it with `sudo`:

```bash
go -C cli run ./cmd/kube-bootstrap install
go -C cli run ./cmd/kube-bootstrap install --apply
```

The `check` module is tested with a mock executor, so the tests do not depend on tools installed on the development machine. The scenarios in `cli/features/` are executed with Godog and use fakes for package installation, so CI never changes the runner system. The CLI code is isolated in the `cli/` directory, separately from the Ansible playbooks and roles.

## Continuous Integration

GitHub Actions runs the CLI tests with the race detector, `go vet`, and a complete build on every pull request and push to `main`. See [CONTRIBUTING.md](CONTRIBUTING.md) for the workflow between feature scenarios, contributor tests, and CI.

The complete playbook is available in `playbooks/site.yml`. To run it without `make`:

```bash
ansible-playbook -i inventory/hosts.ini playbooks/site.yml
```

## Expected Output Example

```text
NAME      STATUS   ROLES           VERSION
master-1  Ready    control-plane   v1.28.x
worker-1  Ready    <none>          v1.28.x
worker-2  Ready    <none>          v1.28.x
```

## Configuration

The main variables are defined in `group_vars/all.yml`:

| Variable                | Purpose                          |
| ----------------------- | -------------------------------- |
| `kubernetes_version`    | Installed Kubernetes version     |
| `container_runtime`     | `containerd` or `docker`         |
| `networking_plugin`     | `calico`, `flannel`, or `weave`  |
| `pod_network_cidr`      | CIDR used by the pod network     |
| `service_cidr`          | CIDR used by Kubernetes services |
| `enable_ingress_nginx`  | Enables the Ingress Controller   |
| `enable_metrics_server` | Enables Metrics Server           |
| `enable_firewall`       | Configures UFW on the nodes      |

Join tokens and certificates are sensitive data. They must not be committed to Git or displayed in CI logs.

## Staged Execution

```bash
make prerequisites
make init
make join
make networking
make validate
```

The same stages are available with the `prerequisites`, `init`, `join`, `workers`, `control-plane`, `networking`, and `finalize` tags.

## Project Structure

```text
inventory/       Existing machine inventory
group_vars/      Shared cluster configuration
playbooks/       Global orchestration
roles/common/    OS preparation and Kubernetes components
roles/kubeadm-init/  First control plane initialization
roles/kubeadm-join/  Additional node joining
roles/networking/   Network add-ons and components
```

## Validation and Troubleshooting

```bash
make validate
kubectl get nodes -o wide
kubectl get pods -A
journalctl -u kubelet -f
journalctl -u containerd -f
```

To change the cluster size, add or remove hosts from the inventory and rerun the playbook. The bootstrap tasks are designed to be rerun without resetting an already configured node.

## Demonstrated Skills

Ansible, Linux, systemd, SSH, containerd, kubeadm, Kubernetes, networking, firewall, idempotence, automation, and infrastructure validation.
