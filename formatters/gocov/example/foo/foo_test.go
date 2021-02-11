package foo_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/codeclimate/test-reporter/formatters/gocov/example/foo"
)

func Test_String(t *testing.T) {
	foo := &foo.Foo{}
	require.Equal(t, "foo", foo.String())
}
