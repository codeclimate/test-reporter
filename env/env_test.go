package env

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Environment_MarshalJSON(t *testing.T) {
	r := require.New(t)
	e, err := New()
	r.NoError(err)

	b, err := e.MarshalJSON()
	r.NoError(err)

	m := map[string]interface{}{}
	err = json.Unmarshal(b, &m)
	r.NoError(err)

	r.NotNil(m["name"])
	r.NotNil(m["committed_at"])
}
