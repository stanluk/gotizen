package main

import (
	"encoding"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
)

// Command implement basic operation of gotizen tool
type Command struct {
	Run       func(context *Context)
	Flag      flag.FlagSet
	Name      string
	UsageLine string
	Short     string
	Long      string
}

// PackageFile interface should be implemented by a entity
// that can be containted in tizen package.
type PackageFile interface {
	encoding.BinaryMarshaler
	PackagePath() string /* Path relative to package root */
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
			err = cmd.Flag.Parse(args[1:])
			if err != nil {
				fmt.Print(err)
				os.Exit(1)
			}
			cmd.Run(ctx)
			os.Exit(0)
		}
	}

	fmt.Printf("Unknown submcommand: '%s'\n", args[0])
	os.Exit(1)
}
