package env

import (
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/require"
)

func Test_loadCIFromENV(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		envy.Set("CI_NAME", "travis")
		envy.Set("CI_BUILD_ID", "a12345")
		envy.Set("CI_BUILD_URL", "http://example.com")
		c := loadCIInfo()
		r.Equal(c.Name, "travis")
		r.Equal(c.BuildID, "a12345")
		r.Equal(c.BuildURL, "http://example.com")
	})
}

func Test_loadCIFromENV_Alt_Vars(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		envy.Set("CIRCLECI", "")
		envy.Set("GITLAB_CI", "gitlab")
		envy.Set("CIRCLE_BUILD_NUM", "b12345")
		envy.Set("CIRCLE_BUILD_URL", "http://example.org")
		c := loadCIInfo()
		r.Equal(c.Name, "gitlab")
		r.Equal(c.BuildID, "b12345")
		r.Equal(c.BuildURL, "http://example.org")
	})
}

func Test_CI_String(t *testing.T) {
	r := require.New(t)
	c := ci{
		Name:     "codeclimate",
		BuildID:  "a12345",
		BuildURL: "http://example.net",
	}
	exp := `CI_NAME=codeclimate
CI_BUILD_ID=a12345
CI_BUILD_URL=http://example.net`
	r.Equal(exp, c.String())
}
