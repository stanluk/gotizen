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
	BinDir    = "bin"
	SharedDir = "shared"
)

func init() {
	packageCmd.Flag.StringVar(&authorCertificateFile, "author-cert", "", "Security profile used to sign package")
	packageCmd.Flag.StringVar(&authorCertificatePass, "author-passwd", "", "Security profile used to sign package")
	packageCmd.Flag.StringVar(&distributorCertificateFile, "dist-cert", "", "Security profile used to sign package")
	packageCmd.Flag.StringVar(&distributorCertificatePass, "dist-passwd", "", "Security profile used to sign package")
}

// makeFileListFromAppsList parses list of Application structure
// and constructs list of DiskFiles
func makeFileListFromAppsList(list []Application) (files []PackageFile, err error) {
	for _, p := range list {
		// 1. Binary files
		if p.Exec != "" {
			df, err := NewDiskFile(path.Join(BinDir, p.Exec))
			if err != nil {
				return nil, err
			}
			files = append(files, df)
		}
		// 2. Icons
		if p.Icon != "" {
			df, err := NewDiskFile(path.Join(SharedDir, p.Icon))
			if err != nil {
				return nil, err
			}
			files = append(files, df)
		}
	}
	return files, nil
}

// makeFileList creates a list of files to be packed into tpk package.
// gotizen do not use any hand-crafted build configuration files like Makfile,
// CMakeList.txt or similar, so only source of information about package files
// is tizen-manifest.xml
func makeFileList(manifest *TizenManifest) (files []PackageFile, err error) {
	list, err := makeFileListFromAppsList(manifest.UIAppEntries)
	if err != nil {
		return nil, err
	}
	files = append(files, list...)
	list, err = makeFileListFromAppsList(manifest.ServiceAppEntries)
	if err != nil {
		return nil, err
	}
	files = append(files, list...)

	// append manifest itself
	files = append(files, manifest)
	return files, nil
}

// writePackageFiles creates new zip package
// and writes raw byte content if files into 'out' writer.
func writePackageFiles(files []PackageFile, out io.Writer) error {
	arch := zip.NewWriter(out)
	defer arch.Close()

	for _, file := range files {
		w, err := arch.Create(file.PackagePath())
		if err != nil {
			return fmt.Errorf("Unable to create archive\n%v", err)
		}
		bytes, err := file.MarshalBinary()
		if err != nil {
			return fmt.Errorf("Unable to get binary data\n%v", err)
		}
		_, err = w.Write(bytes)
		if err != nil {
			return fmt.Errorf("Write failed\n%v", err)
		}
	}
	return nil
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

	all_files, err := makeFileList(context.Manifest)
	if err != nil {
		log.Fatal("Failed to create files list\n", err)
	}

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
		log.Fatalf("Unable to create '%s' file\n%v", zip.Name(), err)
	}
	fmt.Printf("Created %s in %s\n", zip.Name(), context.ProjectRootPath)
}
