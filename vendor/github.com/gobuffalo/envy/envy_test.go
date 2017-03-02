package envy_test

import (
	"os"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/require"
)

func Test_Get(t *testing.T) {
	r := require.New(t)
	r.NotZero(os.Getenv("GOPATH"))
	r.Equal(os.Getenv("GOPATH"), envy.Get("GOPATH", "foo"))
	r.Equal("bar", envy.Get("IDONTEXIST", "bar"))
}

func Test_MustGet(t *testing.T) {
	r := require.New(t)
	r.NotZero(os.Getenv("GOPATH"))
	v, err := envy.MustGet("GOPATH")
	r.NoError(err)
	r.Equal(os.Getenv("GOPATH"), v)

	_, err = envy.MustGet("IDONTEXIST")
	r.Error(err)
}

func Test_Set(t *testing.T) {
	r := require.New(t)
	_, err := envy.MustGet("FOO")
	r.Error(err)

	envy.Set("FOO", "foo")
	r.Equal("foo", envy.Get("FOO", "bar"))
}

func Test_MustSet(t *testing.T) {
	r := require.New(t)

	r.Zero(os.Getenv("FOO"))

	err := envy.MustSet("FOO", "BAR")
	r.NoError(err)

	r.Equal("BAR", os.Getenv("FOO"))
}

func Test_Temp(t *testing.T) {
	r := require.New(t)

	_, err := envy.MustGet("BAR")
	r.Error(err)

	envy.Temp(func() {
		envy.Set("BAR", "foo")
		r.Equal("foo", envy.Get("BAR", "bar"))
		_, err = envy.MustGet("BAR")
		r.NoError(err)
	})

	_, err = envy.MustGet("BAR")
	r.Error(err)
}
