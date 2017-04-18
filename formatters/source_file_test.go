package formatters

import (
	"testing"

	"github.com/markbates/pop/nulls"
	"github.com/stretchr/testify/require"
)

func Test_SourceFile_Merge(t *testing.T) {
	r := require.New(t)
	a := SourceFile{
		BlobID:   "a",
		Coverage: Coverage{nulls.Int{}, nulls.NewInt(2), nulls.NewInt(3), nulls.NewInt(0)},
	}
	b := SourceFile{
		BlobID:   "b",
		Coverage: Coverage{nulls.NewInt(1), nulls.Int{}, nulls.NewInt(3), nulls.Int{}},
	}

	c, err := a.Merge(b)
	r.NoError(err)
	r.Equal("a", c.BlobID)
	r.Equal(4, len(c.Coverage))
	r.InDelta(75.0, c.CoveredPercent, 1)
	r.Equal(LineCounts{Total: 4, Missed: 1, Covered: 3}, c.LineCounts)
}

func Test_SourceFile_BlobID(t *testing.T) {
	r := require.New(t)
	sf, err := NewSourceFile("./coverage.go")
	r.NoError(err)
	r.NotZero(sf.BlobID)
	r.NotContains(sf.BlobID, "blob")
}
