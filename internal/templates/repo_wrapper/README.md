# {{.RepoName}} Development Wrapper

This is a development wrapper for the {{.RepoName}} repository.

## Original Repository
- URL: {{.RepoURL}}
- Cloned at: {{.CloneDate}}
- Added as Git submodule in repo/

## Structure
- repo/: Contains the original repository (as a Git submodule)
- local/: Directory for local development files and overrides
- scripts/: Development and utility scripts

## Development
Add your development notes and instructions here.

## Submodule Management
To update the submodule to the latest version:

```bash
cd repo
git pull origin main  # or master, depending on the branch
cd ..
git add repo
git commit -m "Update submodule"
``` 