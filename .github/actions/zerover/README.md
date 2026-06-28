# ZeroVer Action

A GitHub Action for managing versioning with a [ZeroVer](https://0ver.org/) philosophy, where the major version always stays at 0.

## Features

- **ZeroVer Philosophy**: Major version never increments (stays at 0.x.x)
- **Two Versioning Strategies**:
  - **Rollover**: Increment patch version with automatic rollover to minor version
  - **Conventional Commits**: Analyze commit messages to determine version bumps
- **Dev Versions**: Automatic dev version generation with timestamps and commit hashes
- **Flexible Tag Creation**: Control when tags are created vs. when versions are calculated

## Prerequisites

**Important**: This action requires full git history to properly detect tags and calculate versions. Always use `fetch-depth: 0` with `actions/checkout`:

```yaml
- uses: actions/checkout@v7
  with:
    fetch-depth: 0 # to calculate version correctly
    persist-credentials: true # for pushing tags
```

The action will emit a warning if a shallow clone is detected.

## Usage

### Basic Usage - Get Current Version

```yaml
- uses: actions/checkout@v7
  with:
    fetch-depth: 0 # to calculate version correctly
    persist-credentials: true # for pushing tags

- uses: ./.github/actions/zerover
  id: version

- run: echo "Version is ${{ steps.version.outputs.version }}"
```

**Outputs**:
- On a tagged commit: `0.1.4`
- On an untagged commit: `0.1.4-dev.20260628143052.abc123def456`

### Create a New Tag

```yaml
- uses: actions/checkout@v7
  with:
    fetch-depth: 0 # to calculate version correctly
    persist-credentials: true # for pushing tags

- uses: ./.github/actions/zerover
  id: version
  with:
    create: true

- run: echo "Created version ${{ steps.version.outputs.version }}"
```

### Using Rollover Strategy

```yaml
- uses: ./.github/actions/zerover
  with:
    create: true
    strategy: rollover
    rollover-max: 100
```

**Behavior**:
- `0.5.50` → `0.5.51` (normal increment)
- `0.5.99` → `0.6.0` (rollover to minor)

**Disable rollover** (patch increments indefinitely):
```yaml
- uses: ./.github/actions/zerover
  with:
    create: true
    strategy: rollover
    rollover-max: 0
```

**Behavior**:
- `0.5.999` → `0.5.1000` (no rollover)

### Using Conventional Commits Strategy

```yaml
- uses: ./.github/actions/zerover
  with:
    create: true
    strategy: conventional-commits
```

**Behavior**:
- `feat: new feature` → increments minor version (`0.5.10` → `0.6.0`)
- `fix: bug fix` → increments patch version (`0.5.10` → `0.5.11`)
- `feat!: breaking change` → increments minor version (`0.5.10` → `0.6.0`)
- Commits with `BREAKING CHANGE:` footer → increments minor version

### Custom Git Identity

```yaml
- uses: ./.github/actions/zerover
  with:
    create: true
    user-name: "My Bot"
    user-email: "bot@example.com"
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `create` | Whether to create a new tag | No | `false` |
| `user-name` | Git user name for creating tags | No | `github-actions[bot]` |
| `user-email` | Git user email for creating tags | No | `41898282+github-actions[bot]@users.noreply.github.com` |
| `strategy` | Versioning strategy (`rollover` or `conventional-commits`) | No | `rollover` |
| `rollover-max` | Maximum patch version before rollover (rollover strategy only). Set to 0 to disable rollover. | No | `100` |

## Outputs

| Output | Description |
|--------|-------------|
| `version` | The current version (without 'v' prefix) |
| `created` | Whether a new tag was created (`true` or `false`) |

## Version Format

### Clean Version
When on a tagged commit or after creating a tag:
```
0.1.4
```

### Dev Version
When on an untagged commit:
```
0.1.4-dev.20260628143052.abc123def456
```

Format: `{base-version}-dev.{timestamp}.{commit-hash}`
- **timestamp**: `YYYYMMDDHHmmss` in UTC
- **commit-hash**: 12-character abbreviated hash

## Examples

### Release Workflow

```yaml
name: Release
on:
  push:
    branches: [main]

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
        with:
          fetch-depth: 0 # to calculate version correctly
          persist-credentials: true # for pushing tags
      
      - uses: ./.github/actions/zerover
        id: version
        with:
          create: true
          strategy: conventional-commits
      
      - name: Create GitHub Release
        if: steps.version.outputs.created == 'true'
        uses: actions/create-release@v1
        with:
          tag_name: v${{ steps.version.outputs.version }}
          release_name: Release ${{ steps.version.outputs.version }}
```

### Build Workflow with Dev Versions

```yaml
name: Build
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
        with:
          fetch-depth: 0 # to calculate version correctly
          persist-credentials: true # for pushing tags
      
      - uses: ./.github/actions/zerover
        id: version
      
      - name: Build with version
        run: |
          echo "Building version ${{ steps.version.outputs.version }}"
          # Use version in build artifacts, docker tags, etc.
```

## License

See repository license.
