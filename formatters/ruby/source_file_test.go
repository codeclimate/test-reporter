package ruby

import (
	"testing"

	"github.com/markbates/pop/nulls"
	"github.com/stretchr/testify/require"
)

var sample = SourceFile{
	Coverage: []nulls.Int{
		nulls.NewInt(1),
		nulls.Int{},
		nulls.NewInt(1),
		nulls.NewInt(1),
		nulls.NewInt(0),
		nulls.Int{},
		nulls.Int{},
		nulls.Int{},
		nulls.NewInt(1),
		nulls.Int{},
		nulls.Int{},
		nulls.Int{},
		nulls.NewInt(1),
		nulls.Int{},
		nulls.Int{},
		nulls.Int{},
		nulls.NewInt(1),
		nulls.Int{},
		nulls.Int{},
		nulls.Int{},
		nulls.Int{},
	},
}

func Test_SourceFile_CoveragePercent(t *testing.T) {
	r := require.New(t)

	per := sample.CoveragePercent()
	r.InDelta(85.71, per, 1.0)
}

func Test_SourceFile_LineCounts(t *testing.T) {
	r := require.New(t)

	counts := sample.LineCounts()
	r.Equal(counts.Covered, 6)
	r.Equal(counts.Missed, 1)
	r.Equal(counts.Total, 7)
}
