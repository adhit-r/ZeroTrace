# GitHub Actions Workflows

This directory contains GitHub Actions workflows for ZeroTrace CI/CD pipeline.

## Workflows

### CI (`ci.yml`)

Continuous Integration workflow that runs on every push and pull request.

**Jobs:**
- **api-go**: Tests, builds, and lints the Go API service
- **agent-go**: Tests, builds, and lints the Go agent service
- **enrichment-python**: Tests and lints the Python enrichment service
- **web-react**: Builds and lints the React frontend
- **integration**: Runs integration tests with Docker Compose

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

### Documentation (`docs.yml`)

Deploys documentation to GitHub Pages.

**Features:**
- Builds documentation from `docs/` directory
- Deploys to GitHub Pages automatically
- Runs on push to `main` branch or manual trigger

**Access:**
- Documentation available at: `https://[username].github.io/ZeroTrace/`

### Release (`release.yml`)

Creates releases when tags are pushed.

**Features:**
- Builds binaries for all platforms (Linux, macOS, Windows)
- Creates GitHub release with binaries
- Uploads release artifacts

**Triggers:**
- Push of tags matching `v*.*.*` pattern
- Manual workflow dispatch

**Usage:**
```bash
# Create and push tag
git tag v1.0.0
git push origin v1.0.0
```

### CodeQL (`codeql.yml`)

Security analysis using GitHub CodeQL.

**Features:**
- Analyzes Go, JavaScript, and Python code
- Detects security vulnerabilities
- Runs on push, PR, and weekly schedule

**Languages:**
- Go
- JavaScript/TypeScript
- Python

### Docker (`docker.yml`)

Builds and pushes Docker images to GitHub Container Registry.

**Images:**
- `ghcr.io/[username]/zerotrace-api`
- `ghcr.io/[username]/zerotrace-enrichment`
- `ghcr.io/[username]/zerotrace-web`
- `ghcr.io/[username]/zerotrace-agent`

**Triggers:**
- Push to `main` or `develop` branches
- Push of tags matching `v*.*.*` pattern
- Pull requests (build only, no push)

## Workflow Status

View workflow status at: https://github.com/[username]/ZeroTrace/actions

## Manual Triggers

Some workflows can be triggered manually from the Actions tab:
- **Documentation**: Deploy documentation manually
- **Release**: Create release with custom version

## Secrets

No secrets required for basic workflows. For advanced features:
- `GITHUB_TOKEN`: Automatically provided
- `CODECOV_TOKEN`: Optional, for code coverage reports

## Badges

Add these badges to your README:

```markdown
![CI](https://github.com/adhit-r/ZeroTrace/workflows/CI/badge.svg)
![CodeQL](https://github.com/adhit-r/ZeroTrace/workflows/CodeQL%20Analysis/badge.svg)
![Release](https://github.com/adhit-r/ZeroTrace/workflows/Release/badge.svg)
```

## Troubleshooting

### Workflow Fails

1. Check workflow logs in Actions tab
2. Verify all dependencies are installed
3. Ensure environment variables are set correctly
4. Check for linting errors

### Docker Build Fails

1. Verify Dockerfile exists in `docker/[component]/Dockerfile`
2. Check Docker Buildx setup
3. Verify container registry permissions

### Documentation Not Deploying

1. Check GitHub Pages settings
2. Verify workflow has correct permissions
3. Check build logs for errors

## Contributing

When adding new workflows:
1. Follow existing workflow patterns
2. Add appropriate triggers
3. Include error handling
4. Document in this README
5. Test on feature branches first

---

**Last Updated**: January 2025

