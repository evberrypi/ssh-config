# SSH-Config

A modern command-line utility for managing SSH configurations with ease. `ssh-config` simplifies the process of managing SSH configurations, eliminating the need to consult manual pages or search online for common SSH tasks.

## Features

- 🔧 Add, list, remove, and edit SSH configurations
- 🔑 Fetch and manage SSH keys from GitHub and GitLab
- 📝 Edit SSH config, authorized_keys, and known_hosts files
- 🛠️ Simple and intuitive command-line interface
- 🧪 Comprehensive test coverage
- 📊 Version information and command aliases

## Installation

### From Source

```bash
go install github.com/evberrypi/ssh-config@latest
```

### From Binary

Download the latest release from the [releases page](https://github.com/evberrypi/ssh-config/releases).

## Usage

### Command Aliases

The following command aliases are available for convenience:
- `list` → `ls`
- `remove` → `rm`
- `edit` → `e`
- `help` → `?`

### Version Information

```bash
# Show version information
ssh-config version
# or
ssh-config -v
```

### Managing SSH Configurations

```bash
# Add a new SSH configuration
ssh-config add config

# List existing configurations
ssh-config list config
# or
ssh-config ls config

# Remove a configuration
ssh-config remove [hostname]
# or
ssh-config rm [hostname]

# Edit configurations
ssh-config edit config
# or
ssh-config e config
```

### Managing SSH Keys

```bash
# Add GitHub keys to authorized_keys
ssh-config add github [username]

# Add GitLab keys to authorized_keys
ssh-config add gitlab [username]

# List GitHub keys
ssh-config list github [username]
# or
ssh-config ls github [username]

# List GitLab keys
ssh-config list gitlab [username]
# or
ssh-config ls gitlab [username]
```

### Editing SSH Files

```bash
# Edit SSH config file
ssh-config edit config
# or
ssh-config e config

# Edit authorized_keys file
ssh-config edit keys
# or
ssh-config e keys

# Edit known_hosts file
ssh-config edit hosts
# or
ssh-config e hosts
```

Note: The edit commands use your default editor (specified by the `EDITOR` environment variable) or fall back to `vim` if not set.

## Project Structure

```
ssh-config/
├── cmd/           # Command implementations
│   ├── add.go
│   ├── list.go
│   ├── remove.go
│   ├── edit.go
│   └── version.go
├── version/       # Version information
│   └── version.go
├── utils/         # Utility functions
│   └── utils.go
├── main.go        # Application entry point
├── go.mod         # Go module file
├── go.sum         # Go module checksum
├── LICENSE        # MIT License
└── README.md      # This file
```

## Development

### Prerequisites

- Go 1.24 or later
- Git

### Building from Source

```bash
git clone https://github.com/evberrypi/ssh-config.git
cd ssh-config
go build
```

### Running Tests

```bash
go test ./...
```

### Building with Version Info

To include build time and git commit in your binary, use:

```sh
go build -ldflags "-X github.com/evberrypi/ssh-config/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X github.com/evberrypi/ssh-config/version.GitCommit=$(git rev-parse --short HEAD)"
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - A Commander for modern Go CLI applications
- [Testify](https://github.com/stretchr/testify) - A toolkit with common assertions and mocks