package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
)

var packageCmd = &Command{
	Run:       MakePkg,
	Name:      "package",
	Short:     "build & package Tizen project.",
	UsageLine: "",
	Long: `
	`,
}

type diskFile struct {
	path string
}

func (this *diskFile) Path() string {
	return this.path
}

func (this *diskFile) WriteContent(w io.Writer) error {
	file, err := os.Open(this.path)
	if err != nil {
		return fmt.Errorf("Unable to write file content: %v", err)
	}
	_, err = io.Copy(w, file)
	file.Close()
	return err
}

// create a list of information for package described in
// in context parameter
func listPackageFiles(context *Context) ([]PackageFile, error) {
	var files []PackageFile

	if context.Manifest == nil {
		return nil, fmt.Errorf("Unable to query manifest.")
	}

	// 1. Manifest
	files = append(files, context.Manifest)

	for _, p := range context.Manifest.UIAppEntries {
		// 2. Binary files
		if p.Exec != "" {
			files = append(files, &diskFile{p.Exec})
		}
		// 3. Icons
		if p.Icon != "" {
			files = append(files, &diskFile{p.Exec})
		}
	}
	return files, nil
}

// createZipArchive creates new zip package from files data
// and writes raw byte content int 'out' writer.
func createZipArchive(files []PackageFile, out io.Writer) error {
	arch := zip.NewWriter(out)
	for _, file := range files {
		w, err := arch.Create(file.Path())
		if err != nil {
			return fmt.Errorf("Unable to create archive: %v", err)
		}
		err = file.WriteContent(w)
		if err != nil {
			return fmt.Errorf("Unable to create archive: %v", err)
		}
		if err != nil {
			return fmt.Errorf("Unable to create archive: %v", err)
		}
	}
	return arch.Close()
}

func MakePkg(context *Context) {
	if context.Manifest == nil {
		log.Fatal("No manifest file found in ", context.ProjectPath)
	}
	archFiles, err := listPackageFiles(context)
	if err != nil {
		log.Fatal(err)
	}
	zip, err := os.OpenFile("package.tpk", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("Unable to create 'tpk' file: ", err)
	}
	defer zip.Close()
	err = createZipArchive(archFiles, zip)
	if err != nil {
		log.Fatal("Unable to create 'tpk' file: ", err)
	}
	fmt.Println("Created package.tpk in ", context.ProjectPath)
}
