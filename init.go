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
func createProjectFiles(ctx *Context, files []PackageFile) error {
	for _, pf := range files {
		full_path := filepath.Join(ctx.ProjectPath, pf.Path())
		file, err := os.Create(full_path)
		if err != nil {
			return fmt.Errorf("Unable to write %s, aborting\n", pf.Path())
		}
		reader, err := pf.GetReader()
		if err != nil {
			return fmt.Errorf("Unable to get reader")
		}
		_, err = io.Copy(file, reader)
		if err != nil {
			return fmt.Errorf("Copy failed")
		}
		reader.Close()
		fmt.Println("Created: ", pf.Path())
		file.Close()
	}
	return nil
}

func initProject(context *Context) {
	if context.Manifest != nil {
		log.Fatalf("Tizen mainfest found in %s. Unable to init project.\n", context.ProjectPath)
	}
	defaultManifest := NewTizenManifest(filepath.Base(context.ProjectPath))

	defaultProjectFiles := make([]PackageFile, 1)
	defaultProjectFiles[0] = defaultManifest
	fmt.Println("Initialized empty Tizen project in: ", context.ProjectPath)

	if err := createProjectFiles(context, defaultProjectFiles); err != nil {
		log.Fatal("Unable to create Tizen project files: ", err)
	}
}
