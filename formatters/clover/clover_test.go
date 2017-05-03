package clover

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	r := require.New(t)

	f := &Formatter{Path: "./example.xml"}
	err := f.Parse()
	r.NoError(err)
	r.Len(f.SourceFiles, 12)

	sf := f.SourceFiles[10]
	r.Equal("/Users/markbates/Dropbox/development/php-test-reporter/src/TestReporter/Entity/CiInfo.php", sf.Name)
	r.InDelta(91.78, sf.CoveredPercent, 1)
}
