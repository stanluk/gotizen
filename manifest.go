package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

const (
	LaunchModeSingle = "single"
)

// tizen-manifest path relative to project directory
var TizenManifestDefaultPath string = "tizen-manifest.xml"
var tizenNamespace string = "http://tizen.org/ns/packages"

type TizenManifest struct {
	XMLName      xml.Name        `xml:"manifest"`
	PackageName  string          `xml:"package,attr"`
	Api          string          `xml:"api-version,attr"`
	Version      string          `xml:"version,attr"`
	Profile      NameNode        `xml:"profile"`
	UIAppEntries []UIApplication `xml:"ui-application"`
	Privileges   []string        `xml:"privileges>privilege"`
	XMLNS        string          `xml:"xmlns,attr"`
}

type NameNode struct {
	Name string `xml:"name,attr"`
}

type UIAppType string

// types of UI Applications
var (
	Capp UIAppType = "capp"
)

type UIApplication struct {
	AppId      string    `xml:"appid,attr"`
	Exec       string    `xml:"exec,attr"`
	LaunchMode string    `xml:"launch_mode,attr"`
	Multiple   bool      `xml:"multiple,attr"`
	NoDisplay  bool      `xml:"nodisplay,attr"`
	TaskManage bool      `xml:"taskmanage,attr"`
	Category   NameNode  `xml:"category"`
	Icon       string    `xml:"icon"`
	Type       UIAppType `xml:"type,attr"`
}

func NewTizenManifest(name string) *TizenManifest {
	return &TizenManifest{
		PackageName: fmt.Sprintf("org.tizen.%s", name),
		Api:         "3.0",
		Version:     "0.0.1",
		Profile:     NameNode{"mobile"},
		XMLNS:       tizenNamespace,
		UIAppEntries: []UIApplication{{AppId: fmt.Sprintf("org.tizen.%s", name), Exec: name,
			LaunchMode: LaunchModeSingle, Multiple: false, NoDisplay: false, TaskManage: true, Type: Capp}},
	}
}

type manifestReaderCloser struct {
	*bytes.Buffer
}

func (this *manifestReaderCloser) Close() error {
	return nil
}

func (this *TizenManifest) GetReader() (io.ReadCloser, error) {
	buf, err := xml.MarshalIndent(this, "", "  ")
	if err != nil {
		return nil, err
	}
	return &manifestReaderCloser{Buffer: bytes.NewBuffer(buf)}, nil
}

func (this *TizenManifest) Path() string {
	return TizenManifestDefaultPath
}

func LoadManifest(reader io.Reader) (*TizenManifest, error) {
	var manifest TizenManifest
	dec := xml.NewDecoder(reader)
	err := dec.Decode(&manifest)
	return &manifest, err
}
