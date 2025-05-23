package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git-workspace/internal/fsutil"
	"git-workspace/internal/templates"

	"github.com/spf13/cobra"
)

type RepoData struct {
	RepoName  string
	RepoURL   string
	CloneDate string
}

// copyRepoLocalDir copies the local directory from the repository to the wrapper if it exists
func copyRepoLocalDir(repoDir, wrapperDir string) error {
	srcLocalDir := filepath.Join(repoDir, "local")
	dstLocalDir := filepath.Join(wrapperDir, "local")

	if _, err := os.Stat(srcLocalDir); err == nil {
		// Remove the default local directory if it exists
		os.RemoveAll(dstLocalDir)

		// Copy the local directory from the repository
		if err := fsutil.CopyDir(srcLocalDir, dstLocalDir); err != nil {
			return fmt.Errorf("failed to copy local directory: %v", err)
		}
		fmt.Printf("Copied local/ directory from repository to wrapper\n")
	}
	return nil
}

var addCmd = &cobra.Command{
	Use:   "add [repository-url]",
	Short: "Add a repository to the workspace",
	Long: `Add a Git repository to the workspace by cloning it as a submodule and wrapping it with a development structure.
The repository will be added as a submodule in the repo directory of its wrapper.
If the repository contains a local/ directory, it will be copied to the wrapper's local/ directory.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoURL := args[0]

		// Extract repository name from URL
		repoName := extractRepoName(repoURL)
		if repoName == "" {
			return fmt.Errorf("could not extract repository name from URL: %s", repoURL)
		}

		// Ensure we're in a workspace (check for repos directory)
		if _, err := os.Stat("repos"); os.IsNotExist(err) {
			return fmt.Errorf("repos directory not found. Are you in a workspace directory?")
		}

		// Create development wrapper directory
		wrapperDir := filepath.Join("repos", repoName)
		if err := os.MkdirAll(wrapperDir, 0755); err != nil {
			return fmt.Errorf("failed to create wrapper directory: %v", err)
		}

		// Get absolute path for Git operations
		absWorkspaceDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get workspace directory: %v", err)
		}

		// Use workspace Git repository
		workspaceGit := NewGitWrapper(absWorkspaceDir)

		// Add the repository as a submodule in the wrapper's repo directory
		repoDir := filepath.Join(wrapperDir, "repo")
		if err := workspaceGit.AddSubmodule(repoURL, repoDir); err != nil {
			return fmt.Errorf("failed to add repository as submodule: %v", err)
		}

		// Initialize and update the submodule
		if err := workspaceGit.InitSubmodules(); err != nil {
			return fmt.Errorf("failed to initialize submodules: %v", err)
		}

		// Prepare template data
		data := RepoData{
			RepoName:  repoName,
			RepoURL:   repoURL,
			CloneDate: time.Now().Format("2006-01-02 15:04:05"),
		}

		// Process development wrapper templates
		if err := processRepoTemplates(wrapperDir, data); err != nil {
			return fmt.Errorf("failed to process templates: %v", err)
		}

		// Copy local directory from repository if it exists
		if err := copyRepoLocalDir(repoDir, wrapperDir); err != nil {
			return err
		}

		// Stage and commit changes
		if err := workspaceGit.CreateInitialCommit(fmt.Sprintf("Add %s repository with development wrapper", repoName)); err != nil {
			return fmt.Errorf("failed to create commit: %v", err)
		}

		fmt.Printf("Successfully added repository '%s' to workspace as a submodule\n", repoName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func extractRepoName(url string) string {
	// Remove .git extension if present
	url = strings.TrimSuffix(url, ".git")

	// Split the URL by / and take the last part
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func processRepoTemplates(wrapperDir string, data RepoData) error {
	// Create required directories
	dirs := []string{
		filepath.Join(wrapperDir, "local"),
		filepath.Join(wrapperDir, "scripts"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Get list of templates
	templateFiles, err := templates.ListTemplatesInDir("repo_wrapper")
	if err != nil {
		return fmt.Errorf("failed to list templates: %v", err)
	}

	// Process each template
	for _, templatePath := range templateFiles {
		tmpl, err := templates.GetTemplate(templatePath)
		if err != nil {
			return fmt.Errorf("failed to get template %s: %v", templatePath, err)
		}

		// Create output file
		outputPath := filepath.Join(wrapperDir, filepath.Base(templatePath))
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %v", filepath.Base(templatePath), err)
		}
		defer outputFile.Close()

		// Execute template
		if err := tmpl.Execute(outputFile, data); err != nil {
			return fmt.Errorf("failed to execute template %s: %v", filepath.Base(templatePath), err)
		}
	}

	return nil
}
