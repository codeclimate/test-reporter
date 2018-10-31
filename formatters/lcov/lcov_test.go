package lcov

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/codeclimate/test-reporter/env"
	"github.com/stretchr/testify/require"
)

func Test_Formatter_Parse(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	r := require.New(t)
	l := Formatter{Path: "./example.info"}
	rep, err := l.Format()
	r.NoError(err)

	r.Len(rep.SourceFiles, 1)
	sf := rep.SourceFiles["/Users/markbates/Dropbox/development/javascript-test-reporter/formatter.js"]
	r.Len(sf.Coverage, 104)
}

func Test_Format(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	r := require.New(t)

	rb := Formatter{
		Path: "./example.info",
	}
	rep, err := rb.Format()
	r.NoError(err)

	r.InDelta(90.38, rep.CoveredPercent, 1)

	sf := rep.SourceFiles["/Users/markbates/Dropbox/development/javascript-test-reporter/formatter.js"]
	r.InDelta(90.19, sf.CoveredPercent, 1)

	lc := rep.LineCounts
	r.Equal(47, lc.Covered)
	r.Equal(5, lc.Missed)
	r.Equal(52, lc.Total)
}

func Benchmark_Format(b *testing.B) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	r := require.New(b)
	inFile, err := genCoverage(1000, 1000)
	r.NoError(err)
	defer os.Remove(inFile.Name())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb := Formatter{
			Path: inFile.Name(),
		}
		rep, err := rb.Format()
		r.NoError(err)
		r.Equal(len(rep.SourceFiles), 1000)
	}
}

func genCoverage(files int, linesPerFile int) (*os.File, error) {
	var buffer bytes.Buffer

	buffer.WriteString("TN:\n")
	for i := 0; i < files; i++ {
		sfLine := fmt.Sprintf("SF:source-file-%v\n", i)
		buffer.WriteString(sfLine)
		for j := 0; j < linesPerFile; j++ {
			lineHits := rand.Intn(42)
			daLine := fmt.Sprintf("DA:%v,%v\n", j+1, lineHits)
			buffer.WriteString(daLine)
		}
		buffer.WriteString("end_of_record\n")
	}

	fh, err := ioutil.TempFile("", "lcov_test")
	if err == nil {
		fh.Write(buffer.Bytes())
	}
	return fh, err
}
