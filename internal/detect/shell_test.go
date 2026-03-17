package detect

import (
	"os"
	"strings"
	"testing"
)

func TestEnvironmentSystemPrompt(t *testing.T) {
	origCWD := os.Getenv("PWD")
	os.Setenv("PWD", "/test/dir")
	defer os.Setenv("PWD", origCWD)

	env := &Environment{
		OS:           "linux",
		Shell:        "bash",
		ShellVersion: "5.1.0",
		CWD:          "/test/dir",
	}

	prompt := env.SystemPrompt()

	if !strings.Contains(prompt, "linux") {
		t.Error("prompt should contain OS")
	}
	if !strings.Contains(prompt, "bash") {
		t.Error("prompt should contain Shell")
	}
	if !strings.Contains(prompt, "/test/dir") {
		t.Error("prompt should contain CWD")
	}
	if strings.Contains(prompt, "```") {
		t.Error("prompt should not contain code fences")
	}
}

func TestGetCWD(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Skip("skipping: could not get cwd")
	}
	result := getCWD()
	if result != cwd {
		t.Errorf("getCWD: got %q, want %q", result, cwd)
	}
}

func TestDetect(t *testing.T) {
	env, err := Detect()
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if env.OS == "" {
		t.Error("OS should not be empty")
	}
	if env.Shell == "" {
		t.Error("Shell should not be empty")
	}
	if env.CWD == "" {
		t.Error("CWD should not be empty")
	}
}
