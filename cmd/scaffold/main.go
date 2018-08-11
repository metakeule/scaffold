package main

import (
	"bytes"
	"fmt"
	"github.com/metakeule/config"
	"github.com/metakeule/scaffold"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	cfg = config.MustNew("scaffold", "1.6.1",
		`scaffold creates files and directories based on a template and json input.
Complete documentation at http://godoc.org/gopkg.in/metakeule/scaffold.v1`)

	templateArg     = cfg.NewString("template", "the file where the template resides", config.Default("scaffold.template"), config.Shortflag('t'))
	dirArg          = cfg.NewString("dir", "directory that is the target/root of the file creations", config.Default("."))
	templatePathArg = cfg.NewString("path", "the path to look for template files, the different directories must be separated with a colon (:)")
	verboseArg      = cfg.NewBool("verbose", "show verbose messages", config.Default(false), config.Shortflag('v'))

	headCmd    = cfg.MustCommand("head", "shows the head section of the given template").Skip("dir")
	testCmd    = cfg.MustCommand("test", "makes a test run without creating any files")
	scanCmd    = cfg.MustCommand("scan", "scan scans a directory and generates a template based on it. placeholders in dirs and files must start with #").Skip("template").Skip("dir")
	scanDirArg = scanCmd.NewString("scandir", "directory which is scanned to create the template", config.Default("."))

	listCmd = cfg.MustCommand("list", "prints a list of template files, residing in path").Skip("template")
)

type notFound string

func (n notFound) Error() string {
	return fmt.Sprintf("could not find template file %#v", string(n))
}

func printTemplates() {
	paths := strings.Split(templatePathArg.Get(), ":")
	for _, path := range paths {
		if path == "" {
			continue
		}
		fileinfos, err := ioutil.ReadDir(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "skipping %s (missing)\n", path)
			} else {
				fmt.Fprintf(os.Stderr, "skipping %s (%s)\n", path, err)
			}
		} else {
			var bf bytes.Buffer
			for _, fi := range fileinfos {
				if !fi.IsDir() {
					name := fi.Name()
					if name[0] != '.' {
						bf.WriteString("  " + name + "\n")
						// fmt.Fprintln(os.Stdout, name)
					}
				}
			}

			if bf.String() == "" {
				fmt.Fprintf(os.Stdout, "no templates inside %s\n", path)
			} else {
				fmt.Fprintf(os.Stdout, "templates inside %s:\n%s\n", path, bf.String())
			}
		}
	}

}

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
	paths := append([]string{""}, strings.Split(templatePathArg.Get(), ":")...)

	file := templateArg.Get()

	for _, p := range paths {
		if findInDir(p, file) {
			return filepath.Join(p, file), nil
		}
		if findInDir(p, file+".template") {
			return filepath.Join(p, file+".template"), nil
		}
	}
	return "", notFound(file)
}

func main() {

	var (
		err         error
		dir         string
		scanDir     string
		file        string
		templateRaw []byte
		templ       []byte
	)

steps:
	for jump := 1; err == nil; jump++ {
		switch jump - 1 {
		default:
			break steps
		case 0:
			err = cfg.Run()
		case 1:
			if cfg.ActiveCommand() == scanCmd {
				scanDir, err = filepath.Abs(scanDirArg.Get())
			}
		case 2:
			if cfg.ActiveCommand() == scanCmd {
				templ, err = scaffold.Scan(scanDir)
				if err == nil {
					fmt.Fprintln(os.Stdout, string(templ))
					os.Exit(0)
				}
				break steps
			}
		case 3:
			if cfg.ActiveCommand() == listCmd {
				printTemplates()
				os.Exit(0)
			}
		case 4:
			dir, err = filepath.Abs(dirArg.Get())
		case 5:
			file, err = findFile()
		case 6:
			println("found ", file)
			templateRaw, err = ioutil.ReadFile(file)
		case 7:
			head, template := scaffold.SplitTemplate(string(templateRaw))
			switch cfg.ActiveCommand() {
			case nil:
				err = scaffold.Run(dir, template, os.Stdin, os.Stdout, false)
			case testCmd:
				err = scaffold.Run(dir, template, os.Stdin, os.Stdout, true)
			case headCmd:
				fmt.Fprintln(os.Stdout, head)
			default:
				panic("unreachable")
			}
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		fmt.Fprintln(os.Stdout, "\n\n--------------------------------------\n\n"+cfg.Usage())
		os.Exit(1)
	}
}
