package swiftcov

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	r := require.New(t)

	f := &Formatter{FileNames: []string{"./Calculator.swift.gcov"}}
	err := f.Parse()
	r.NoError(err)

	r.Len(f.SourceFiles, 1)

	sf := f.SourceFiles[0]
	r.Equal("./Calculator.swift.gcov", sf.Name)
	r.InDelta(70.8, sf.CoveredPercent, 1)
	r.Len(sf.Coverage, 61)
	r.False(sf.Coverage[15].Valid)
	r.False(sf.Coverage[27].Valid)
	r.True(sf.Coverage[26].Valid)
	r.Equal(0, sf.Coverage[53].Int)
	r.Equal(1, sf.Coverage[48].Int)
	r.Equal(2, sf.Coverage[18].Int)
}
