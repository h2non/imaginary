package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestExitWithError(t *testing.T) {
	switch os.Getenv("IMAGINARY_TEST_EXIT") {
	case "plain":
		exitWithError("no arguments")
	case "args":
		exitWithError("one argument: %s", "value")
	}

	for mode, want := range map[string]string{
		"plain": "no arguments\n",
		"args":  "one argument: value\n",
	} {
		cmd := exec.Command(os.Args[0], "-test.run=^TestExitWithError$")
		cmd.Env = append(os.Environ(), "IMAGINARY_TEST_EXIT="+mode)
		var stderr strings.Builder
		cmd.Stderr = &stderr

		err := cmd.Run()
		if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 {
			t.Errorf("%s: exit %v, want status 1", mode, err)
		}
		if stderr.String() != want {
			t.Errorf("%s: printed %q, want %q", mode, stderr.String(), want)
		}
	}
}
