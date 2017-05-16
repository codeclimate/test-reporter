package gcov

import (
	"testing"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	r := require.New(t)

	f := &Formatter{}
	files, err := f.Search("examples")
	t.Log(files)
	if err != nil {
		t.Fatal(err)
	}
	err = f.Parse()
	r.NoError(err)
	r.Len(f.SourceFiles, 3)

	testCalculator(r, f.SourceFiles[0])
	testHamming(r, f.SourceFiles[2])
	testReport(r, f)
}

func testCalculator(r *require.Assertions, sf formatters.SourceFile) {
	r.Equal("examples/Calculator.swift.gcov", sf.Name)
	r.InDelta(70.8, sf.CoveredPercent, 1)
	r.Len(sf.Coverage, 61)
	r.False(sf.Coverage[15].Valid)
	r.False(sf.Coverage[27].Valid)
	r.True(sf.Coverage[26].Valid)
	r.Equal(0, sf.Coverage[53].Int)
	r.Equal(1, sf.Coverage[48].Int)
	r.Equal(2, sf.Coverage[18].Int)
}

func testHamming(r *require.Assertions, sf formatters.SourceFile) {
	r.Equal("examples/hamming.c.gcov", sf.Name)
	r.InDelta(83.3, sf.CoveredPercent, 1)
	r.Len(sf.Coverage, 25)
	r.False(sf.Coverage[2].Valid)
	r.True(sf.Coverage[5].Valid)
	r.Equal(0, sf.Coverage[10].Int)
	r.Equal(2, sf.Coverage[13].Int)
}

func testReport(r *require.Assertions, f *Formatter) {
	report, _ := f.Format()
	r.InDelta(71.7, report.CoveredPercent, 1)
}
