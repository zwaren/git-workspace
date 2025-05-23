# Git Workspace CLI Tool

A command-line interface tool for managing Git workspaces. Helps organize and manage multiple Git repositories in a structured workspace environment.

## Features

- Create structured workspaces with standardized layouts
- Add repositories as Git submodules with development wrappers
- Preserve repository-specific local development files
- Template-based workspace and repository configuration

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

# Initialize a new workspace
git-workspace init my-workspace

# Add a repository to the workspace
cd my-workspace
git-workspace add https://github.com/user/repo.git
```

### Command Details

#### Initialize Workspace (`init`)

Creates a new workspace with a standardized directory structure:

```
workspace/
├── README.md          # Workspace documentation
├── .gitignore        # Git ignore patterns
├── scripts/          # Workspace-specific scripts
└── repos/            # Directory for Git repositories
```

Usage:
```bash
git-workspace init <workspace-name>
```

#### Add Repository (`add`)

Adds a Git repository to the workspace with a development wrapper structure:

```
repos/
└── repo-name/
    ├── README.md     # Wrapper documentation
    ├── repo/         # The actual repository (as a Git submodule)
    ├── local/        # Local development files and overrides
    └── scripts/      # Repository-specific scripts
```

Features:
- Adds the repository as a Git submodule
- Creates a structured development wrapper
- Preserves repository's local/ directory if it exists
- Generates helpful documentation and scripts

Usage:
```bash
git-workspace add <repository-url>
```

## Project Structure

```
.
├── cmd/              # Command implementations
│   ├── root.go      # Root command definition
│   ├── version.go   # Version command
│   ├── init.go      # Init command
│   └── add.go       # Add command
├── internal/         # Internal packages
│   ├── templates/   # Template management
│   └── fsutil/      # File system utilities
├── main.go          # Entry point
└── README.md        # This file
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

Apache License 2.0

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 