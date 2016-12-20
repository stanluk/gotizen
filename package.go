package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

var authorCertificateFile string      // author certificate file
var authorCertificatePass string      // author certificate password
var distributorCertificateFile string // distributor certificate file
var distributorCertificatePass string // distributor certificate passowrd

var packageCmd = &Command{
	Run:       MakePkg,
	Name:      "package",
	Short:     "build & package Tizen project.",
	UsageLine: "",
	Long: `
	`,
}

const (
	BinDir = "bin"
	ResDir = "res"
)

type diskFile struct {
	realPath string // path relative to project root dir
	path     string // path relative to package root dir
	file     *os.File
}

func init() {
	packageCmd.Flag.StringVar(&authorCertificateFile, "author-cert", "", "Security profile used to sign package")
	packageCmd.Flag.StringVar(&authorCertificatePass, "author-passwd", "", "Security profile used to sign package")
	packageCmd.Flag.StringVar(&distributorCertificateFile, "dist-cert", "", "Security profile used to sign package")
	packageCmd.Flag.StringVar(&distributorCertificatePass, "dist-passwd", "", "Security profile used to sign package")
}

func (this *diskFile) Path() string {
	return this.path
}

func (this *diskFile) GetReader() (io.ReadCloser, error) {
	file, err := os.Open(this.realPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to write file content: %v", err)
	}
	return file, nil
}

// create a list of package files described in manifest
// gotizen do not use any hand-creafted build configuration files,
// so only source of information about packages is tizen-manifest.xml
func makeFileList(manifest *TizenManifest) (files []PackageFile) {
	for _, p := range manifest.UIAppEntries {
		// 1. Binary files
		if p.Exec != "" {
			var df diskFile
			df.path = path.Join(BinDir, p.Exec)
			df.realPath = p.Exec
			files = append(files, &df)
		}
		// 2. Icons
		if p.Icon != "" {
			var df diskFile
			df.path = path.Join(ResDir, p.Icon)
			df.realPath = p.Icon
			files = append(files, &df)
		}
	}
	for _, p := range manifest.ServiceAppEntries {
		// 1. Binary files
		if p.Exec != "" {
			var df diskFile
			df.path = path.Join(BinDir, p.Exec)
			df.realPath = p.Exec
			files = append(files, &df)
		}
		// 2. Icons
		if p.Icon != "" {
			var df diskFile
			df.path = path.Join(ResDir, p.Icon)
			df.realPath = p.Icon
			files = append(files, &df)
		}
	}

	// 3. append manifest itself
	files = append(files, manifest)
	return files
}

// writePackageFiles creates new zip package
// and writes raw byte content if files into 'out' writer.
func writePackageFiles(files []PackageFile, out io.Writer) error {
	arch := zip.NewWriter(out)
	for _, file := range files {
		w, err := arch.Create(file.Path())
		if err != nil {
			return fmt.Errorf("Unable to create archive: %v", err)
		}
		reader, err := file.GetReader()
		if err != nil {
			return fmt.Errorf("Unable to get reader %v", err)
		}
		_, err = io.Copy(w, reader)
		if err != nil {
			return fmt.Errorf("Copy failed")
		}
		reader.Close()
		if err != nil {
			return fmt.Errorf("Unable to create archive: %v", err)
		}
	}
	return arch.Close()
}

func createSignature(name SignatureType, files []PackageFile) (*Signature, error) {
	s, err := NewSignature(name, files)
	if err != nil {
		return nil, err
	}

	if name == AuthorSignature {
		s.AuthorCertificate = authorCertificateFile
		s.AuthorPass = authorCertificatePass
	} else if name == DistributorSignature {
		s.AuthorCertificate = distributorCertificateFile
		s.AuthorPass = distributorCertificatePass
	}

	return s, nil
}

func MakePkg(context *Context) {
	if context.Manifest == nil {
		log.Fatal("No manifest file found in ", context.ProjectRootPath)
	}
	zip, err := os.OpenFile(context.Manifest.PackageName+".tpk", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("Unable to create 'tpk' file: ", err)
	}
	defer zip.Close()

	all_files := makeFileList(context.Manifest)

	authorSignature, err := createSignature(AuthorSignature, all_files)
	if err != nil {
		log.Fatal("Unable to sign package, ", err)
	}
	all_files = append(all_files, authorSignature)

	distributorSignature, err := createSignature(DistributorSignature, all_files)
	if err != nil {
		log.Fatal("Unable to sign package, ", err)
	}
	all_files = append(all_files, distributorSignature)

	err = writePackageFiles(all_files, zip)
	if err != nil {
		log.Fatalf("Unable to create '%s' file: %v", zip.Name(), err)
	}
	fmt.Printf("Created %s in %s\n", zip.Name(), context.ProjectRootPath)
}
