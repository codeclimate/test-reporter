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
		Coverage: Coverage{1, nulls.Int{}, 3, nulls.Int{}, 5},
	}

	c := a.Merge(b)
	r.Equal("b", c.BlobID)
	r.Equal(5, len(c.Coverage))
	r.Equal(Coverage{1, 2, 6, nulls.Int{}, 5}, c.Coverage)
	r.Equal(LineCounts{Total: 5, Missed: 1, Covered: 4}, c.LineCounts)
}
