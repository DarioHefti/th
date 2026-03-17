# th - Terminal Help

Get shell commands from an LLM directly in your terminal.

## Installation

### macOS/Linux

```bash
curl -sSL https://raw.githubusercontent.com/terminal-help/th/main/scripts/install.sh | bash
```

### Windows

```powershell
irm https://raw.githubusercontent.com/terminal-help/th/main/scripts/install.ps1 | iex
```

### From Source

```bash
go install github.com/terminal-help/th@latest
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

1. Create an app registration in Azure AD
2. Grant the app access to Azure AI Foundry
3. Run `th --config` to configure your credentials

## Options

- `--no-clipboard` - Don't copy result to clipboard
- `--config` - Run setup wizard
