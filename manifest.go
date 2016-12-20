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
var TizenManifestDefaultPath = "tizen-manifest.xml"
var tizenNamespace = "http://tizen.org/ns/packages"

type TizenManifest struct {
	XMLName           xml.Name      `xml:"manifest"`
	PackageName       string        `xml:"package,attr"`
	Api               string        `xml:"api-version,attr"`
	Version           string        `xml:"version,attr"`
	Profile           NameAttr      `xml:"profile"`
	UIAppEntries      []Application `xml:"ui-application"`
	ServiceAppEntries []Application `xml:"service-application"`
	Privileges        []string      `xml:"privileges>privilege"`
	XMLNS             string        `xml:"xmlns,attr"`
}

type NameAttr struct {
	Name string `xml:"name,attr"`
}

type ValueAttr struct {
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
	Category           NameAttr  `xml:"category,omitempty"`
	Icon               string    `xml:"icon,omitempty"`
	Type               UIAppType `xml:"type,attr"`
	BackgroundCategory ValueAttr `xml:"background-category"`
}

func NewTizenManifest(name string) *TizenManifest {
	return &TizenManifest{
		PackageName: fmt.Sprintf("org.tizen.%s", name),
		Api:         "3.0",
		Version:     "0.0.1",
		Profile:     NameAttr{"mobile"},
		XMLNS:       tizenNamespace,
		ServiceAppEntries: []Application{{AppId: fmt.Sprintf("org.tizen.%s", name), Exec: name,
			LaunchMode: LaunchModeSingle, Multiple: false, NoDisplay: false, TaskManage: true, Type: Capp, BackgroundCategory: ValueAttr{"system"}}},
	}
}

func (manif *TizenManifest) MarshalBinary() ([]byte, error) {
	buf, err := xml.MarshalIndent(manif, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("MarshalIndent failed\n%v", err)
	}
	return buf, nil
}

func LoadManifest(reader io.Reader) (*TizenManifest, error) {
	manif := &TizenManifest{}
	dec := xml.NewDecoder(reader)
	return manif, dec.Decode(manif)
}

func (manif *TizenManifest) PackagePath() string {
	return TizenManifestDefaultPath
}
