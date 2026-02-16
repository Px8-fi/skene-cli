# Installation Testing Guide

This guide is for team members testing the new installation script.

## 🧪 Testing Instructions

### Test Scenario 1: Remote Installation (Recommended)

Test the one-liner installation from GitHub:

```bash
curl -fsSL https://raw.githubusercontent.com/Px8-fi/skene-cli/main/install.sh | bash
```

**What to check:**
- ✅ Script downloads successfully
- ✅ Platform is detected correctly (macOS ARM64/Intel, Linux, Windows)
- ✅ Binary downloads from GitHub releases (if releases exist)
- ✅ Installation completes successfully
- ✅ `skene` command works after installation

**If GitHub releases don't exist yet:**
- The script should fall back to building from source (if you have Go installed)
- Or show a helpful error message

### Test Scenario 2: Manual Script Download

```bash
# Download the script
curl -fsSL https://raw.githubusercontent.com/Px8-fi/skene-cli/main/install.sh -o install.sh

# Review the script (optional)
cat install.sh

# Run it
chmod +x install.sh
./install.sh
```

### Test Scenario 3: Custom Install Location

```bash
# Install to home directory (no sudo needed)
INSTALL_DIR=~/bin ./install.sh

# Verify it's in ~/bin
~/bin/skene --version
```

### Test Scenario 4: Local Development Build

If you've cloned the repository:

```bash
cd skene-cli

# Build locally first
make build

# Use the install script (should detect local build)
./install.sh

# Or force local build
USE_LOCAL=true ./install.sh
```

## 🐛 Reporting Issues

If you encounter any issues, please report:

1. **Your platform:**
   ```bash
   uname -s  # OS
   uname -m  # Architecture
   ```

2. **Error messages:** Copy the full error output

3. **What you tried:** Which installation method did you use?

4. **Expected vs Actual:** What should have happened vs what actually happened

## ✅ Success Criteria

The installation is successful if:

- [ ] Script runs without errors
- [ ] Binary is installed to `/usr/local/bin/skene` (or custom location)
- [ ] `skene --version` or `skene --help` works
- [ ] Binary is executable
- [ ] No permission errors

## 📝 Notes

- The script requires `curl` or `wget` to download binaries
- Installing to `/usr/local/bin` requires `sudo` (you'll be prompted)
- If you don't have `sudo` access, use `INSTALL_DIR=~/bin` instead
- The script automatically detects your platform and downloads the correct binary

## 🚀 Quick Test

Run this to test everything:

```bash
# Test installation
curl -fsSL https://raw.githubusercontent.com/Px8-fi/skene-cli/main/install.sh | bash

# Verify it works
skene --version || skene --help || echo "skene command found at: $(which skene)"
```
