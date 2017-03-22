package ruby

import (
	"github.com/codeclimate/test-reporter/formatters"
)

func (r Formatter) Format() (formatters.Report, error) {
	rep, err := formatters.NewReport()
	if err != nil {
		return rep, err
	}

	for _, tt := range r.Tests {
		for _, f := range tt.SourceFiles {
			sf := formatters.NewSourceFile(f.Name)
			sf.LineCounts = f.LineCounts()
			sf.Coverage = f.Coverage
			sf.CoveredPercent = f.CoveragePercent()
			rep.AddSourceFile(sf)
		}
	}

	return rep, nil
}
