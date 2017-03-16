package formatters

import (
	"testing"

	"github.com/markbates/pop/nulls"
	"github.com/stretchr/testify/require"
)

func Test_Report_Merge(t *testing.T) {
	r := require.New(t)
	a := &Report{
		CoveredPercent: 62.5,
		SourceFiles:    SourceFiles{},
	}
	a.AddSourceFile(SourceFile{
		Name:     "a.go",
		Coverage: Coverage{nulls.Int{}, 2, 3, nulls.Int{}},
	})
	a.AddSourceFile(SourceFile{
		Name:     "b.go",
		Coverage: Coverage{1, 2, 3, nulls.Int{}},
	})

	b := &Report{
		CoveredPercent: 50,
		LineCounts:     LineCounts{Missed: 2, Covered: 2, Total: 4},
		SourceFiles:    SourceFiles{},
	}
	b.AddSourceFile(SourceFile{
		Name:     "b.go",
		Coverage: Coverage{1, nulls.Int{}, 3, nulls.Int{}},
	})

	c := &Report{
		CoveredPercent: 66.6,
		LineCounts:     LineCounts{Missed: 2, Covered: 4, Total: 6},
		SourceFiles:    SourceFiles{},
	}
	c.AddSourceFile(SourceFile{
		Name:     "b.go",
		Coverage: Coverage{1, 2, 3, nulls.Int{}},
	})
	c.AddSourceFile(SourceFile{
		Name:     "c.go",
		Coverage: Coverage{nulls.Int{}, nulls.Int{}},
	})
	a.Merge(b, c)
	// 3 files
	r.Equal(3, len(a.SourceFiles))
	sf := a.SourceFiles["b.go"]
	r.NotNil(sf)
	lc := sf.LineCounts
	r.Equal(3, lc.Covered)
	r.Equal(1, lc.Missed)
	r.Equal(4, lc.Total)

	r.InDelta(75, lc.CoveredPercent(), 1.0)
}
