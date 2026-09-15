Feature: Vérification des outils locaux

  Scenario: L'environnement de déploiement est complet
    Given ansible, ansible-playbook, ssh et ssh-keygen sont installés
    When l'utilisateur lance "kube-bootstrap check"
    Then l'environnement est déclaré prêt

  Scenario: Les outils optionnels de validation sont absents
    Given les outils de déploiement obligatoires sont installés
    And make, kubectl, ansible-lint et yamllint sont absents
    When l'utilisateur lance "kube-bootstrap check"
    Then l'environnement reste déclaré prêt
    And les outils optionnels sont signalés comme manquants

  Scenario: Un outil obligatoire est absent
    Given ansible-playbook est absent
    When l'utilisateur lance "kube-bootstrap check"
    Then l'environnement est déclaré incomplet
    And l'installation des dépendances est recommandée