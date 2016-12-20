package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Context holds basic information about package being build
type Context struct {
	ProjectRootPath string         /* Absolute path to project's root directory */
	Manifest        *TizenManifest /* Project manifest, may be null */
}

// BuildContext builds context structure using projectPath directory
func BuildContext(projectPath string) (*Context, error) {
	fullPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to get absolute project path: %v", err)
	}
	manifestFullPath := filepath.Join(fullPath, TizenManifestDefaultPath)
	file, err := os.Open(manifestFullPath)
	if err != nil {
		return nil, fmt.Errorf("Opening manifest failed\n%v", err)
	}
	defer file.Close()

	manifest, err := LoadManifest(file)
	if err != nil {
		return nil, fmt.Errorf("Unmarshaling manifest failed\n%v", err)
	}

	return &Context{ProjectRootPath: fullPath, Manifest: manifest}, nil
}
