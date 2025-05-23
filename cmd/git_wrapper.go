package cmd

import (
	"os"
	"os/exec"
)

// GitWrapper wraps Git command-line operations for the workspace
type GitWrapper struct {
	workDir string
}

// NewGitWrapper creates a new GitWrapper instance
func NewGitWrapper(workDir string) *GitWrapper {
	return &GitWrapper{
		workDir: workDir,
	}
}

// Init initializes a new Git repository in the current directory
func (g *GitWrapper) Init() error {
	cmd := exec.Command("git", "init")
	cmd.Dir = g.workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CreateInitialCommit stages all files and creates the initial commit
func (g *GitWrapper) CreateInitialCommit(message string) error {
	// Stage all files
	stageCmd := exec.Command("git", "add", ".")
	stageCmd.Dir = g.workDir
	stageCmd.Stdout = os.Stdout
	stageCmd.Stderr = os.Stderr
	if err := stageCmd.Run(); err != nil {
		return err
	}

	// Create commit
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = g.workDir
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	return commitCmd.Run()
}

// CloneRepository clones a Git repository into the specified directory
func (g *GitWrapper) CloneRepository(url, destDir string) error {
	cmd := exec.Command("git", "clone", url, destDir)
	cmd.Dir = g.workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// AddSubmodule adds a Git submodule at the specified path
func (g *GitWrapper) AddSubmodule(url, path string) error {
	cmd := exec.Command("git", "submodule", "add", url, path)
	cmd.Dir = g.workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// InitSubmodules initializes and updates all submodules
func (g *GitWrapper) InitSubmodules() error {
	// Initialize submodules
	initCmd := exec.Command("git", "submodule", "init")
	initCmd.Dir = g.workDir
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr
	if err := initCmd.Run(); err != nil {
		return err
	}

	// Update submodules
	updateCmd := exec.Command("git", "submodule", "update")
	updateCmd.Dir = g.workDir
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	return updateCmd.Run()
}
