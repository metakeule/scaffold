package scaffold

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var fileVar = regexp.MustCompile("^#([a-zA-Z_]+)$")

// the scanner scans a directory recursively
// and creates a template based on the structure of the files and directories
type scanner struct {
	bf             bytes.Buffer
	skipDirRegex   *regexp.Regexp
	openDirs       []string
	currentDirPath string
}

func (s *scanner) walkDir(path string, info os.FileInfo, err error) error {
	s.closeDir(path)
	if err == filepath.SkipDir {
		return nil
	}

	var nstr = info.Name()

	if s.skipDirRegex != nil && s.skipDirRegex.MatchString(nstr) {
		return filepath.SkipDir
	}

	s.currentDirPath, _ = filepath.Abs(path)

	nstr = fixName(nstr)

	s.openDirs = append(s.openDirs, nstr)

	s.bf.WriteString(fmt.Sprintf("\n>>>%s/\n", nstr))

	return nil
}

func (s *scanner) _closeDir() {
	cdir := s.openDirs[len(s.openDirs)-1]

	s.bf.WriteString(fmt.Sprintf("\n<<<%s/\n", cdir))
	if len(s.openDirs) > 1 {
		s.openDirs = s.openDirs[:len(s.openDirs)-1]
	} else {
		s.openDirs = []string{}
	}
}

// closeDir() closes the last currentdir if needed
func (s *scanner) closeDir(currentFile string) {
	currentFile, _ = filepath.Abs(currentFile)
	dir := filepath.Dir(currentFile)

	// file is not part of the currentDir, that means
	// we left the currentDir, so close the dir and set the currentDirPath accordingly
	if !strings.Contains(dir, s.currentDirPath) {
		s._closeDir()
	}

}

// returns if s starts with an ascii lowercase letter
func isLowercase(s string) bool {
	return s[0] > 90
}

func fixName(in string) string {
	if fileVar.MatchString(in) {
		prefix := "filename"
		in = fileVar.FindString(in)[1:]
		if isLowercase(in) {
			prefix = "filenameLower"
		}
		return fmt.Sprintf("{{%s .%s}}", prefix, CamelCase1(in))
	}
	return in
}

func splitFilename(file string) (withoutext, ext string) {
	idx := strings.LastIndex(file, ".")

	return file[:idx], file[idx:]
}

func (s *scanner) walkFile(path string, info os.FileInfo, err error) error {
	s.closeDir(path)

	if err != nil {
		return err
	}

	bare, ext := splitFilename(info.Name())
	bare = fixName(bare)

	s.bf.WriteString(fmt.Sprintf("\n>>>%s\n", bare+ext))

	fc, err2 := ioutil.ReadFile(path)

	if err2 != nil {
		return err2
	}

	s.bf.Write(fc)

	s.bf.WriteString(fmt.Sprintf("\n<<<%s\n", bare+ext))
	return nil
}

func (s *scanner) walk(path string, info os.FileInfo, err error) error {
	if err == filepath.SkipDir {
		return nil
	}
	if err != nil {
		return err
	}

	if info.IsDir() {
		return s.walkDir(path, info, err)
	}

	return s.walkFile(path, info, err)
}

// scans a directory recursively
// and creates a template based on the structure of the files and directories
func Scan(dirname string, opts ...ScanOption) (template []byte, err error) {
	s := &scanner{}

	for _, opt := range opts {
		opt(s)
	}

	err = filepath.Walk(dirname, s.walk)

	if err == io.EOF || err == filepath.SkipDir {
		err = nil
	}

	if err != nil {
		return nil, err
	}

	for len(s.openDirs) > 0 {
		s._closeDir()
	}

	return s.bf.Bytes(), err
}

type ScanOption func(*scanner)

func SkipDirs(regexstring string) ScanOption {
	re := regexp.MustCompile(regexstring)
	return func(s *scanner) {
		s.skipDirRegex = re
	}
}
