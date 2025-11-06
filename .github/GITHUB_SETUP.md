# GitHub Actions Setup

Complete CI/CD pipeline for ZeroTrace using GitHub Actions.

## Overview

ZeroTrace now has a comprehensive CI/CD pipeline that automates:
- Testing and building all components
- Security analysis with CodeQL
- Documentation deployment to GitHub Pages
- Docker image building and publishing
- Release creation and binary distribution

## Workflows Created

### 1. CI Workflow (`ci.yml`)

**Purpose**: Continuous Integration for all components

**Runs On**:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**Jobs**:
- **api-go**: Tests, builds, and lints Go API service
- **agent-go**: Tests, builds, and lints Go agent service
- **enrichment-python**: Tests and lints Python enrichment service
- **web-react**: Builds and lints React frontend
- **integration**: Runs Docker Compose integration tests

**Features**:
- Code coverage reporting (Codecov)
- Linting with golangci-lint, ruff, ESLint
- Cross-platform builds

### 2. Documentation Workflow (`docs.yml`)

**Purpose**: Deploy documentation to GitHub Pages

**Runs On**:
- Push to `main` branch (when docs change)
- Manual workflow dispatch

**Features**:
- Builds documentation from `docs/` directory
- Deploys to GitHub Pages automatically
- Accessible at: `https://[username].github.io/ZeroTrace/`

**Setup Required**:
1. Go to repository Settings → Pages
2. Enable GitHub Pages
3. Select source: "GitHub Actions"

### 3. Release Workflow (`release.yml`)

**Purpose**: Create releases with binaries

**Runs On**:
- Push of tags matching `v*.*.*` pattern
- Manual workflow dispatch

**Features**:
- Builds binaries for all platforms (Linux, macOS, Windows)
- Creates GitHub release with binaries
- Uploads release artifacts

**Usage**:
```bash
# Create and push tag
git tag v1.0.0
git push origin v1.0.0
```

### 4. CodeQL Workflow (`codeql.yml`)

**Purpose**: Security vulnerability analysis

**Runs On**:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches
- Weekly schedule (Sunday at midnight)

**Features**:
- Analyzes Go, JavaScript, and Python code
- Detects security vulnerabilities
- Automatic scanning

### 5. Docker Workflow (`docker.yml`)

**Purpose**: Build and push Docker images

**Runs On**:
- Push to `main` or `develop` branches
- Push of tags matching `v*.*.*` pattern
- Pull requests (build only, no push)

**Images**:
- `ghcr.io/[username]/zerotrace-api`
- `ghcr.io/[username]/zerotrace-enrichment`
- `ghcr.io/[username]/zerotrace-web`
- `ghcr.io/[username]/zerotrace-agent`

**Access**:
- Images available at: `https://github.com/[username]?tab=packages`

## Status Badges

Badges are automatically added to README.md:

```markdown
[![CI](https://github.com/adhit-r/ZeroTrace/actions/workflows/ci.yml/badge.svg)](https://github.com/adhit-r/ZeroTrace/actions/workflows/ci.yml)
[![CodeQL](https://github.com/adhit-r/ZeroTrace/actions/workflows/codeql.yml/badge.svg)](https://github.com/adhit-r/ZeroTrace/actions/workflows/codeql.yml)
[![Documentation](https://github.com/adhit-r/ZeroTrace/actions/workflows/docs.yml/badge.svg)](https://github.com/adhit-r/ZeroTrace/actions/workflows/docs.yml)
```

## Setup Instructions

### 1. Enable GitHub Pages

1. Go to repository Settings → Pages
2. Source: "GitHub Actions"
3. Save

### 2. Enable GitHub Packages

1. Go to repository Settings → Actions → General
2. Enable "Read and write permissions" for GitHub Packages
3. Save

### 3. Verify Workflows

1. Push to repository
2. Go to Actions tab
3. Verify workflows run successfully

### 4. (Optional) Code Coverage

1. Sign up at [codecov.io](https://codecov.io)
2. Add repository
3. Get token
4. Add as repository secret: `CODECOV_TOKEN`

## Workflow Status

View all workflows at:
- https://github.com/[username]/ZeroTrace/actions

## Manual Triggers

Some workflows can be triggered manually:

1. Go to Actions tab
2. Select workflow
3. Click "Run workflow"
4. Select branch and options
5. Click "Run workflow"

## Troubleshooting

### Workflow Fails

1. Check workflow logs in Actions tab
2. Verify all dependencies are installed
3. Check for linting errors
4. Ensure environment variables are set

### Docker Build Fails

1. Verify Dockerfile exists: `docker/[component]/Dockerfile`
2. Check Docker Buildx setup
3. Verify container registry permissions

### Documentation Not Deploying

1. Check GitHub Pages settings
2. Verify workflow has correct permissions
3. Check build logs for errors
4. Ensure `docs/` directory exists

### Release Not Creating

1. Verify tag format: `v*.*.*`
2. Check release workflow logs
3. Ensure binaries build successfully
4. Verify GitHub token permissions

## Next Steps

1. **Push to GitHub**: Commit and push all workflows
2. **Enable GitHub Pages**: Set up documentation hosting
3. **Test Workflows**: Push changes and verify workflows run
4. **Create Release**: Tag and push to test release workflow
5. **Monitor**: Watch workflows in Actions tab

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitHub Pages Documentation](https://docs.github.com/en/pages)
- [GitHub Packages Documentation](https://docs.github.com/en/packages)
- [CodeQL Documentation](https://codeql.github.com/docs/)

---

**Last Updated**: January 2025

