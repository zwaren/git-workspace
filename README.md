# Git Workspace CLI Tool

A command-line interface tool for managing Git workspaces.

## Installation

### Prerequisites

- Go 1.16 or higher
- Git

### Building from Source

1. Clone the repository:
```bash
git clone <your-repository-url>
cd git-workspace
```

2. Build the binary:
```bash
go build
```

3. (Optional) Add to your PATH for global access:
```bash
# On Unix-like systems
sudo mv git-workspace /usr/local/bin/

# On Windows
# Move the executable to a location in your PATH
```

## Usage

### Basic Commands

```bash
# Show help information
git-workspace --help

# Check version
git-workspace version
```

### Command Structure

```bash
git-workspace [command] [subcommand] [flags]
```

## Project Structure

```
.
├── cmd/
│   ├── root.go    # Root command definition
│   └── version.go # Version command implementation
├── main.go        # Entry point
├── go.mod         # Go module definition
└── README.md      # This file
```

## Development

### Adding New Commands

To add a new command:

1. Create a new file in the `cmd` directory
2. Define your command using Cobra
3. Add it to the root command in the `init()` function

Example:
```go
var newCmd = &cobra.Command{
    Use:   "new",
    Short: "Short description",
    Run: func(cmd *cobra.Command, args []string) {
        // Command implementation
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
```

## License

[Your chosen license]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 