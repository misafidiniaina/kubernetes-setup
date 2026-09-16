Feature: Local tool verification

  Scenario: The deployment environment is complete
    Given ansible, ansible-playbook, ssh, and ssh-keygen are installed
    When the user runs "kube-bootstrap check"
    Then the environment is reported as ready

  Scenario: Optional validation tools are missing
    Given the required deployment tools are installed
    And make, kubectl, ansible-lint, and yamllint are missing
    When the user runs "kube-bootstrap check"
    Then the environment is still reported as ready
    And the optional tools are reported as missing

  Scenario: A required tool is missing
    Given ansible-playbook is missing
    When the user runs "kube-bootstrap check"
    Then the environment is reported as incomplete
    And dependency installation is recommended