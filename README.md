# th - Terminal Help

Get shell commands from an LLM directly in your terminal.

## Installation

### macOS/Linux

```bash
curl -sSL https://raw.githubusercontent.com/DarioHefti/th/refs/heads/main/scripts/install.sh | bash
```

### Windows

```powershell
irm https://raw.githubusercontent.com/DarioHefti/th/refs/heads/main/scripts/install.ps1 | iex
```

### From Source

```bash
go install github.com/DarioHefti/th@latest
```

## Usage

```bash
# Get a command for listing all files modified today
th "list all files modified today"

# Find large files over 100MB and copy to clipboard
th "find large files over 100MB" --c

# Re-run setup wizard
th --config
```

## Setup

1. Run `th --config` to select a free model
2. Available models:
   - `minimax-m2.5` (MiniMax M2.5, default)
   - `big-pickle` (Stealth model)
   - `nemotron-3-super-free` (Nemotron 3 Super)

## Options

- `--c` - Copy result to clipboard
- `--config` - Run setup wizard

## Privacy

When you use th, then the following data will be sent to an llm.

- **Your prompt** - The text query you type (e.g., "list all files modified today")
- **OS type** - Sent as part of the API request for context-aware responses
- **File tree (depth 3)** - Top 3 levels of your current directory structure
- **Git info** - Current git branch and status (if inside a git repository)

## Development

### Creating a Release

```bash
git tag v1.0.0
git push --tags
```

This will trigger the release workflow which builds binaries for Linux, macOS, and Windows.
