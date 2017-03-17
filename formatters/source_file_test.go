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
		Coverage: Coverage{nulls.Int{}, 2, 3, nulls.Int{}},
	}
	b := SourceFile{
		BlobID:   "b",
		Coverage: Coverage{1, nulls.Int{}, 3, nulls.Int{}},
	}

	c, err := a.Merge(b)
	r.NoError(err)
	r.Equal("a", c.BlobID)
	r.Equal(4, len(c.Coverage))
	r.Equal(Coverage{1, 2, 6, nulls.Int{}}, c.Coverage)
	r.Equal(LineCounts{Total: 4, Missed: 1, Covered: 3}, c.LineCounts)
}
