package simplecov

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
			rep.AddSourceFile(f)
		}
	}

	return rep, nil
}
