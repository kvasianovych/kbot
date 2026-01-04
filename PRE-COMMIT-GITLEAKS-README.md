# Gitleaks Pre-commit Hook

A Python-based Git pre-commit hook that automatically detects secrets and sensitive information in your commits using [Gitleaks](https://github.com/gitleaks/gitleaks).

## Features

- 🔍 Automatically scans staged files for secrets before each commit
- 📦 Auto-installs Gitleaks if not already present
- 🖥️ Cross-platform support (Linux, macOS, Windows)
- ⚙️ Easy enable/disable via Git configuration
- 🏗️ Supports multiple architectures (x64, ARM64)

## Installation

### Step 1: Copy the Hook Script

Save the `pre-commit` Python script to your repository's hooks directory:

```bash
# Navigate to your repository
cd /path/to/your/repo

# Copy the script to the hooks directory
cp pre-commit.py .git/hooks/pre-commit

# Make it executable (Linux/macOS)
chmod +x .git/hooks/pre-commit
```

### Step 2: First Run

The hook will automatically install Gitleaks on the first commit if it's not already installed:

```bash
git commit -m "Test commit"
```

You'll see output like:
```
Running gitleaks pre-commit hook...
Gitleaks not found. Installing...
Detected OS: linux, Architecture: x64
Downloading gitleaks from https://github.com/gitleaks/gitleaks/releases/...
Extracting gitleaks...
Gitleaks installed successfully to /home/user/.local/bin
```

### Step 3: Add to PATH (if needed)

If you see a warning about PATH, add `~/.local/bin` to your shell profile:

```bash
# For bash (~/.bashrc or ~/.bash_profile)
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
source ~/.bashrc

# For zsh (~/.zshrc)
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.zshrc
source ~/.zshrc
```

## Usage

### Normal Operation

Once installed, the hook runs automatically on every commit:

```bash
git add file.txt
git commit -m "Add new feature"
```

**If no secrets are found:**
```
Running gitleaks pre-commit hook...
✓ No leaks detected by gitleaks
[main abc1234] Add new feature
 1 file changed, 10 insertions(+)
```

**If secrets are detected:**
```
Running gitleaks pre-commit hook...

============================================================
⚠️  GITLEAKS DETECTED POTENTIAL SECRETS
============================================================
[Detailed gitleaks output showing the detected secret]
============================================================

Commit rejected. Please remove secrets before committing.
```

### Controlling the Hook

#### Disable the Hook

To disable the hook for this repository:

```bash
git config hooks.gitleaks false
```

#### Re-enable the Hook

To re-enable it:

```bash
git config hooks.gitleaks true
```

#### Check Current Status

To see if the hook is enabled:

```bash
git config hooks.gitleaks
# Returns: true, false, or nothing (defaults to enabled)
```

#### Bypass for a Single Commit

To skip the hook for just one commit (use sparingly):

```bash
git commit --no-verify -m "Emergency fix"
```

⚠️ **Warning:** Using `--no-verify` bypasses all pre-commit checks. Only use this when absolutely necessary.

## Configuration Levels

You can set the configuration at different levels:

```bash
# Repository level (default, stored in .git/config)
git config hooks.gitleaks false

# Global level (applies to all your repos)
git config --global hooks.gitleaks false

# System level (applies to all users)
git config --system hooks.gitleaks false
```

## Troubleshooting

### Hook Doesn't Run

1. **Check if the file is executable:**
   ```bash
   ls -la .git/hooks/pre-commit
   # Should show: -rwxr-xr-x
   ```

2. **Make it executable if needed:**
   ```bash
   chmod +x .git/hooks/pre-commit
   ```

3. **Verify Python is available:**
   ```bash
   python3 --version
   ```

### Gitleaks Not Found After Installation

If you see "gitleaks not found after installation":

1. **Check if it was installed:**
   ```bash
   ls ~/.local/bin/gitleaks
   ```

2. **Add to PATH manually:**
   ```bash
   export PATH="$PATH:$HOME/.local/bin"
   ```

3. **Or create a symlink to a directory already in PATH:**
   ```bash
   sudo ln -s ~/.local/bin/gitleaks /usr/local/bin/gitleaks
   ```

### Manual Gitleaks Installation

If auto-installation fails, install manually:

**macOS (Homebrew):**
```bash
brew install gitleaks
```

**Linux:**
```bash
# Download and install manually
wget https://github.com/gitleaks/gitleaks/releases/download/v8.21.2/gitleaks_8.21.2_linux_x64.tar.gz
tar -xzf gitleaks_8.21.2_linux_x64.tar.gz
sudo mv gitleaks /usr/local/bin/
```

**Windows:**
```bash
# Using Chocolatey
choco install gitleaks

# Or using Scoop
scoop install gitleaks
```

### Permission Denied on Windows

On Windows, you may need to:
1. Run Git Bash as Administrator
2. Or adjust execution policies in PowerShell

## Team Setup

To share this hook with your team, you can:

### Option 1: Manual Installation

Document the installation steps in your README and have each team member install it manually.

### Option 2: Setup Script

Create a `setup-hooks.sh` script in your repository:

```bash
#!/bin/bash
cp scripts/pre-commit.py .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
echo "Pre-commit hook installed successfully!"
```

Team members run:
```bash
./setup-hooks.sh
```

### Option 3: Git Template Directory

Set up a template directory that applies to all new repos:

```bash
# Create template
mkdir -p ~/.git-templates/hooks
cp pre-commit.py ~/.git-templates/hooks/pre-commit
chmod +x ~/.git-templates/hooks/pre-commit

# Configure Git to use it
git config --global init.templatedir ~/.git-templates

# Apply to existing repos
git init
```

## What Gitleaks Detects

Gitleaks scans for common secrets including:

- API keys and tokens
- AWS credentials
- Private keys
- Database passwords
- OAuth tokens
- Slack tokens
- Generic passwords in code
- And many more patterns

## Customization

To customize Gitleaks behavior, create a `.gitleaks.toml` file in your repository root. See the [Gitleaks documentation](https://github.com/gitleaks/gitleaks#configuration) for configuration options.

## Version Information

- **Gitleaks Version:** 8.21.2
- **Supported Platforms:** Linux, macOS, Windows
- **Supported Architectures:** x64, ARM64

To update to a newer version of Gitleaks, edit the `GITLEAKS_VERSION` variable in the script.

## Uninstalling

To remove the hook:

```bash
rm .git/hooks/pre-commit
```

To uninstall Gitleaks:

```bash
rm ~/.local/bin/gitleaks
```

## Support

For issues with:
- **This hook:** Check the troubleshooting section above
- **Gitleaks itself:** Visit [Gitleaks GitHub Issues](https://github.com/gitleaks/gitleaks/issues)
- **Git hooks in general:** See [Git Hooks Documentation](https://git-scm.com/docs/githooks)

## License

This hook script is provided as-is for use in your projects. Gitleaks itself is licensed under the MIT License.
