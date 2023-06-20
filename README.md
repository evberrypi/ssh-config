# SSH-Config

`ssh-config` is a command-line utility to manage SSH configurations to prevent users from needing to consult Man pages, or the internet. It allows you to easily add, list, remove, and edit configurations in `~/.ssh/config`. Additionally, `ssh-config` allows you to fetch public SSH keys from GitHub and GitLab user accounts.

## Structure

The project is structured as follows:
```
ssh-config-tool/
│ README.md
│ main.go
│
└───cmd/
│ add.go
│ list.go
│ remove.go
│ edit.go
│ add_test.go
│ list_test.go
│ remove_test.go
│ edit_test.go
│
└───utils/
│ utils.go
| utils_test.go
```

## Installation

To install `ssh-config`, use the following command:

```bash
go get github.com/evberrypi/ssh-config
```

## Commands

### Add

To add a new SSH configuration, use the add command.

```bash
ssh-config-tool add config
```

This command will prompt you for the SSH host name, IP address, username, SSH key path, and any additional SSH arguments. It will then add the configuration to the ~/.ssh/config file.

To add GitHub keys to authorized_keys, use the following command:

```bash
ssh-config-tool add github-key [username]
```

This command will fetch the user's keys from GitHub and add them to the ~/.ssh/authorized_keys file.

To add GitLab keys to authorized_keys, use the following command:

```bash
ssh-config-tool add gitlab-key [username]
```

This command will fetch the user's keys from GitLab and add them to the ~/.ssh/authorized_keys file.

Remember to replace `[username]` with the actual GitHub or GitLab username.

You can also refer to the updated "Commands" section in this README for the other commands.

### List

To list existing SSH configurations, use the list command.

```bash
ssh-config list
```

You can also fetch public SSH keys from GitHub or GitLab user accounts.

```bash
ssh-config list github-keys [username]
ssh-config list gitlab-keys [username]
```

### Remove

To remove an SSH configuration, use the remove command followed by the host name.

```bash
ssh-config remove [name]
```

### Edit

Use the `edit` command to edit SSH configurations.

```bash
ssh-config edit
```
By default, this command will edit the SSH configuration file. This is equivalent to running:

```bash
ssh-config edit config
```

To edit the `~/.ssh/authorized_keys` file by using the following command:

```bash
ssh-config edit keys
```

You can also edit the `~/.ssh/known_hosts` file by using the following command:

```bash
ssh-config edit hosts
```

These commands use the editor specified by the `EDITOR` environment variable. If unset, the editor will default to `vim`.


## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a pull request.