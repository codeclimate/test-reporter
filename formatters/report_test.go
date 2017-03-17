package formatters

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/markbates/pop/nulls"
	"github.com/stretchr/testify/require"
)

func Test_Report_Merge_Bad_GitHead(t *testing.T) {
	r := require.New(t)
	a := &Report{
		Git: ccGit{
			Head: "a",
		},
	}
	b := &Report{
		Git: ccGit{
			Head: "b",
		},
	}
	err := a.Merge(b)
	r.Error(err)
	r.Equal("git heads do not match", err.Error())
}

func Test_Report_Merge_MismatchedCoverageLength(t *testing.T) {
	r := require.New(t)
	a := &Report{
		Git: ccGit{
			Head: "a",
		},
		SourceFiles: SourceFiles{
			"a.go": {
				Name:     "a.go",
				Coverage: Coverage{1},
			},
		},
	}
	b := &Report{
		Git: ccGit{
			Head: "a",
		},
		SourceFiles: SourceFiles{
			"a.go": {
				Name:     "a.go",
				Coverage: Coverage{1, 2},
			},
		},
	}
	err := a.Merge(b)
	r.Error(err)
	r.Equal("coverage length mismatch for a.go", err.Error())
}

func Test_Report_Merge(t *testing.T) {
	r := require.New(t)
	a := &Report{
		Git: ccGit{
			Head: "a",
		},
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
		Git: ccGit{
			Head: "a",
		},
		CoveredPercent: 50,
		LineCounts:     LineCounts{Missed: 2, Covered: 2, Total: 4},
		SourceFiles:    SourceFiles{},
	}
	b.AddSourceFile(SourceFile{
		Name:     "b.go",
		Coverage: Coverage{1, nulls.Int{}, 3, nulls.Int{}},
	})

	c := &Report{
		Git: ccGit{
			Head: "a",
		},
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

func Test_Report_JSON_Unmarshal(t *testing.T) {
	r := require.New(t)
	f, err := os.Open("../examples/codeclimate.json")
	r.NoError(err)

	rep, err := NewReport()
	r.NoError(err)
	err = json.NewDecoder(f).Decode(&rep)
	r.NoError(err)

	r.Equal(20, len(rep.SourceFiles))
	r.Equal("/go/src/github.com/codeclimate/test-reporter/ruby-test-reporter", rep.Environment.PWD)

	sf := rep.SourceFiles["lib/code_climate/test_reporter/client.rb"]
	r.NotNil(sf)
	r.InDelta(87.87, sf.CoveredPercent, 1)

	lc := sf.LineCounts
	r.Equal(8, lc.Missed)
	r.Equal(58, lc.Covered)
	r.Equal(66, lc.Total)
}
