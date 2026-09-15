# Makefile for Kubernetes Setup with Ansible

.PHONY: help setup prerequisites init join workers networking finalize validate check test all clean syntax-check ping lint inventory

# Default target
help:
	@echo "Available targets:"
	@echo "  setup          - Run complete cluster setup"
	@echo "  prerequisites  - Install prerequisites on all nodes"
	@echo "  init           - Initialize Kubernetes control plane"
	@echo "  join           - Join all nodes to cluster"
	@echo "  workers        - Join only worker nodes"
	@echo "  networking     - Install networking components"
	@echo "  finalize       - Final cluster validation"
	@echo "  validate       - Validate nodes and system pods"
	@echo "  check          - Check local CLI dependencies"
	@echo "  test           - Run Go tests"
	@echo "  clean          - Clean up temporary files"
	@echo "  syntax-check   - Check playbook syntax"
	@echo "  ping           - Test connectivity to all hosts"

# Variables
INVENTORY ?= inventory/hosts.ini
PLAYBOOK ?= playbooks/site.yml
ANSIBLE_OPTS ?=

# Main setup
setup: syntax-check
	ansible-playbook $(ANSIBLE_OPTS) -i $(INVENTORY) $(PLAYBOOK)

# Individual components
prerequisites:
	ansible-playbook $(ANSIBLE_OPTS) -i $(INVENTORY) $(PLAYBOOK) --tags prerequisites

init:
	ansible-playbook $(ANSIBLE_OPTS) -i $(INVENTORY) $(PLAYBOOK) --tags init

join:
	ansible-playbook $(ANSIBLE_OPTS) -i $(INVENTORY) $(PLAYBOOK) --tags join

workers:
	ansible-playbook $(ANSIBLE_OPTS) -i $(INVENTORY) $(PLAYBOOK) --tags workers

networking:
	ansible-playbook $(ANSIBLE_OPTS) -i $(INVENTORY) $(PLAYBOOK) --tags networking

finalize:
	ansible-playbook $(ANSIBLE_OPTS) -i $(INVENTORY) $(PLAYBOOK) --tags finalize

validate: syntax-check
	ansible-playbook $(ANSIBLE_OPTS) -i $(INVENTORY) $(PLAYBOOK) --tags finalize

check:
	go run ./cmd/kube-bootstrap check

test:
	go test ./...

# Utilities
syntax-check:
	ansible-playbook --syntax-check -i $(INVENTORY) $(PLAYBOOK)

ping:
	ansible -i $(INVENTORY) all -m ping

clean:
	find . -name "*.retry" -delete
	find . -name "*.log" -delete

# Development
lint:
	ansible-lint $(PLAYBOOK)
	yamllint .

# Show inventory
inventory:
	ansible-inventory -i $(INVENTORY) --list