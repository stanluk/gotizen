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

// validates if given directory contains Tizen project files
func hasTizenProjectFiles(rootPath string, files []PackageFile) bool {
	for _, pf := range files {
		full_path := filepath.Join(rootPath, pf.Path())
		if _, err := os.Stat(full_path); err == nil || !os.IsNotExist(err) {
			return true
		}
	}
	return false
}

// creates all project files
func createFile(rootdir string, file PackageFile) error {
	full_path := filepath.Join(rootdir, file.Path())
	f, err := os.Create(full_path)
	if err != nil {
		return fmt.Errorf("Unable to write %s, aborting\n", file.Path())
	}
	reader, err := file.GetReader()
	if err != nil {
		return fmt.Errorf("Unable to get reader")
	}
	_, err = io.Copy(f, reader)
	if err != nil {
		return fmt.Errorf("Copy failed")
	}
	reader.Close()
	fmt.Println("Created: ", file.Path())
	return f.Close()
}

func initProject(context *Context) {
	if context.Manifest != nil {
		log.Fatalf("Tizen mainfest found in %s. Unable to init project.\n", context.ProjectPath)
	}
	defaultManifest := NewTizenManifest(filepath.Base(context.ProjectPath))

	defaultProjectFiles := make([]PackageFile, 1)
	defaultProjectFiles[0] = defaultManifest
	fmt.Println("Initialized empty Tizen project in: ", context.ProjectPath)

	// create project files
	for _, pf := range defaultProjectFiles {
		if err := createFile(context.ProjectPath, pf); err != nil {
			log.Fatal("Unable to create Tizen project files: ", err)
		}
	}
}
