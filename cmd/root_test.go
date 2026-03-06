package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestExecuteExitsWhenNoTargetProvided(t *testing.T) {
	if os.Getenv("TEST_EXECUTE_NO_TARGET") == "1" {
		rootCmd.SetArgs(nil)
		Execute()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExecuteExitsWhenNoTargetProvided")
	cmd.Env = append(os.Environ(), "TEST_EXECUTE_NO_TARGET=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected subprocess to exit with error")
	}
	if !strings.Contains(string(output), "请指定目标 URL") {
		t.Fatalf("output = %q, want missing target message", output)
	}
}
