package main

import (
	"encoding/xml"
	"io"
	"os"
	"testing"
)

var files []string = []string{
	"bin/lockscreen",
}

func getReader(path string) io.Reader {
	file, err := os.Open(path)
	if err != nil {
		panic("Unable to open file")
	}
	return file
}

func TestSignature(t *testing.T) {
	if err := os.Chdir("signature_test"); err != nil {
		panic("Unable to chdir")
	}

	sig := NewXMLSignature()

	for _, f := range files {
		err := sig.AppendReference(f, getReader(f))
		if err != nil {
			panic("AppendReference failed")
		}
	}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent(" ", "")
	enc.Encode(sig)
}
