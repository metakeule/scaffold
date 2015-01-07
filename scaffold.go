package scaffold

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var FuncMap = template.FuncMap{
	"replace":    Replace,
	"camelCase1": CamelCase1,
	"camelCase2": CamelCase2,
	"title":      strings.Title,
	"toLower":    strings.ToLower,
	"toUpper":    strings.ToUpper,
	"trim":       strings.Trim,
}

func CamelCase1(src string) string {
	s := strings.Split(src, "_")
	for i, _ := range s {
		s[i] = strings.Title(s[i])
	}
	return strings.Join(s, "")
}

func CamelCase2(src string) string {
	s := strings.Split(src, "_")
	for i, _ := range s {
		if i != 0 {
			s[i] = strings.Title(s[i])
		}
	}
	return strings.Join(s, "")
}

func Replace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

func writeFile(file string, content []byte, log io.Writer) error {
	dir := filepath.Dir(file)
	s, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0770)
	} else {
		if !s.IsDir() {
			return fmt.Errorf("not a directory: %#v", dir)
		}
	}
	if log != nil {
		log.Write([]byte(file + "\n"))
	}

	return ioutil.WriteFile(file, content, 0664)
}

func parseGenerator(startDir string, rd io.Reader, log io.Writer) error {
	scanner := bufio.NewScanner(rd)
	var file string
	var dir = startDir
	var bf bytes.Buffer
	var line = -1
	for scanner.Scan() {
		line++
		s := scanner.Text()
		if strings.HasPrefix(s, ">>>") {
			fd := strings.TrimSpace(strings.TrimPrefix(s, ">>>"))
			if fd[len(fd)-1] == '/' {
				dir = filepath.Join(dir, fd)
			} else {
				file = filepath.Join(dir, fd)
			}
			continue
		}

		if strings.HasPrefix(s, "<<<") {
			fd := strings.TrimSpace(strings.TrimPrefix(s, "<<<"))
			if fd[len(fd)-1] == '/' {
				dirName := filepath.Base(dir) + "/"
				if dirName != fd {
					return fmt.Errorf("syntax error in line %d closing dir %#v but should close dir %#v", line, fd, dirName)
				}
				dir = filepath.Dir(dir)
			} else {
				base := filepath.Base(file)
				if base != fd {
					return fmt.Errorf("syntax error in line %d closing file %#v but should close file %#v", line, fd, base)
				}
				err := writeFile(file, bf.Bytes(), log)
				if err != nil {
					return err
				}
				file = ""
				bf.Reset()
			}
			continue
		}

		bf.WriteString(s + "\n")
		// fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func SplitTemplate(s string) (help, templ string) {
	spl := strings.SplitN(s, "\n\n", 2)
	return spl[0], spl[1]
}

func convertJSON(rd io.Reader) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}
	err = json.NewDecoder(rd).Decode(&data)
	return
}

func substitute(templ string, data map[string]interface{}) (rd io.Reader, err error) {
	var bf bytes.Buffer
	t := template.New("x").Funcs(FuncMap)
	t, err = t.Parse(templ)
	if err != nil {
		return
	}
	err = t.Execute(&bf, data)
	rd = &bf
	return
}

func Run(baseDir string, template string, input io.Reader, log io.Writer) error {

	var (
		err          error
		placeholders map[string]interface{}
		generator    io.Reader
	)

steps:
	for jump := 1; err == nil; jump++ {
		switch jump - 1 {
		default:
			break steps
		case 0:
			placeholders, err = convertJSON(input)
		case 1:
			generator, err = substitute(template, placeholders)
		case 2:
			err = parseGenerator(baseDir, generator, log)
		}
	}
	return err
}
