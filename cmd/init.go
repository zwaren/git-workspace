package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git-workspace/internal/templates"

	"github.com/spf13/cobra"
)

type WorkspaceData struct {
	Name string
}

var initCmd = &cobra.Command{
	Use:   "init [workspace-name]",
	Short: "Initialize a new workspace",
	Long: `Initialize a new workspace with the standard directory structure and template files.
This will create a new directory with the given name and set it up as a Git repository.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceName := args[0]

		// Create workspace directory
		workspaceDir, err := filepath.Abs(workspaceName)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %v", err)
		}

		if err := os.MkdirAll(workspaceName, 0755); err != nil {
			return fmt.Errorf("failed to create workspace directory: %v", err)
		}

		// Change to workspace directory
		if err := os.Chdir(workspaceName); err != nil {
			return fmt.Errorf("failed to change to workspace directory: %v", err)
		}

		// Create required directories
		dirs := []string{"repos", "scripts"}
		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create %s directory: %v", dir, err)
			}
		}

		// Create .gitkeep in repos directory
		reposGitKeep := filepath.Join("repos", ".gitkeep")
		if err := os.WriteFile(reposGitKeep, []byte{}, 0644); err != nil {
			return fmt.Errorf("failed to create .gitkeep: %v", err)
		}

		// Initialize Git repository
		git := NewGitWrapper(workspaceDir)
		if err := git.Init(); err != nil {
			return fmt.Errorf("failed to initialize git repository: %v", err)
		}

		// Process template files
		data := WorkspaceData{
			Name: strings.Title(workspaceName),
		}
		if err := processWorkspaceTemplates(data); err != nil {
			return fmt.Errorf("failed to process templates: %v", err)
		}

		// Initial git commit
		if err := git.CreateInitialCommit("Initial workspace setup"); err != nil {
			return fmt.Errorf("failed to create initial commit: %v", err)
		}

		fmt.Printf("Successfully initialized workspace '%s'\n", workspaceName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func processWorkspaceTemplates(data WorkspaceData) error {
	// Get list of templates
	templateFiles, err := templates.ListTemplatesInDir("workspace")
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
		outputPath := filepath.Base(templatePath)
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %v", outputPath, err)
		}
		defer outputFile.Close()

		// Execute template
		if err := tmpl.Execute(outputFile, data); err != nil {
			return fmt.Errorf("failed to execute template %s: %v", outputPath, err)
		}
	}

	return nil
}
