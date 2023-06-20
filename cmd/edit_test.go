package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// getenvFunc is a function variable for getting environment variables.
// It is assigned os.Getenv by default, but can be changed for testing.
func TestEditCmd(t *testing.T) {
	defer func() {
		getenvFunc = os.Getenv
	}()

	// Mock execCommand
	execCommand = func(command string, args ...string) *exec.Cmd {
		cmd := exec.Command("echo", "mocked")
		return cmd
	}

	// Mock getenvFunc
	getenvFunc = func(key string) string {
		if key == "EDITOR" {
			return "mock-editor"
		}
		return ""
	}

	tests := []struct {
		args           []string
		expectedOutput string
		expectedError  string
	}{
		{
			args:           []string{"config"},
			expectedOutput: "mocked\n",
		},
		{
			args:           []string{"keys"},
			expectedOutput: "mocked\n",
		},
		{
			args:           []string{"hosts"},
			expectedOutput: "mocked\n",
		},
		{
			args:           []string{"invalid"},
			expectedOutput: "Invalid argument. Use 'config' or 'keys'.\n",
		},
	}

	for _, test := range tests {
		buf := &bytes.Buffer{}
		cobraCmd := &cobra.Command{}
		cobraCmd.SetOut(buf)
		cobraCmd.SetErr(buf)

		EditCmd.Run(cobraCmd, test.args)

		output := buf.String()

		if test.expectedError != "" {
			assert.Contains(t, output, test.expectedError)
		} else {
			assert.Equal(t, test.expectedOutput, output)
		}
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprintln(os.Stdout, "mocked")
	os.Exit(0)
}
