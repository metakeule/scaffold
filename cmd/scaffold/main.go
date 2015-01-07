package main

import (
	"fmt"
	"github.com/metakeule/scaffold"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {

	var (
		err            error
		workingDir     string
		templateRaw    []byte
		help, template string
	)

steps:
	for jump := 1; err == nil; jump++ {
		switch jump - 1 {
		default:
			break steps
		case 0:
			if len(os.Args) < 2 {
				err = fmt.Errorf("missing template argument")
			}
		case 1:
			workingDir, err = os.Getwd()
		case 2:
			workingDir, err = filepath.Abs(workingDir)
		case 3:
			templateRaw, err = ioutil.ReadFile(os.Args[1])
		case 4:
			help, template = scaffold.SplitTemplate(string(templateRaw))
			if len(os.Args) > 2 && os.Args[2] == "help" {
				fmt.Fprintln(os.Stdout, help)
				os.Exit(0)
			}
		case 5:
			err = scaffold.Run(workingDir, template, os.Stdin, os.Stdout)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, help)
		os.Exit(1)
	}
}
