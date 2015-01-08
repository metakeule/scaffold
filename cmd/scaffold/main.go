package main

import (
	"fmt"
	"gopkg.in/metakeule/config.v1"
	"gopkg.in/metakeule/scaffold.v1"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	cfg = config.MustNew("scaffold", "1.0",
		`scaffold creates files and directories based on a template and json input.
		Complete documentation at http://godoc.org/gopkg.in/metakeule/scaffold.v0.`)

	templateArg = cfg.NewString("template", "the file where the template resides", config.Required, config.Shortflag('t'))
	dirArg      = cfg.NewString("dir", "directory that is the target/root of the file creations", config.Default("."))
	exampleCmd  = cfg.MustCommand("example", "shows the example section of the given template")
)

func main() {

	var (
		err               error
		dir               string
		templateRaw       []byte
		example, template string
	)

steps:
	for jump := 1; err == nil; jump++ {
		switch jump - 1 {
		default:
			break steps
		case 0:
			err = cfg.Run()
		case 1:
			dir = dirArg.Get()
		case 2:
			dir, err = filepath.Abs(dir)
		case 3:
			templateRaw, err = ioutil.ReadFile(templateArg.Get())
		case 4:
			example, template = scaffold.SplitTemplate(string(templateRaw))
		case 5:
			switch cfg.ActiveCommand() {
			case nil:
				err = scaffold.Run(dir, template, os.Stdin, os.Stdout)
			case exampleCmd:
				fmt.Fprintln(os.Stdout, example)
			}
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
