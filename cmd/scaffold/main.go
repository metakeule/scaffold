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
	cfg = config.MustNew("scaffold", "1.1",
		`scaffold creates files and directories based on a template and json input.
Complete documentation at http://godoc.org/gopkg.in/metakeule/scaffold.v1`)

	templateArg = cfg.NewString("template", "the file where the template resides", config.Required, config.Shortflag('t'))
	dirArg      = cfg.NewString("dir", "directory that is the target/root of the file creations", config.Default("."))
	headCmd     = cfg.MustCommand("head", "shows the head section of the given template").Skip("dir")
	testCmd     = cfg.MustCommand("test", "makes a test run without creating any files")
)

func main() {

	var (
		err         error
		dir         string
		templateRaw []byte
	)

steps:
	for jump := 1; err == nil; jump++ {
		switch jump - 1 {
		default:
			break steps
		case 0:
			err = cfg.Run()
		case 1:
			dir, err = filepath.Abs(dirArg.Get())
		case 2:
			templateRaw, err = ioutil.ReadFile(templateArg.Get())
		case 3:
			head, template := scaffold.SplitTemplate(string(templateRaw))
			switch cfg.ActiveCommand() {
			case nil:
				err = scaffold.Run(dir, template, os.Stdin, os.Stdout, false)
			case testCmd:
				err = scaffold.Run(dir, template, os.Stdin, os.Stdout, true)
			case headCmd:
				fmt.Fprintln(os.Stdout, head)
			}
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stdout, " -> run 'scaffold help' to get more help")
		os.Exit(1)
	}
}
