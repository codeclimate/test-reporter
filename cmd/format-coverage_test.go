package cmd

import (
	"testing"

	"github.com/codeclimate/test-reporter/formatters/clover"
	"github.com/stretchr/testify/require"
)

func Test_CoverageFormatter_Save_Fail_On_Empty_Report(t *testing.T) {
	r := require.New(t)
	a := CoverageFormatter{}
	a.In = &clover.Formatter{}
	a.In.Search("../formatters/clover/empty_example.xml")

	err := a.Save()
	r.Error(err)
	r.Equal("could not find coverage info for source files", err.Error())
}
