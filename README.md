# Kubernetes Setup with Ansible and Kubeadm

This Ansible project provides a scalable and production-ready way to deploy Kubernetes clusters using kubeadm.

## Features

- **Scalable Architecture**: Support for multi-master HA clusters
- **Multiple Container Runtimes**: Containerd and Docker support
- **Networking Plugins**: Calico, Flannel, and Weave support
- **Security Best Practices**: Audit logging, RBAC, pod security standards
- **Add-ons**: Ingress Nginx, Metrics Server, Cert Manager
- **Idempotent**: Safe to run multiple times
- **Tagged Tasks**: Selective execution of components

## Prerequisites

- Ansible 2.9+
- Ubuntu/Debian servers
- SSH access to all nodes
- At least 2GB RAM and 2 CPUs per node
- Network connectivity between all nodes

## Quick Start

1. **Clone this repository**

   ```bash
   git clone <repository-url>
   cd kubernetes-setup
   ```

2. **Configure inventory**
   Edit `inventory/hosts.ini` with your server details:

   ```ini
   [masters]
   master1 ansible_host=192.168.1.10
   master2 ansible_host=192.168.1.11
   master3 ansible_host=192.168.1.12

   [workers]
   worker1 ansible_host=192.168.1.20
   worker2 ansible_host=192.168.1.21
   ```

3. **Configure variables**
   Edit `group_vars/all.yml` to customize your deployment

4. **Run the playbook**
   ```bash
   ansible-playbook -i inventory/hosts.ini playbooks/site.yml
   ```

## Project Structure

```
.
├── ansible.cfg                 # Ansible configuration
├── inventory/
│   └── hosts.ini              # Inventory file
├── group_vars/
│   └── all.yml                # Global variables
├── host_vars/                 # Host-specific variables
├── roles/
│   ├── common/                # Common setup (prerequisites, container runtime, k8s components)
│   ├── kubeadm-init/          # Control plane initialization
│   ├── kubeadm-join/          # Node joining
│   └── networking/            # Networking and add-ons
└── playbooks/
    └── site.yml               # Main playbook
```

## Configuration Options

### Cluster Configuration

- `kubernetes_version`: Kubernetes version to install
- `container_runtime`: `containerd` or `docker`
- `networking_plugin`: `calico`, `flannel`, or `weave`
- `pod_network_cidr`: Pod network CIDR
- `service_cidr`: Service network CIDR

### High Availability

- `load_balancer_address`: External load balancer IP for HA control plane
- Multiple masters in inventory for HA setup

### Security

- `enable_audit_log`: Enable Kubernetes audit logging
- `enable_pod_security_standards`: Enable pod security standards
- `disable_swap`: Disable swap (required for Kubernetes)

### Add-ons

- `enable_ingress_nginx`: Install NGINX Ingress Controller
- `enable_metrics_server`: Install Kubernetes Metrics Server
- `enable_cert_manager`: Install cert-manager

## Running Specific Tasks

Use tags to run specific parts of the deployment:

```bash
# Run only prerequisites
ansible-playbook -i inventory/hosts.ini playbooks/site.yml --tags prerequisites

# Initialize control plane only
ansible-playbook -i inventory/hosts.ini playbooks/site.yml --tags init

# Join workers only
ansible-playbook -i inventory/hosts.ini playbooks/site.yml --tags workers

# Install networking only
ansible-playbook -i inventory/hosts.ini playbooks/site.yml --tags networking
```

## Scaling the Cluster

### Adding Worker Nodes

1. Add new workers to `inventory/hosts.ini`
2. Run the join playbook:
   ```bash
   ansible-playbook -i inventory/hosts.ini playbooks/site.yml --tags workers
   ```

### Adding Control Plane Nodes

1. Add new masters to `inventory/hosts.ini`
2. Run the control plane join:
   ```bash
   ansible-playbook -i inventory/hosts.ini playbooks/site.yml --tags control-plane
   ```

## Troubleshooting

### Common Issues

1. **Kubelet fails to start**
   - Check container runtime is running
   - Verify cgroup driver configuration
   - Check system requirements

2. **Nodes fail to join**
   - Verify network connectivity
   - Check firewall rules
   - Ensure token is valid (tokens expire after 24 hours)

3. **Networking issues**
   - Verify pod network CIDR configuration
   - Check Calico/node status
   - Ensure required ports are open

### Logs and Debugging

```bash
# Check kubelet logs
journalctl -u kubelet -f

# Check container runtime logs
journalctl -u containerd -f  # or docker

# Check cluster status
kubectl get nodes
kubectl get pods -A
```

## Security Considerations

- Use SSH key authentication
- Restrict SSH access
- Configure firewall rules
- Enable audit logging
- Use RBAC for access control
- Keep Kubernetes and OS updated
- Use secrets management for sensitive data

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License.
