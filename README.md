# gitops-bootstrap

> From zero to GitOps in one command — opinionated CLI to bootstrap a production-ready GitOps workflow on any Kubernetes cluster.

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Python 3.11+](https://img.shields.io/badge/python-3.11%2B-blue.svg)](https://www.python.org/)

---

## The Problem

Adopting GitOps is widely accepted as a best practice, but getting started is surprisingly painful:

- **Too many choices**: ArgoCD vs. Flux, Helm vs. Kustomize, Sealed Secrets vs. SOPS vs. External Secrets, mono-repo vs. multi-repo…
- **Hours of glue work**: Installing the controller, structuring the repo, wiring environments, setting up secrets management, configuring RBAC, health checks, notifications…
- **Tribal knowledge**: Most teams figure it out through blog posts, trial and error, and copying from previous jobs. The "right" structure lives in someone's head, not in code.
- **Inconsistency**: Every team in the organization ends up with a slightly different GitOps setup, making platform support harder.

**gitops-bootstrap** solves this by encoding opinionated best practices into a single CLI that scaffolds a complete, production-ready GitOps workflow in minutes.

## Core Principles

| Principle | Description |
|---|---|
| **Opinionated defaults, escape hatches everywhere** | Sensible defaults for 90% of cases, with every choice overridable via flags or config. |
| **Convention over configuration** | Standard directory structure and naming so teams across the org speak the same language. |
| **Day-2 ready** | Not just initial setup — includes patterns for promotions, rollbacks, secrets rotation, and drift detection. |
| **Cluster-agnostic** | Works on EKS, GKE, AKS, k3s, kind, or any conformant Kubernetes cluster. |
| **Idempotent** | Safe to re-run. Applies only what's missing, never overwrites existing customizations. |

## What It Sets Up

Running `gitops-bootstrap init` on a cluster produces:

### In the Cluster

- GitOps controller installed and configured (ArgoCD by default, Flux as alternative)
- Namespace structure for the controller and managed environments
- RBAC for the GitOps controller (least privilege)
- Secrets management operator (Sealed Secrets by default, External Secrets Operator as alternative)
- Optional: Ingress for the ArgoCD UI with TLS

### In the Git Repository

```
gitops-repo/
├── bootstrap/                      # One-time cluster setup
│   ├── argocd/                     # ArgoCD installation manifests
│   │   ├── namespace.yaml
│   │   ├── install.yaml            # Pinned ArgoCD version
│   │   └── argocd-cm.yaml         # Custom configuration
│   └── sealed-secrets/             # Secrets management setup
│       ├── namespace.yaml
│       └── controller.yaml
│
├── apps/                           # Application definitions (App of Apps pattern)
│   ├── _root.yaml                  # Root Application pointing to this directory
│   ├── my-api.yaml                 # ArgoCD Application for my-api
│   └── my-frontend.yaml           # ArgoCD Application for my-frontend
│
├── environments/                   # Per-environment configuration
│   ├── base/                       # Shared base manifests
│   │   ├── my-api/
│   │   │   ├── kustomization.yaml
│   │   │   ├── deployment.yaml
│   │   │   ├── service.yaml
│   │   │   └── hpa.yaml
│   │   └── my-frontend/
│   │       ├── kustomization.yaml
│   │       ├── deployment.yaml
│   │       └── service.yaml
│   │
│   ├── dev/                        # Dev overrides
│   │   ├── my-api/
│   │   │   ├── kustomization.yaml  # patches: replicas=1, resources=small
│   │   │   └── sealed-secret.yaml
│   │   └── kustomization.yaml
│   │
│   ├── staging/                    # Staging overrides
│   │   ├── my-api/
│   │   │   ├── kustomization.yaml
│   │   │   └── sealed-secret.yaml
│   │   └── kustomization.yaml
│   │
│   └── production/                 # Production overrides
│       ├── my-api/
│       │   ├── kustomization.yaml  # patches: replicas=3, resources=large, PDB
│       │   └── sealed-secret.yaml
│       └── kustomization.yaml
│
├── platform/                       # Platform-level services (managed by platform team)
│   ├── cert-manager/
│   ├── external-dns/
│   ├── ingress-nginx/
│   └── monitoring/
│
├── policies/                       # OPA/Kyverno policies (optional)
│   ├── require-labels.yaml
│   ├── disallow-latest-tag.yaml
│   └── require-resource-limits.yaml
│
└── docs/
    ├── ARCHITECTURE.md             # How this repo is structured and why
    ├── ADDING-AN-APP.md            # Step-by-step guide for dev teams
    ├── SECRETS.md                  # How to manage secrets in this setup
    └── TROUBLESHOOTING.md          # Common issues and fixes
```

## Features (Planned)

### Phase 1 — Core Bootstrap

- [ ] Interactive CLI wizard: choose GitOps controller, secrets manager, environments
- [ ] Non-interactive mode via flags/config file for CI/automation
- [ ] Install ArgoCD (Helm-based, pinned version) with opinionated defaults
- [ ] Generate repo structure following App of Apps pattern
- [ ] Kustomize-based environment management (base + overlays)
- [ ] Sealed Secrets setup with key generation and backup instructions
- [ ] Generate RBAC manifests for the GitOps controller
- [ ] Scaffold example application with full environment promotion path
- [ ] Generate documentation (ARCHITECTURE.md, ADDING-AN-APP.md, SECRETS.md)

### Phase 2 — Flux Support & Advanced Secrets

- [ ] Flux CD as alternative GitOps controller
- [ ] External Secrets Operator integration (AWS Secrets Manager, Vault)
- [ ] SOPS-based secrets as a third option
- [ ] Multi-cluster support (hub-spoke model)
- [ ] Helm chart support alongside Kustomize

### Phase 3 — Day-2 Operations

- [ ] `gitops-bootstrap add-app <name>` — scaffold a new application with all environments
- [ ] `gitops-bootstrap add-env <name>` — add a new environment to all applications
- [ ] `gitops-bootstrap validate` — lint the repo structure and check for common mistakes
- [ ] `gitops-bootstrap diff <env-a> <env-b>` — compare configuration between environments
- [ ] `gitops-bootstrap promote <app> --from dev --to staging` — generate promotion PR
- [ ] Pre-commit hooks for manifest validation (kubeval, kustomize build)

### Phase 4 — Platform Integration

- [ ] Notifications setup (Slack, Teams) for sync status
- [ ] ArgoCD Notifications integration
- [ ] Image Updater configuration for automated image promotions
- [ ] GitHub Actions / GitLab CI workflow templates for PR-based promotions
- [ ] Webhook configuration for automatic sync on push
- [ ] Dashboard: terminal-based overview of sync status across environments

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│                  gitops-bootstrap CLI                     │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────┐   │
│  │   Wizard     │  │  Scaffolder  │  │  Installer    │   │
│  │              │  │              │  │               │   │
│  │ - Interactive│  │ - Repo       │  │ - Helm        │   │
│  │ - Config file│  │   structure  │  │ - kubectl     │   │
│  │ - Flags     │  │ - Templates  │  │ - Wait/health │   │
│  │              │  │ - Docs       │  │               │   │
│  └──────┬───────┘  └──────┬───────┘  └──────┬────────┘   │
│         │                 │                  │            │
│         └─────────────────┼──────────────────┘            │
│                           │                               │
│  ┌────────────────────────┴─────────────────────────┐     │
│  │              Template Engine                      │     │
│  │  Jinja2 templates for all generated manifests     │     │
│  └───────────────────────────────────────────────────┘     │
│                           │                               │
└───────────────────────────┼───────────────────────────────┘
                            │
               ┌────────────┼───────────────┐
               │            │               │
               ▼            ▼               ▼
         ┌──────────┐ ┌──────────┐   ┌────────────┐
         │  Git     │ │  K8s API │   │   Helm     │
         │  Repo    │ │          │   │   (ArgoCD)  │
         └──────────┘ └──────────┘   └────────────┘
```

### Component Responsibilities

- **Wizard**: Gathers user preferences through interactive prompts or config file/flags. Validates input and produces a normalized configuration object.
- **Scaffolder**: Generates the Git repository structure from Jinja2 templates. Handles directory creation, manifest rendering, and documentation generation. Idempotent — detects existing files and skips them.
- **Installer**: Applies bootstrap manifests to the cluster. Installs ArgoCD/Flux via Helm, sets up secrets management, configures RBAC. Includes health checks to verify successful installation.
- **Template Engine**: Jinja2-based rendering layer used by both Scaffolder and Installer. All generated manifests are templates with sensible defaults that can be customized.

## Tech Stack

| Component | Technology | Rationale |
|---|---|---|
| Language | **Python 3.11+** | Consistent with team's tooling ecosystem, fast iteration |
| CLI framework | **Typer** | Modern, type-hinted, supports both interactive and non-interactive modes |
| Terminal UI | **Rich** | Interactive prompts, progress indicators, beautiful output |
| Template engine | **Jinja2** | Industry standard for manifest templating, familiar to most engineers |
| K8s client | **kubernetes (official)** | For health checks and cluster introspection |
| Helm integration | **subprocess (helm CLI)** | Helm's Python bindings are immature; wrapping the CLI is more reliable |
| Config file | **YAML (Pydantic)** | Natural format for K8s engineers, validated with Pydantic models |
| Package manager | **uv** | Fast dependency resolution |
| Build system | **pyproject.toml** (hatch/hatchling) | Modern Python packaging standard |
| Testing | **pytest** | With tmp_path fixtures for repo scaffolding tests |
| Linting | **Ruff** | Fast, comprehensive |

## Planned CLI Interface

```bash
# Interactive wizard (recommended for first-time setup)
gitops-bootstrap init

# Non-interactive with flags
gitops-bootstrap init \
  --controller argocd \
  --controller-version 2.13.1 \
  --secrets sealed-secrets \
  --environments dev,staging,production \
  --repo-path ./gitops-repo \
  --cluster-context prod-eu-west-1

# From config file (for reproducibility / team standardization)
gitops-bootstrap init --config bootstrap-config.yaml

# Add a new application to the existing structure
gitops-bootstrap add-app my-new-service \
  --type deployment \
  --port 8080 \
  --has-ingress \
  --has-hpa

# Add a new environment
gitops-bootstrap add-env qa --base staging

# Validate repo structure
gitops-bootstrap validate ./gitops-repo

# Compare environments
gitops-bootstrap diff dev production --app my-api

# Promote an app between environments
gitops-bootstrap promote my-api --from staging --to production
```

### Example Config File

```yaml
# bootstrap-config.yaml
controller:
  type: argocd
  version: "2.13.1"
  ingress:
    enabled: true
    host: argocd.internal.company.com
    tls: true

secrets:
  type: sealed-secrets

environments:
  - name: dev
    auto_sync: true
    prune: true
  - name: staging
    auto_sync: true
    prune: false
  - name: production
    auto_sync: false      # manual sync for production
    prune: false
    require_pr: true

applications:
  - name: example-api
    type: deployment
    port: 8080
    replicas:
      dev: 1
      staging: 2
      production: 3
    has_ingress: true
    has_hpa: true
    hpa:
      min_replicas: 2
      max_replicas: 10
      target_cpu: 70

platform_services:
  cert_manager: true
  external_dns: false
  ingress_nginx: true
  monitoring: false

policies:
  enabled: true
  engine: kyverno
```

### Interactive Wizard Flow (Planned)

```
$ gitops-bootstrap init

  ╭─────────────────────────────────────────╮
  │        gitops-bootstrap v0.1.0          │
  │   From zero to GitOps in one command    │
  ╰─────────────────────────────────────────╯

  ? Select GitOps controller:
    ❯ ArgoCD (recommended)
      Flux CD

  ? ArgoCD version: (2.13.1)

  ? Select secrets management:
    ❯ Sealed Secrets (simple, self-contained)
      External Secrets Operator (AWS SM, Vault, etc.)
      SOPS (git-native encryption)

  ? Environments to create: (dev, staging, production)

  ? Scaffold an example application? (Y/n)

  ? Target repo path: (./gitops-repo)

  ? Cluster context: (current: prod-eu-west-1)

  ⠸ Installing ArgoCD v2.13.1...          ✓
  ⠸ Setting up Sealed Secrets...          ✓
  ⠸ Generating repo structure...          ✓
  ⠸ Creating example application...       ✓
  ⠸ Generating documentation...           ✓
  ⠸ Verifying cluster health...           ✓

  ✓ GitOps bootstrap complete!

  Next steps:
    1. cd ./gitops-repo && git init && git add -A && git commit -m "feat: initial gitops structure"
    2. Push to your Git provider
    3. ArgoCD UI: https://argocd.internal.company.com
    4. Read docs/ADDING-AN-APP.md to onboard your first real application
```

## Design Decisions

### Why App of Apps?

The [App of Apps pattern](https://argo-cd.readthedocs.io/en/stable/operator-manual/cluster-bootstrapping/#app-of-apps-pattern) is the most widely adopted approach for managing multiple applications with ArgoCD. A single root Application watches the `apps/` directory, and each file there defines an Application pointing to its environment-specific manifests. This provides:
- Single entry point for the entire cluster state.
- Self-service: dev teams add an Application YAML to onboard.
- Declarative: the list of applications is version-controlled.

### Why Kustomize over Helm for app manifests?

- Kustomize works with plain YAML — no templating language to learn.
- Overlays make environment differences explicit and auditable.
- Better for GitOps: `kustomize build` output is deterministic.
- Helm is used for installing third-party software (ArgoCD, cert-manager), not for application manifests.

### Why Sealed Secrets as default?

- Zero external dependencies (no Vault, no cloud provider secrets manager).
- Works on any cluster, including local development (kind, k3s).
- Encrypted secrets can live in Git safely.
- Simple mental model: `kubeseal` encrypts, controller decrypts.
- External Secrets Operator is offered as an alternative for teams already using AWS Secrets Manager, Vault, etc.

### Why Python over Go?

- Faster development cycle for a scaffolding/templating tool.
- Jinja2 is the best templating engine available, and it's Python-native.
- The CLI doesn't need to be compiled or distributed as a single binary (pipx handles this well).
- Consistent with the rest of the tooling ecosystem (kubeshield, kube-cost-lens, platformhub).

### Why not just use a Helm chart for everything?

Helm charts are great for distributing reusable software, but they're a poor fit for representing an organization's unique GitOps repository structure. Each team's environments, naming conventions, and promotion workflows are different. A scaffolding tool that generates plain YAML with Kustomize overlays gives teams full ownership and visibility over their manifests, without the abstraction layer of `values.yaml`.

## Project Structure (Planned)

```
gitops-bootstrap/
├── src/
│   └── gitops_bootstrap/
│       ├── __init__.py
│       ├── cli.py                  # Typer CLI entrypoint
│       ├── wizard/
│       │   ├── __init__.py
│       │   ├── prompts.py          # Interactive prompts (Rich)
│       │   └── config.py           # Config file parsing (Pydantic)
│       ├── scaffolder/
│       │   ├── __init__.py
│       │   ├── repo.py             # Repo structure generation
│       │   ├── apps.py             # Application manifest generation
│       │   ├── environments.py     # Environment overlay generation
│       │   └── docs.py             # Documentation generation
│       ├── installer/
│       │   ├── __init__.py
│       │   ├── argocd.py           # ArgoCD installation via Helm
│       │   ├── flux.py             # Flux installation
│       │   ├── secrets.py          # Sealed Secrets / ESO setup
│       │   └── health.py           # Post-install health checks
│       ├── templates/              # Jinja2 templates for all manifests
│       │   ├── argocd/
│       │   ├── flux/
│       │   ├── apps/
│       │   ├── environments/
│       │   ├── platform/
│       │   ├── policies/
│       │   └── docs/
│       └── models.py               # Pydantic models for configuration
├── tests/
│   ├── conftest.py
│   ├── test_wizard.py
│   ├── test_scaffolder.py
│   ├── test_installer.py
│   └── fixtures/
│       └── sample-configs/
├── examples/
│   ├── bootstrap-config.yaml       # Example config file
│   └── github-actions-promote.yml  # Example promotion workflow
├── pyproject.toml
├── LICENSE
└── README.md
```

## Related & Prior Art

| Tool | Comparison |
|---|---|
| [ArgoCD Autopilot](https://github.com/argoproj-labs/argocd-autopilot) | Go-based, ArgoCD-only, opinionated but less flexible. Inspiration for the App of Apps approach. |
| [Flux Bootstrap](https://fluxcd.io/flux/cmd/flux_bootstrap/) | Built into Flux CLI, Flux-only, focused on controller installation. |
| [Kubefirst](https://kubefirst.io/) | Full platform (CI/CD, secrets, IDP), heavier scope, includes cloud provisioning. |
| [Backstage](https://backstage.io/) | Developer portal with scaffolding capabilities, much larger scope. |

**gitops-bootstrap** differentiates by being lightweight, controller-agnostic, and focused exclusively on the GitOps repo structure + controller setup — without trying to be an entire platform.

## Contributing

Contributions are welcome. Please open an issue to discuss your idea before submitting a PR.

This project follows:
- [Conventional Commits](https://www.conventionalcommits.org/) for commit messages.
- Trunk-based development with short-lived feature branches.
- All code must pass `ruff check`, `ruff format --check`, and `mypy` before merge.

## License

[Apache License 2.0](LICENSE)
