package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"text/template"
)

//go:embed repo_wrapper/* workspace/*
var templateFS embed.FS

// GetTemplate returns a parsed template for the given template path
func GetTemplate(templatePath string) (*template.Template, error) {
	content, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %v", templatePath, err)
	}

	tmpl, err := template.New(path.Base(templatePath)).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %v", templatePath, err)
	}

	return tmpl, nil
}

// ListTemplatesInDir returns a list of template files in the specified directory
func ListTemplatesInDir(dir string) ([]string, error) {
	entries, err := fs.ReadDir(templateFS, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %v", dir, err)
	}

	var templates []string
	for _, entry := range entries {
		if !entry.IsDir() {
			templates = append(templates, path.Join(dir, entry.Name()))
		}
	}

	return templates, nil
}
