package upload

import "github.com/codeclimate/test-reporter/formatters"

func NewTestReport(rep formatters.Report) *TestReport {
	tr := &TestReport{
		Data: newData(rep),
	}
	return tr
}

type TestReport struct {
	Data Data `json:"data"`
}
