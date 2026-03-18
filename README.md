# th - Terminal Help

Get shell commands from an LLM directly in your terminal.

## Installation

### macOS/Linux

```bash
curl -sSL https://raw.githubusercontent.com/DarioHefti/th/main/scripts/install.sh | bash
```

### Windows

```powershell
irm https://raw.githubusercontent.com/DarioHefti/th/main/scripts/install.ps1 | iex
```

### From Source

```bash
go install github.com/DarioHefti/th@latest
```

## Usage

```bash
# Get a command for listing all files modified today
th "list all files modified today"

# Find large files over 100MB
th "find large files over 100MB"

# Re-run setup wizard
th --config
```

## Setup

1. Run `th --config` to select a free model
2. Available models:
   - `minimax-m2.5-free` (MiniMax M2.5)
   - `big-pickle` (Stealth model)
   - `mimo-v2-flash-free` (MiMo V2 Flash)
   - `nemotron-3-super-free` (Nemotron 3 Super)

## Options

- `--c` - Copy result to clipboard
- `--config` - Run setup wizard

## Development

### Creating a Release

```bash
git tag v1.0.0
git push --tags
```

This will trigger the release workflow which builds binaries for Linux, macOS, and Windows.
