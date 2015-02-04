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

// FuncMap provides functions to the template.
// New functions can be added as needed. The usual restrictions for
// text/template.FuncMap apply (see http://golang.org/pkg/text/template/#FuncMap)
var FuncMap = template.FuncMap{
	"replace":          Replace,
	"camelCase1":       CamelCase1,
	"camelCase2":       CamelCase2,
	"title":            strings.Title,
	"toLower":          strings.ToLower,
	"toUpper":          strings.ToUpper,
	"trim":             strings.Trim,
	"doubleCurlyOpen":  DoubleCurlyOpen,
	"doubleCurlyClose": DoubleCurlyClose,
	"dollar":           Dollar,
}

// Dollar returns a dollar char
func Dollar() string {
	return "$"
}

// DoubleCurlyOpen returns two open curly braces
func DoubleCurlyOpen() string {
	return "{{"
}

// DoubleCurlyClose returns two closed curly braces
func DoubleCurlyClose() string {
	return "}}"
}

// CamelCase1 converts a string in snake_case to CamelCase where the first letter of each word is capitalized
func CamelCase1(src string) string {
	s := strings.Split(src, "_")
	for i, _ := range s {
		s[i] = strings.Title(s[i])
	}
	return strings.Join(s, "")
}

// CamelCase2 converts a string in snake_case to camelCase where the first letter of each but the first word is capitalized
func CamelCase2(src string) string {
	s := strings.Split(src, "_")
	for i, _ := range s {
		if i != 0 {
			s[i] = strings.Title(s[i])
		}
	}
	return strings.Join(s, "")
}

// Replace replaces every occurence of old in s by new
func Replace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

// writeFile creates the given file with the given content if isTest is true.
// Needed directories are created on the fly and each file name is written to log
// if log is not nil.
// If isTest is true, no files and directories are created.
func writeFile(file string, content []byte, log io.Writer, isTest bool) error {
	dir := filepath.Dir(file)

	if s, err := os.Stat(dir); err != nil || !s.IsDir() {
		if err != nil {
			if os.IsNotExist(err) {
				if !isTest {
					err = os.MkdirAll(dir, 0770)
				} else {
					err = nil
				}
			}
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("not a directory: %#v", dir)
		}
	}

	if log != nil {
		log.Write([]byte(file + "\n"))
	}

	if !isTest {
		return ioutil.WriteFile(file, content, 0664)
	}
	return nil
}

// parseGenerator creates files and directories beneath baseDir as defined in the reader.
// The file names are written to log if it is not nil.
// If isTest is true, no files and directories are created.
func parseGenerator(baseDir string, rd io.Reader, log io.Writer, isTest bool) error {
	scanner := bufio.NewScanner(rd)
	var file string
	var dir = baseDir
	var bf bytes.Buffer
	var line = -1
	for scanner.Scan() {
		line++
		s := scanner.Text()
		if strings.HasPrefix(s, ">>>") {
			fd := strings.TrimSpace(strings.TrimPrefix(s, ">>>"))
			if fd[len(fd)-1] == '/' {
				dir = filepath.Join(dir, fd)
				file = ""
			} else {
				if file != "" {
					return fmt.Errorf("syntax error in line %d embedding file within file is not allowed (%#v inside %#v)", line, fd, file)
				}
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
				err := writeFile(file, bf.Bytes(), log, isTest)
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

// SplitTemplate splits the given template on the first empty line.
// It returns the head and body of the template.
// Templates must be UTF8 without byte order marker and have \n (linefeed) as line terminator.
func SplitTemplate(template string) (head, body string) {
	spl := strings.SplitN(template, "\n\n", 2)
	return spl[0], spl[1]
}

// convertJSON converts json to a map
func convertJSON(rd io.Reader) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}
	err = json.NewDecoder(rd).Decode(&data)
	return
}

// mix mixes the given data to the template body
func mix(body string, data map[string]interface{}) (rd io.Reader, err error) {
	var bf bytes.Buffer
	t := template.New("x").Funcs(FuncMap)
	t, err = t.Parse(body)
	if err != nil {
		return
	}
	err = t.Execute(&bf, data)
	rd = &bf
	return
}

// Run mixes the properties of the json object to the template body. The result is then used
// to create files and directories beneath baseDir.
// If isTest is true the files and directories are not really created.
// If log is not nil a list of files that will be created is written to log.
func Run(baseDir string, body string, json io.Reader, log io.Writer, isTest bool) error {

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
			placeholders, err = convertJSON(json)
		case 1:
			generator, err = mix(body, placeholders)
		case 2:
			err = parseGenerator(baseDir, generator, log, isTest)
		}
	}
	return err
}
