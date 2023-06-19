package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/evberrypi/ssh-config/utils"
	"github.com/spf13/cobra"
)

type ConfigOptions struct {
	HostName  string
	IPAddress string
	Username  string
	SSHKey    string
}

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH configuration or keys",
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Add a new SSH configuration",
	Run:   runConfigCmd,
}

var configOptions ConfigOptions

var gitHubKeyCmd = &cobra.Command{
	Use:   "github-key [username]",
	Short: "Add GitHub keys to authorized_keys",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addServiceKey("github", args[0])
	},
}

var gitLabKeyCmd = &cobra.Command{
	Use:   "gitlab-key [username]",
	Short: "Add GitLab keys to authorized_keys",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addServiceKey("gitlab", args[0])
	},
}

var serviceURLs = map[string]string{
	"github": "https://github.com/%s.keys",
	"gitlab": "https://gitlab.com/%s.keys",
}

func runConfigCmd(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	if configOptions.HostName == "" {
		fmt.Print("Enter the SSH host name: ")
		configOptions.HostName, _ = reader.ReadString('\n')
		configOptions.HostName = strings.TrimSpace(configOptions.HostName)
	}

	if configOptions.IPAddress == "" {
		fmt.Print("Enter the IP address: ")
		configOptions.IPAddress, _ = reader.ReadString('\n')
		configOptions.IPAddress = strings.TrimSpace(configOptions.IPAddress)
	}

	if configOptions.Username == "" {
		fmt.Print("Enter the username: ")
		configOptions.Username, _ = reader.ReadString('\n')
		configOptions.Username = strings.TrimSpace(configOptions.Username)
	}

	if configOptions.SSHKey == "" {
		fmt.Print("Enter the SSH key path (leave empty for default): ")
		configOptions.SSHKey, _ = reader.ReadString('\n')
		configOptions.SSHKey = strings.TrimSpace(configOptions.SSHKey)
		if configOptions.SSHKey == "" {
			configOptions.SSHKey = "~/.ssh/id_rsa.pub"
		}
	}
	configOptions.SSHKey = utils.ExpandUser(configOptions.SSHKey)

	fmt.Print("Enter the SSH host name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter the IP address: ")
	ip, _ := reader.ReadString('\n')
	ip = strings.TrimSpace(ip)

	fmt.Print("Enter the username: ")
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)

	fmt.Print("Enter the SSH key path (leave empty for default): ")
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)
	if key == "" {
		key = "~/.ssh/id_rsa.pub"
	}
	key = utils.ExpandUser(key)

	extraArgs := make(map[string]string)
	fmt.Println("Enter extra SSH arguments in format key=value, type 'done' to finish:")
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "done" {
			break
		}

		parts := strings.Split(input, "=")
		if len(parts) == 2 {
			extraArgs[parts[0]] = parts[1]
		}
	}

	configBlock := fmt.Sprintf("Host %s\n    HostName %s\n    User %s\n    IdentityFile %s\n", name, ip, user, key)
	for arg, value := range extraArgs {
		configBlock += fmt.Sprintf("    %s %s\n", arg, value)
	}

	configPath := utils.ExpandUser(utils.SshConfigPath)
	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(configBlock)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func addServiceKey(service, username string) {
	url, found := serviceURLs[service]
	if !found {
		fmt.Println("Invalid service specified.")
		return
	}

	url = fmt.Sprintf(url, username)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching keys: %v\n", err)
		return
	}
	defer resp.Body.Close()

	keys, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading keys: %v\n", err)
		return
	}

	if len(keys) == 0 {
		fmt.Println("No keys found.")
		return
	}

	authorizedKeysPath := utils.ExpandUser("~/.ssh/authorized_keys")

	// First read the entire authorized_keys file to check for existing keys
	file, err := os.Open(authorizedKeysPath)
	if err != nil {
		fmt.Printf("Error opening authorized_keys for reading: %v\n", err)
		return
	}

	existingKeys, err := io.ReadAll(file)
	file.Close()
	if err != nil {
		fmt.Printf("Error reading authorized_keys: %v\n", err)
		return
	}

	// Check if key already exists
	if strings.Contains(string(existingKeys), string(keys)) {
		fmt.Println("Key already exists in authorized_keys.")
		return
	}

	// Append the key if it doesn't exist
	file, err = os.OpenFile(authorizedKeysPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Error opening authorized_keys for appending: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("\n" + string(keys))
	if err != nil {
		fmt.Printf("Error writing to authorized_keys: %v\n", err)
		return
	}

	fmt.Println("Keys added to authorized_keys successfully.")
}

func init() {
	AddCmd.AddCommand(configCmd, gitHubKeyCmd, gitLabKeyCmd)

	configCmd.Flags().StringVarP(&configOptions.HostName, "host", "H", "", "SSH host name")
	configCmd.Flags().StringVarP(&configOptions.IPAddress, "ip", "I", "", "IP address")
	configCmd.Flags().StringVarP(&configOptions.Username, "user", "U", "", "Username")
	configCmd.Flags().StringVarP(&configOptions.SSHKey, "key", "K", "", "SSH key path")
}
