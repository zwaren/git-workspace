package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git-workspace/internal/templates"

	"github.com/spf13/cobra"
)

type RepoData struct {
	RepoName  string
	RepoURL   string
	CloneDate string
}

var addCmd = &cobra.Command{
	Use:   "add [repository-url]",
	Short: "Add a repository to the workspace",
	Long: `Add a Git repository to the workspace by cloning it as a submodule and wrapping it with a development structure.
The repository will be added as a submodule in the src directory of its wrapper.`,
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
		absWrapperDir, err := filepath.Abs(wrapperDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %v", err)
		}

		// Initialize Git repository in the wrapper directory
		wrapperGit := NewGitWrapper(absWrapperDir)
		if err := wrapperGit.Init(); err != nil {
			return fmt.Errorf("failed to initialize wrapper repository: %v", err)
		}

		// Add the repository as a submodule in the src directory
		srcDir := "src"
		if err := wrapperGit.AddSubmodule(repoURL, srcDir); err != nil {
			return fmt.Errorf("failed to add repository as submodule: %v", err)
		}

		// Initialize and update the submodule
		if err := wrapperGit.InitSubmodules(); err != nil {
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

		// Create initial commit with templates and submodule
		if err := wrapperGit.CreateInitialCommit("Initial wrapper setup with submodule"); err != nil {
			return fmt.Errorf("failed to create initial commit: %v", err)
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

	// Create .gitkeep in local directory
	localGitKeep := filepath.Join(wrapperDir, "local", ".gitkeep")
	if err := os.WriteFile(localGitKeep, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create .gitkeep: %v", err)
	}

	return nil
}
