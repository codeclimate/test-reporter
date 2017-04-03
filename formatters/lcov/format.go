package lcov

import "github.com/codeclimate/test-reporter/formatters"

func (r Formatter) Format() (formatters.Report, error) {
	rep, err := formatters.NewReport()
	if err != nil {
		return rep, err
	}

	for _, f := range r.SourceFiles {
		rep.AddSourceFile(f)
	}

	return rep, nil
}
