package swiftcov

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/markbates/pop/nulls"
	"github.com/pkg/errors"
)

// Formatter collects SwiftCov files, parses them, then formats them into a single report.
type Formatter struct {
	FileNames   []string
	SourceFiles []formatters.SourceFile
}

var searchPaths = []string{"./"}
var search = ".swift.gcov" // look for these file extensions

// Search searches the designated paths for SwiftCov files,
// appending them to the list of filenames.
func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)
	for _, p := range paths {
		logrus.Debugf("checking search path %s for SwiftCov formatter", p)
		files, err := ioutil.ReadDir(p)
		if err != nil {
			return "", errors.WithStack(err)
		}
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), search) {
				f.FileNames = append(f.FileNames, file.Name())
			}
		}
	}

	if len(f.FileNames) == 0 {
		return "",
			errors.WithStack(
				errors.Errorf(
					"could not find any files in search paths for SwiftCov. search paths were: %s",
					strings.Join(paths, ", ")))
	}

	return fmt.Sprint(f.FileNames), nil
}

// Parse parses each file the formatter has found in turn.
func (f *Formatter) Parse() error {
	for _, file := range f.FileNames {
		sf, err := parseSourceFile(file)
		if err != nil {
			return errors.WithStack(err)
		}
		f.SourceFiles = append(f.SourceFiles, *sf)
	}
	return nil
}

// Parse a single Swift source file.
func parseSourceFile(fileName string) (*formatters.SourceFile, error) {
	gitHead, _ := env.GetHead()
	sf, err := formatters.NewSourceFile(fileName, gitHead)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		split := strings.SplitN(string(line), ":", 3)
		if len(split) != 3 {
			return nil, errors.New(
				"SwiftCov file expected to have 3 parts to each line, separated by ':'" + string(line))
		}

		coverage := strings.TrimSpace(split[0])
		lineNum, _ := strconv.Atoi(strings.TrimSpace(split[1]))
		if lineNum < 1 { // pre code metadata
			continue
		}

		switch coverage {
		case "-":
			sf.Coverage = append(sf.Coverage, nulls.Int{})
		case "#####":
			sf.Coverage = append(sf.Coverage, nulls.NewInt(0))
		default: // coverage is number of hits
			num, err := strconv.Atoi(coverage)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			sf.Coverage = append(sf.Coverage, nulls.NewInt(num))
		}

	}

	sf.CalcLineCounts()
	return &sf, nil
}

// Format combines the source files into a report.
func (f *Formatter) Format() (formatters.Report, error) {
	rep, err := formatters.NewReport()
	if err != nil {
		return rep, err
	}

	for _, sf := range f.SourceFiles {
		rep.AddSourceFile(sf)
	}

	return rep, nil
}
