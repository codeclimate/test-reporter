package ruby

import (
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/markbates/pop/nulls"
)

type SourceFile struct {
	Name     string
	Coverage []nulls.Int
}

func (f SourceFile) LineCounts() formatters.LineCounts {
	lc := formatters.LineCounts{}

	for _, c := range f.Coverage {
		if !c.Valid {
			continue
		}
		lc.Total++
		if c.Int == 0 {
			lc.Missed++
			continue
		}
		lc.Covered++
	}

	return lc
}

func (f SourceFile) CoveragePercent() float64 {
	var lc float64
	var hc float64
	for _, x := range f.Coverage {
		if x.Valid {
			lc++
			if x.Int > 0 {
				hc++
			}
		}
	}
	return (hc / lc) * 100
}
