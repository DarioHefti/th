# Installation Locations

This document explains where files are installed when running the installer scripts.

## Linux / macOS (`scripts/install.sh`)

| Item | Location |
|------|----------|
| Binary | `~/.local/bin/th` |
| Custom directory | User-specified via `-d` or `--dir` option |

### Examples

```bash
# Default installation
./scripts/install.sh
# Installs to: ~/.local/bin/th

# Custom directory
./scripts/install.sh -d /usr/local/bin
# Installs to: /usr/local/bin/th
```

### PATH

After installation, ensure the installation directory is in your PATH:

```bash
# Bash
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
source ~/.bashrc

# Zsh
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.zshrc
source ~/.zshrc
```

---

## Windows (`scripts/install.ps1`)

| Item | Location |
|------|----------|
| Binary | `%LOCALAPPDATA%\th\bin\th.exe` |
| Custom directory | User-specified via `-InstallDir` parameter |
| PATH | Installation directory is added to user PATH |

### Examples

```powershell
# Default installation
.\scripts\install.ps1
# Installs to: C:\Users\<username>\AppData\Local\th\bin\th.exe

# Custom directory
.\scripts\install.ps1 -InstallDir "C:\Program Files\th"
# Installs to: C:\Program Files\th\th.exe
```

### PATH

The installer automatically adds the installation directory to your user PATH. You may need to restart your terminal for changes to take effect.

---

## Uninstallation

### Linux / macOS

```bash
rm ~/.local/bin/th
# Remove from PATH manually if added
```

### Windows

```powershell
# Remove the binary
Remove-Item "$env:LOCALAPPDATA\th\bin\th.exe" -Recurse -Force

# Remove from PATH (manual)
# Open System Properties → Environment Variables → User variables → Path
# Remove the entry: %LOCALAPPDATA%\th\bin
```
