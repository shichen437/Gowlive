## Contributing

First off, thanks for taking the time to contribute!

## Quick Development

### 1. Environment Preparation

- [Go](https://go.dev/) (version >= 1.24)
- [mise](https://mise.jdx.dev/getting-started.html)
- [Node.js](https://nodejs.org/) (version >= 24)
- [Yarn](https://yarnpkg.com/)

### 2. Clone Project

```bash
git clone https://github.com/shichen437/Gowlive.git
cd Gowlive
```

### 3. Backend Startup

```bash
# Copy configuration file
cp manifest/config/config.yaml.example manifest/config/config.yaml

# Unify development environment (requires mise, optional)
mise trust
mise install

# Install Go dependencies
go mod tidy

# Start backend service (hot update can be achieved using GoFrame CLI tool)
# CLI installation is required for the first run
mise run cli.install

# Start!
go run main.go # For hot update: gf run main.go
```

### 4. Frontend Startup

```bash
# Enter the frontend directory
cd web

# Unify development environment (requires mise, optional)
mise trust
mise install

# Install Node.js dependencies
npm install

# Start frontend development service
npm run dev # or yarn dev
```

## PR Policy

- Commit Convention: Follow [Conventional Commits_↗](https://www.conventionalcommits.org/zh-hans/v1.0.0/).
- Commit Count: Keep a PR to no more than 2 commits; squash if it exceeds.
- Scope of Changes: Keep the PR focused on a single topic; avoid mixed, omnibus changes.
- Required Notes: List key changes, compatibility impacts, and migration steps.
- Performance & Security: Ensure no noticeable performance regressions or security risks.
