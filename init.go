package main

import (
	"fmt"
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

var projectDirs = []string{
	"bin",
	"shared",
}

// creates all project files
func createFile(path string, file PackageFile) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Unable to write %s, aborting\n", file.PackagePath())
	}
	bytes, err := file.MarshalBinary()
	if err != nil {
		f.Close()
		return fmt.Errorf("Marshaling %s failed.\n%v", file.PackagePath(), err)
	}
	_, err = f.Write(bytes)
	if err != nil {
		f.Close()
		return fmt.Errorf("File write failed")
	}
	return f.Close()
}

func initProject(context *Context) {
	if context.Manifest != nil {
		log.Fatalf("Tizen mainfest found in %s. Unable to init project.\n", context.ProjectRootPath)
	}
	defaultManifest := NewTizenManifest(filepath.Base(context.ProjectRootPath))

	defaultProjectFiles := make([]PackageFile, 1)
	defaultProjectFiles[0] = defaultManifest

	// create project dirs
	for _, dir := range projectDirs {
		fullPath := filepath.Join(context.ProjectRootPath, dir)
		err := os.Mkdir(fullPath, 0777)
		if err != nil {
			log.Fatalf("Unable to create Tizen project directories\n%v", err)
		}
	}
	// create project files
	for _, pf := range defaultProjectFiles {
		fullPath := filepath.Join(context.ProjectRootPath, pf.PackagePath())
		if err := createFile(fullPath, pf); err != nil {
			log.Fatal("Unable to create Tizen project files: ", err)
		}
	}
	fmt.Println("Initialized empty Tizen project in: ", context.ProjectRootPath)
}
