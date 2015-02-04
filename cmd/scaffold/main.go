package main

import (
	"errors"
	"fmt"
	"gopkg.in/metakeule/config.v1"
	"gopkg.in/metakeule/scaffold.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	cfg = config.MustNew("scaffold", "1.4",
		`scaffold creates files and directories based on a template and json input.
Complete documentation at http://godoc.org/gopkg.in/metakeule/scaffold.v1`)

	templateArg     = cfg.NewString("template", "the file where the template resides", config.Required, config.Shortflag('t'))
	templatePathArg = cfg.NewString("path", "the path to look for template files, the different directories must be separated with a semicolon (;)")
	dirArg          = cfg.NewString("dir", "directory that is the target/root of the file creations", config.Default("."))
	verboseArg      = cfg.NewBool("verbose", "show verbose messages", config.Default(false), config.Shortflag('v'))
	headCmd         = cfg.MustCommand("head", "shows the head section of the given template").Skip("dir")
	testCmd         = cfg.MustCommand("test", "makes a test run without creating any files")

	FileNotFound = errors.New("template file not found")
)

func findInDir(path, file string) bool {
	if verboseArg.Get() {
		println("looking for ", filepath.Join(path, file))
	}
	fullPath, err := filepath.Abs(filepath.Join(path, file))
	if err != nil {
		return false
	}

	var info os.FileInfo
	info, err = os.Stat(fullPath)

	if err != nil {
		return false
	}

	return !info.IsDir()
}

// findFile finds the file inside the given path and returns the found file or an error
func findFile() (fullPath string, err error) {
	paths := append([]string{""}, strings.Split(templatePathArg.Get(), ";")...)

	file := templateArg.Get()

	for _, p := range paths {
		if findInDir(p, file) {
			return filepath.Join(p, file), nil
		}
		if findInDir(p, file+".template") {
			return filepath.Join(p, file+".template"), nil
		}
	}
	return "", FileNotFound
}

func main() {

	var (
		err         error
		dir         string
		file        string
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
			file, err = findFile()
		case 3:
			println("found ", file)
			templateRaw, err = ioutil.ReadFile(file)
		case 4:
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
