package cmd

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/spf13/cobra"
	"github.com/evberrypi/ssh-config/utils"
	"os"
	"os/exec"
)

// Mock for exec.Command
type MockCommandExecutor struct {
	mock.Mock
}

func (m *MockCommandExecutor) Command(name string, arg ...string) *exec.Cmd {
	m.Called(name, arg)
	// We return nil since we're just verifying the call
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
		mockExec.On("Command", editor, []string{test.expectedArg}).Return(nil)
		
		EditCmd.Run(&cobra.Command{}, test.args)
		
		// Assert that the mock exec.Command was called with correct arguments
		mockExec.AssertExpectations(t)
	}
}
