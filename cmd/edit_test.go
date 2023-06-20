package cmd

import (
	"os"
	"os/exec"
	"testing"

	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
)

type MockCommandExecutor struct {
	mock.Mock
}

func (m *MockCommandExecutor) Command(name string, arg ...string) *exec.Cmd {
	cmd := m.Called(name, arg).Get(0)
	if cmd != nil {
		return cmd.(*exec.Cmd)
	}
	return nil
}

func TestEditCmd(t *testing.T) {
	mockExec := new(MockCommandExecutor)

	// Replace exec.Command with our mock version
	execCommand = mockExec.Command

	editor := "vim"

	// Set environment variable EDITOR to vim for testing
	os.Setenv("EDITOR", editor)

	configPath := utils.ExpandUser(utils.SshConfigPath)
	keysPath := utils.ExpandUser("~/.ssh/authorized_keys")
	hostsPath := utils.ExpandUser("~/.ssh/known_hosts")

	tests := []struct {
		args        []string
		expectedArg string
	}{
		{[]string{}, configPath},
		{[]string{"keys"}, keysPath},
		{[]string{"hosts"}, hostsPath},
	}

	for _, test := range tests {
		// Mock exec.Command call to check if it is called with correct arguments
		mockCmd := &exec.Cmd{Path: "true"} // A no-op command
		mockExec.On("Command", editor, []string{test.expectedArg}).Return(mockCmd)

		EditCmd.Run(&cobra.Command{}, test.args)

		// Assert that the mock exec.Command was called with correct arguments
		mockExec.AssertExpectations(t)
	}
}
