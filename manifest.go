package main

import (
	"encoding/xml"
	"fmt"
	"io"
)

const (
	LaunchModeSingle = "single"
)

// tizen-manifest path relative to project directory
var TizenManifestDefaultPath string = "tizen-manifest.xml"

type TizenManifest struct {
	XMLName      xml.Name        `xml:"manifest"`
	PackageName  string          `xml:"package,attr"`
	Api          string          `xml:"api-version,attr"`
	Version      string          `xml:"version,attr"`
	Profile      NameNode        `xml:"profile"`
	UIAppEntries []UIApplication `xml:"ui-application"`
	Privileges   []string        `xml:"privileges>privilege"`
}

type NameNode struct {
	Name string `xml:"name,attr"`
}

type UIApplication struct {
	AppId      string   `xml:"appid,attr"`
	Exec       string   `xml:"exec,attr"`
	LaunchMode string   `xml:"launch_mode,attr"`
	Multiple   bool     `xml:"multiple,attr"`
	NoDisplay  bool     `xml:"nodisplay,attr"`
	TaskManage bool     `xml:"taskmanage,attr"`
	Category   NameNode `xml:"category"`
	Icon       string   `xml:"icon"`
}

func NewTizenManifest(name string) *TizenManifest {
	return &TizenManifest{
		PackageName: fmt.Sprintf("org.tizen.%s", name),
		Api:         "3.0",
		Version:     "0.0.1",
		Profile:     NameNode{"mobile"},
		UIAppEntries: []UIApplication{{AppId: fmt.Sprintf("org.tizen.%s", name), Exec: name,
			LaunchMode: LaunchModeSingle, Multiple: false, NoDisplay: false, TaskManage: true}},
	}
}

func (this *TizenManifest) WriteContent(writer io.Writer) error {
	bytes, err := xml.MarshalIndent(this, "", "  ")
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	if err != nil {
		return fmt.Errorf("Unable to write file content: %v", err)
	}
	return nil
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
