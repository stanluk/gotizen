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
var TizenManifestDefaultPath = "tizen-manifest.xml"
var tizenNamespace = "http://tizen.org/ns/packages"

type TizenManifest struct {
	XMLName           xml.Name      `xml:"manifest"`
	PackageName       string        `xml:"package,attr"`
	Api               string        `xml:"api-version,attr"`
	Version           string        `xml:"version,attr"`
	Profile           NameNode      `xml:"profile"`
	UIAppEntries      []Application `xml:"ui-application"`
	ServiceAppEntries []Application `xml:"service-application"`
	Privileges        []string      `xml:"privileges>privilege"`
	XMLNS             string        `xml:"xmlns,attr"`
}

type NameNode struct {
	Name string `xml:"name,attr"`
}

type ValueNode struct {
	Name string `xml:"value,attr"`
}

type UIAppType string

// types of UI Applications
var (
	Capp UIAppType = "capp"
)

type Application struct {
	AppId              string    `xml:"appid,attr"`
	Exec               string    `xml:"exec,attr"`
	LaunchMode         string    `xml:"launch_mode,attr"`
	Multiple           bool      `xml:"multiple,attr"`
	NoDisplay          bool      `xml:"nodisplay,attr"`
	TaskManage         bool      `xml:"taskmanage,attr"`
	Category           NameNode  `xml:"category,omitempty"`
	Icon               string    `xml:"icon,omitempty"`
	Type               UIAppType `xml:"type,attr"`
	BackgroundCategory ValueNode `xml:"background-category"`
}

func NewTizenManifest(name string) *TizenManifest {
	return &TizenManifest{
		PackageName: fmt.Sprintf("org.tizen.%s", name),
		Api:         "3.0",
		Version:     "0.0.1",
		Profile:     NameNode{"mobile"},
		XMLNS:       tizenNamespace,
		ServiceAppEntries: []Application{{AppId: fmt.Sprintf("org.tizen.%s", name), Exec: name,
			LaunchMode: LaunchModeSingle, Multiple: false, NoDisplay: false, TaskManage: true, Type: Capp, BackgroundCategory: ValueNode{"system"}}},
	}
}

type manifestReaderCloser struct {
	*bytes.Buffer
}

func (manif *manifestReaderCloser) Close() error {
	return nil
}

func (manif *TizenManifest) GetReadCloser() (io.ReadCloser, error) {
	buf, err := xml.MarshalIndent(manif, "", "  ")
	if err != nil {
		return nil, err
	}
	return &manifestReaderCloser{Buffer: bytes.NewBuffer(buf)}, nil
}

func (manif *TizenManifest) PackagePath() string {
	return TizenManifestDefaultPath
}

// LoadManifest constructs TizenManifest structure form binary
// xml representation
func LoadManifest(reader io.Reader) (*TizenManifest, error) {
	var manifest TizenManifest
	dec := xml.NewDecoder(reader)
	err := dec.Decode(&manifest)
	return &manifest, err
}
