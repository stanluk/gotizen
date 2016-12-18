package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Context struct {
	ProjectPath string         /* Absolute path to project's root directory */
	Manifest    *TizenManifest /* Project manifest, may be null */
}

func BuildContext(projectPath string) (*Context, error) {
	var manifest *TizenManifest

	full_path, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to get absolute project path: ", err)
	}
	manifest_full_path := filepath.Join(full_path, TizenManifestDefaultPath)
	file, err := os.Open(manifest_full_path)
	if err == nil {
		manifest, _ = LoadManifest(file)
		file.Close()
	}

	return &Context{ProjectPath: full_path, Manifest: manifest}, nil
}
