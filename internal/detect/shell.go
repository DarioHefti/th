package detect

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Environment struct {
	OS           string
	Shell        string
	ShellVersion string
	CWD          string
	GitBranch    string
	GitStatus    string
}

func Detect() (*Environment, error) {
	env := &Environment{
		OS:  runtime.GOOS,
		CWD: getCWD(),
	}

	shell, err := detectShell()
	if err != nil {
		return nil, fmt.Errorf("detecting shell: %w", err)
	}
	env.Shell = shell

	shellVersion, err := detectShellVersion(shell)
	if err != nil {
		env.ShellVersion = "unknown"
	} else {
		env.ShellVersion = shellVersion
	}

	env.GitBranch, env.GitStatus = detectGit()

	return env, nil
}

func getCWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "/"
	}
	return cwd
}

func getFileTree(cwd string, maxDepth int) string {
	var sb strings.Builder
	var walk func(path string, depth int)
	walk = func(path string, depth int) {
		if depth > maxDepth {
			return
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return
		}
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), ".") && entry.Name() != ".git" {
				continue
			}
			indent := strings.Repeat("  ", depth)
			sb.WriteString(indent)
			if entry.IsDir() {
				sb.WriteString(entry.Name() + "/\n")
				walk(filepath.Join(path, entry.Name()), depth+1)
			} else {
				sb.WriteString(entry.Name() + "\n")
			}
		}
	}
	sb.WriteString(cwd + "/\n")
	walk(cwd, 0)
	return sb.String()
}

func detectShell() (string, error) {
	if runtime.GOOS == "windows" {
		return detectWindowsShell()
	}

	shellEnv := os.Getenv("SHELL")
	if shellEnv != "" {
		parts := strings.Split(shellEnv, "/")
		return parts[len(parts)-1], nil
	}

	shell := os.Getenv("SHELL_SPECIAL")
	if shell != "" {
		return shell, nil
	}

	cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", os.Getppid()), "-o", "comm=")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	return "bash", nil
}

func detectWindowsShell() (string, error) {
	comspec := os.Getenv("COMSPEC")
	if strings.Contains(strings.ToLower(comspec), "powershell") {
		return "powershell", nil
	}
	return "cmd", nil
}

func detectShellVersion(shell string) (string, error) {
	var cmd *exec.Cmd

	switch shell {
	case "bash":
		cmd = exec.Command("bash", "--version")
	case "zsh":
		cmd = exec.Command("zsh", "--version")
	case "fish":
		cmd = exec.Command("fish", "--version")
	case "powershell", "pwsh":
		cmd = exec.Command("pwsh", "-Version")
	default:
		return "", fmt.Errorf("unknown shell: %s", shell)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func detectGit() (branch, status string) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err == nil {
		branch = strings.TrimSpace(string(output))
	}

	cmd = exec.Command("git", "status", "--porcelain")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		modified := 0
		added := 0
		untracked := 0
		for _, line := range lines {
			if len(line) >= 2 {
				switch line[:2] {
				case "M ", "MM":
					modified++
				case "A ", "AM":
					added++
				case "??":
					untracked++
				}
			}
		}
		if modified > 0 || added > 0 || untracked > 0 {
			status = fmt.Sprintf("%d modified, %d added, %d untracked", modified, added, untracked)
		}
	}

	return branch, status
}

func (e *Environment) SystemPrompt() string {
	fileTree := getFileTree(e.CWD, 3)

	gitInfo := ""
	if e.GitBranch != "" {
		gitInfo = fmt.Sprintf("- Git branch: %s", e.GitBranch)
		if e.GitStatus != "" {
			gitInfo += fmt.Sprintf("\n- Git status: %s", e.GitStatus)
		}
	}

	prompt := fmt.Sprintf(`You are a command-line expert. The user is on:
- OS: %s
- Shell: %s
- Working directory: %s
%s
- File tree:
%s

Given a natural language request, respond with ONLY the shell command appropriate for their environment. 
- Do NOT include any explanation, comments, or markdown.
- Do NOT use code fences or backticks.
- The response should be ONLY the raw command.
- If platform-specific, prefer %s commands.`, e.OS, e.Shell, e.CWD, gitInfo, fileTree, e.Shell)

	return prompt
}
