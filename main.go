package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"
)

type Command struct {
	Run       func(context *Context)
	Flag      flag.FlagSet
	Name      string
	UsageLine string
	Short     string
	Long      string
}

type PackageFile interface {
	Path() string                   /* Path relative to project root directory */
	WriteContent(w io.Writer) error /* Write package file content */
}

var commands = []*Command{
	initCmd,
	packageCmd,
}

var usageTempl = `gotizen is a tool for bulding and deploying Tizen OS packages.

Usage:

	gotizen command [arguments]

The commands are:
{{range .}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}

`

func usage() {
	t := template.New("tmpl")
	_, err := t.Parse(usageTempl)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(os.Stderr, commands)
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, err := BuildContext(wd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)
	args := flag.Args()

	if len(args) < 1 {
		usage()
		os.Exit(1)
	}
	for _, cmd := range commands {
		if cmd.Name == args[0] {
			cmd.Flag.Parse(args[1:])
			cmd.Run(ctx)
			os.Exit(0)
		}
	}

	fmt.Printf("Unknown submcommand: '%s'\n", args[0])
	os.Exit(1)
}
