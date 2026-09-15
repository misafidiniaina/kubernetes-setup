# Kubernetes Cluster Bootstrap with Ansible

Déployer rapidement un cluster Kubernetes sur des machines Linux existantes, sans configurer chaque serveur manuellement. L'inventaire décrit les nœuds accessibles en SSH ; Ansible prépare les systèmes, installe containerd et Kubernetes, initialise le control plane, rejoint les nœuds et installe le réseau du cluster.

## Ce que le projet automatise

- préparation Ubuntu/Debian : swap, modules kernel, sysctl, NTP et firewall ;
- installation de containerd, kubelet, kubeadm et kubectl ;
- initialisation du premier control plane avec kubeadm ;
- ajout de control planes supplémentaires et de workers ;
- installation de Calico, Flannel ou Weave ;
- installation optionnelle d'Ingress NGINX, Metrics Server et cert-manager ;
- validation finale des nœuds et de l'état du cluster ;
- exécution par étapes grâce aux tags Ansible.

Le projet vise un déploiement reproductible sur une infrastructure existante. Il ne provisionne pas les machines et ne remplace pas une architecture de production complète avec load balancer externe, sauvegardes et gestion centralisée des secrets.

## Prérequis

- Ansible 2.9 ou une version plus récente ;
- nœuds Ubuntu/Debian avec au moins 2 CPU et 2 Go de RAM ;
- accès SSH depuis la machine qui exécute Ansible ;
- un utilisateur avec privilèges `sudo` ;
- connectivité réseau entre les nœuds ;
- une seule version de Kubernetes définie dans `group_vars/all.yml`.

## Déploiement rapide

```bash
cp inventory/hosts.ini.example inventory/hosts.ini
# Adapter les adresses, l'utilisateur et la clé SSH
vim inventory/hosts.ini
vim group_vars/all.yml

make ping
make syntax-check
make setup
make validate
```

## CLI de vérification

Le premier module Go vérifie les dépendances de la machine de contrôle avant de lancer Ansible :

```bash
go run ./cmd/kube-bootstrap check
```

Les commandes obligatoires sont `ansible`, `ansible-playbook`, `ssh` et `ssh-keygen`. `kubectl` est optionnel et sert à la validation locale du cluster. Un outil obligatoire absent bloque le déploiement et indique le futur module `kube-bootstrap install`.

Pour exécuter le même contrôle avec le Makefile :

```bash
make test
make check
```

Le module `check` est testé avec un exécuteur simulé ; les tests ne dépendent donc pas des outils installés sur la machine de développement.

Le playbook complet est disponible dans `playbooks/site.yml`. Pour une exécution sans `make` :

```bash
ansible-playbook -i inventory/hosts.ini playbooks/site.yml
```

## Exemple de résultat attendu

```text
NAME      STATUS   ROLES           VERSION
master-1  Ready    control-plane   v1.28.x
worker-1  Ready    <none>          v1.28.x
worker-2  Ready    <none>          v1.28.x
```

## Configuration

Les variables principales sont définies dans `group_vars/all.yml` :

| Variable                | Rôle                                     |
| ----------------------- | ---------------------------------------- |
| `kubernetes_version`    | Version installée de Kubernetes          |
| `container_runtime`     | `containerd` ou `docker`                 |
| `networking_plugin`     | `calico`, `flannel` ou `weave`           |
| `pod_network_cidr`      | CIDR utilisé par le réseau des pods      |
| `service_cidr`          | CIDR utilisé par les services Kubernetes |
| `enable_ingress_nginx`  | Active l'Ingress Controller              |
| `enable_metrics_server` | Active Metrics Server                    |
| `enable_firewall`       | Configure UFW sur les nœuds              |

Les tokens et certificats de join sont des données sensibles. Ils ne doivent pas être commités dans Git ni affichés dans les logs CI.

## Exécution par étapes

```bash
make prerequisites
make init
make join
make networking
make validate
```

Les mêmes étapes sont disponibles avec les tags `prerequisites`, `init`, `join`, `workers`, `control-plane`, `networking` et `finalize`.

## Structure

```text
inventory/       Inventaire des machines existantes
group_vars/      Configuration commune du cluster
playbooks/       Orchestration globale
roles/common/    Préparation OS et composants Kubernetes
roles/kubeadm-init/  Initialisation du premier control plane
roles/kubeadm-join/  Ajout des autres nœuds
roles/networking/   Add-ons et composants réseau
```

## Validation et dépannage

```bash
make validate
kubectl get nodes -o wide
kubectl get pods -A
journalctl -u kubelet -f
journalctl -u containerd -f
```

Pour modifier la taille du cluster, ajouter ou supprimer des hôtes dans l'inventaire puis relancer le playbook. Les tâches de bootstrap sont conçues pour être rejouées sans réinitialiser un nœud déjà configuré.

## Compétences démontrées

Ansible, Linux, systemd, SSH, containerd, kubeadm, Kubernetes, réseau, firewall, idempotence, automatisation et validation d'infrastructure.
