package upload

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TransformPostBatchURL_Secure(t *testing.T) {
	r := require.New(t)

	uploader := Uploader{
		Insecure: false,
	}

	rawURL := "https://example.com/"
	actualURL, err := uploader.TransformPostBatchURL(rawURL)

	r.Equal("https://example.com/", actualURL)
	r.Nil(err)
}
func Test_TransformPostBatchURL_Insecure_Success(t *testing.T) {
	r := require.New(t)

	uploader := Uploader{
		Insecure: true,
	}

	rawURL := "https://example.com/"
	actualURL, err := uploader.TransformPostBatchURL(rawURL)

	r.Equal("http://example.com/", actualURL)
	r.Nil(err)
}

func Test_TransformPostBatchURL_Insecure_Error(t *testing.T) {
	r := require.New(t)

	uploader := Uploader{
		Insecure: true,
	}

	rawURL := "://example.com/"
	actualURL, err := uploader.TransformPostBatchURL(rawURL)

	r.Equal("", actualURL)
	r.Equal("parse ://example.com/: missing protocol scheme", err.Error())
}
