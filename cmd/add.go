// Package cmd provides command-line interfaces for interacting with SSH configurations.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// ConfigOptions represents the options for SSH configuration
type ConfigOptions struct {
	HostName  string
	IPAddress string
	Username  string
	SSHKey    string
}

// AddCmd represents the Cobra command for adding a new SSH configuration or keys.
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH configuration or keys",
	Long: `Add a new SSH configuration or keys to your SSH setup.
This command can be used to:
- Add a new SSH configuration to ~/.ssh/config
- Add GitHub keys to ~/.ssh/authorized_keys
- Add GitLab keys to ~/.ssh/authorized_keys`,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Add a new SSH configuration",
	Long:  "Add a new SSH configuration to your ~/.ssh/config file",
	RunE:  runConfigCmd,
}

var configOptions ConfigOptions

var gitHubKeyCmd = &cobra.Command{
	Use:   "github [username]",
	Short: "Add GitHub keys to authorized_keys",
	Long:  "Fetch and add public SSH keys from a GitHub user to your ~/.ssh/authorized_keys file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return addServiceKey("github", args[0], afero.NewOsFs())
	},
}

var gitLabKeyCmd = &cobra.Command{
	Use:   "gitlab [username]",
	Short: "Add GitLab keys to authorized_keys",
	Long:  "Fetch and add public SSH keys from a GitLab user to your ~/.ssh/authorized_keys file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return addServiceKey("gitlab", args[0], afero.NewOsFs())
	},
}

// promptForExtraArgs is a function variable for testability
var promptForExtraArgs = func(reader *bufio.Reader) (map[string]string, error) {
	extraArgs := make(map[string]string)
	fmt.Println("Enter extra SSH arguments in format key=value, type 'done' to finish:")

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read extra arguments: %w", err)
		}
		input = strings.TrimSpace(input)

		if input == "done" {
			break
		}

		parts := strings.Split(input, "=")
		if len(parts) == 2 {
			extraArgs[parts[0]] = parts[1]
		} else {
			fmt.Println("Invalid format. Please use key=value format.")
		}
	}

	return extraArgs, nil
}

func runConfigCmd(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// Prompt for required information if not provided via flags
	if err := promptForConfigOptions(reader); err != nil {
		return fmt.Errorf("failed to get configuration options: %w", err)
	}

	// Get extra SSH arguments
	extraArgs, err := promptForExtraArgs(reader)
	if err != nil {
		return fmt.Errorf("failed to get extra arguments: %w", err)
	}

	// Format the configuration block
	configBlock := utils.FormatSSHConfig(
		configOptions.HostName,
		configOptions.IPAddress,
		configOptions.Username,
		configOptions.SSHKey,
		extraArgs,
	)

	// Ensure the config file exists with correct permissions
	if err := utils.EnsureFileExists(utils.SSHPaths.Config, 0644); err != nil {
		return fmt.Errorf("failed to ensure config file exists: %w", err)
	}

	// Append the configuration
	file, err := os.OpenFile(utils.ExpandUser(utils.SSHPaths.Config), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(configBlock); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	fmt.Println("Configuration added successfully.")
	return nil
}

func promptForConfigOptions(reader *bufio.Reader) error {
	if configOptions.HostName == "" {
		fmt.Print("Enter the SSH host name: ")
		host, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read host name: %w", err)
		}
		configOptions.HostName = strings.TrimSpace(host)
	}

	if configOptions.IPAddress == "" {
		fmt.Print("Enter the IP address: ")
		ip, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read IP address: %w", err)
		}
		configOptions.IPAddress = strings.TrimSpace(ip)
	}

	if configOptions.Username == "" {
		fmt.Print("Enter the username: ")
		user, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read username: %w", err)
		}
		configOptions.Username = strings.TrimSpace(user)
	}

	if configOptions.SSHKey == "" {
		fmt.Print("Enter the SSH key path (leave empty for default): ")
		key, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read SSH key path: %w", err)
		}
		configOptions.SSHKey = strings.TrimSpace(key)
		if configOptions.SSHKey == "" {
			configOptions.SSHKey = "~/.ssh/id_rsa.pub"
		}
	}
	configOptions.SSHKey = utils.ExpandUser(configOptions.SSHKey)

	return nil
}

func addServiceKey(service, username string, fs afero.Fs) error {
	url, found := utils.ServiceURLs[service]
	if !found {
		return fmt.Errorf("invalid service specified: %s", service)
	}

	url = fmt.Sprintf(url, username)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch keys: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch keys: HTTP %d", resp.StatusCode)
	}

	keys, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read keys: %w", err)
	}

	if len(keys) == 0 {
		return fmt.Errorf("no keys found for user %s", username)
	}

	// Ensure the authorized_keys file exists with correct permissions
	if err := utils.EnsureFileExists(utils.SSHPaths.AuthorizedKeys, 0600); err != nil {
		return fmt.Errorf("failed to ensure authorized_keys file exists: %w", err)
	}

	// Open the authorized_keys file for appending
	file, err := fs.OpenFile(utils.ExpandUser(utils.SSHPaths.AuthorizedKeys), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open authorized_keys file: %w", err)
	}
	defer file.Close()

	// Add a comment and the keys
	comment := fmt.Sprintf("\n# Keys added from %s user %s via ssh-config\n", service, username)
	if _, err := file.WriteString(comment + string(keys)); err != nil {
		return fmt.Errorf("failed to write to authorized_keys: %w", err)
	}

	fmt.Printf("Successfully added %s keys for user %s\n", service, username)
	return nil
}

func init() {
	AddCmd.AddCommand(configCmd, gitHubKeyCmd, gitLabKeyCmd)

	configCmd.Flags().StringVarP(&configOptions.HostName, "host", "H", "", "SSH host name")
	configCmd.Flags().StringVarP(&configOptions.IPAddress, "ip", "I", "", "IP address")
	configCmd.Flags().StringVarP(&configOptions.Username, "user", "U", "", "Username")
	configCmd.Flags().StringVarP(&configOptions.SSHKey, "key", "K", "", "SSH key path")
}
