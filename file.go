package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type DiskFile struct {
	path string // file path relative to package root dir
}

func checkFile(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("Stat failed on %s\n%v", path, err)
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("%s is not regular file\n%v", path, err)
	}
	return nil
}

func NewDiskFile(filePath string) (*DiskFile, error) {
	err := checkFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Invalid file %s\n%v", filePath, err)
	}
	return &DiskFile{path: filePath}, nil
}

func (diskFile *DiskFile) MarshalBinary() ([]byte, error) {
	err := checkFile(diskFile.path)
	if err != nil {
		return nil, fmt.Errorf("Invalid file %s\n%v", diskFile.path, err)
	}
	file, err := os.Open(diskFile.path)
	if err != nil {
		return nil, fmt.Errorf("Open file %s failed \n%v", diskFile.path, err)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Reading file %s failed \n%v", diskFile.path, err)
	}
	return bytes, nil
}

func (diskFile *DiskFile) PackagePath() string {
	return diskFile.path
}
