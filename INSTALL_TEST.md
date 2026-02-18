# Installation Testing Guide

Testing checklist for the Skene CLI install script.

## Test Scenarios

### 1. One-liner Install (Remote)

```bash
curl -fsSL https://raw.githubusercontent.com/SkeneTechnologies/skene-cli/main/install.sh | bash
```

**Check:**
- Script downloads and runs without errors
- Platform detected correctly (macOS ARM64/Intel, Linux, Windows)
- Binary installed to `/usr/local/bin/skene`
- `skene` command works

If no GitHub release exists yet the script will try to build from source (requires Go) or print a helpful error.

### 2. Clone and Install

```bash
git clone https://github.com/SkeneTechnologies/skene-cli
cd skene-cli
./install.sh
```

**Check:**
- Detects you are inside the repository
- Uses existing `build/skene` if available, otherwise builds from source
- Installation succeeds

### 3. Custom Install Location

```bash
INSTALL_DIR=~/bin ./install.sh
~/bin/skene --version
```

### 4. Local Development Build

```bash
cd skene-cli
make build
./install.sh          # uses local build automatically
# or force local build:
USE_LOCAL=true ./install.sh
```

### 5. Specific Version

```bash
VERSION=v0.2.0 ./install.sh
```

## Post-Install Verification

```bash
# Any of these should confirm the binary is working:
skene --version
skene --help
which skene
```

## Reporting Issues

Include the following when filing a bug:

```bash
uname -s   # OS
uname -m   # Architecture
```

Plus the full error output, which install method you used, and what you expected vs what happened.

## Success Criteria

- [ ] Script runs without errors
- [ ] Binary installed and executable
- [ ] `skene` command is available in PATH
- [ ] No permission errors (or clean sudo prompt)

## Notes

- Requires `curl` or `wget`
- `/usr/local/bin` needs `sudo`; use `INSTALL_DIR=~/bin` to avoid it
- The script auto-detects your OS and architecture
