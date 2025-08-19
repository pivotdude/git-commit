# GitHub Actions for building git-commit

## CI Workflow (`.github/workflows/ci.yml`)

**Triggers:**
- Push to `master` branch
- Pull request to `master` branch

**What it does:**
- Sets up Go 1.24.6
- Caches Go modules for faster builds
- Builds the application: `go build -v -o git-commit ./main.go`
- Runs tests: `go test -v ./...`

## Release Workflow (`.github/workflows/release.yml`)

**Triggers:**
- Release publication (GitHub Release)

**What it does:**
- Builds binaries for the following platforms:
  - **Linux**: amd64, arm64
  - **Windows**: amd64 (with .exe extension)
  - **macOS**: amd64, arm64
- Creates checksums file
- Uploads all binaries to GitHub Release

## How to create a release

1. Make a commit with your changes
2. Go to the "Releases" section on GitHub
3. Click "Create a new release"
4. Enter a version tag (e.g., v1.0.0)
5. Click "Publish release"

After publishing the release, the build will automatically start and binaries will be added to the release.

## Final files in the release:

- `git-commit-linux-amd64` - Linux 64-bit Intel
- `git-commit-linux-arm64` - Linux ARM64
- `git-commit-windows-amd64.exe` - Windows 64-bit
- `git-commit-macos-amd64` - macOS Intel 64-bit
- `git-commit-macos-arm64` - macOS Apple Silicon
- `checksums.txt` - SHA256 checksums

Dockerfile is NOT required for these builds, as native Go compilation is used.