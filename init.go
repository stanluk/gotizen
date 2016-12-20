package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var initCmd = &Command{
	Run:       initProject,
	Name:      "init",
	Short:     "initializes empty Tizen project.",
	UsageLine: "",
	Long: `
	`,
}

// creates all project files
func createFile(rootdir string, file PackageFile) error {
	fullPath := filepath.Join(rootdir, file.PackagePath())
	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("Unable to write %s, aborting\n", file.PackagePath())
	}
	reader, err := file.GetReadCloser()
	if err != nil {
		return fmt.Errorf("Unable to get reader")
	}
	_, err = io.Copy(f, reader)
	if err != nil {
		return fmt.Errorf("Copy failed")
	}
	err = reader.Close()
	fmt.Println("Created: ", file.PackagePath())
	return f.Close()
}

func initProject(context *Context) {
	if context.Manifest != nil {
		log.Fatalf("Tizen mainfest found in %s. Unable to init project.\n", context.ProjectRootPath)
	}
	defaultManifest := NewTizenManifest(filepath.Base(context.ProjectRootPath))

	defaultProjectFiles := make([]PackageFile, 1)
	defaultProjectFiles[0] = defaultManifest
	fmt.Println("Initialized empty Tizen project in: ", context.ProjectRootPath)

	// create project files
	for _, pf := range defaultProjectFiles {
		if err := createFile(context.ProjectRootPath, pf); err != nil {
			log.Fatal("Unable to create Tizen project files: ", err)
		}
	}
}
