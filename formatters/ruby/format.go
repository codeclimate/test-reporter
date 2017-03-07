package ruby

import (
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
)

func (r Formatter) Format() (formatters.Report, error) {
	rep := formatters.Report{
		SourceFiles: []formatters.SourceFile{},
		LineCounts:  formatters.LineCounts{},
	}

	env, err := env.New()
	if err != nil {
		return rep, err
	}
	rep.CIService = env
	rep.Git = env.Git

	var covPer float64
	for _, tt := range r.Tests {
		for _, f := range tt.SourceFiles {
			sf := formatters.SourceFile{
				Name:       f.Name,
				LineCounts: f.LineCounts(),
			}
			for _, i := range f.Coverage {
				sf.Coverage = append(sf.Coverage, i.Interface())
			}
			sf.CoveredPercent = f.CoveragePercent()
			covPer += sf.CoveredPercent
			rep.SourceFiles = append(rep.SourceFiles, sf)
			rep.LineCounts.Covered += sf.LineCounts.Covered
			rep.LineCounts.Missed += sf.LineCounts.Missed
			rep.LineCounts.Total += sf.LineCounts.Total
		}
	}

	rep.CoveredPercent = covPer / float64(len(rep.SourceFiles))

	return rep, nil
}
