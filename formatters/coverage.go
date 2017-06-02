package formatters

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

type Coverage []NullInt

// MarshalJSON marshals the coverage into JSON. Since the Code Climate
// API requires this as a string "[1,2,null]" and not just a straight
// JSON array we have to do a bunch of work to coerce into that format
func (c Coverage) MarshalJSON() ([]byte, error) {
	cc := make([]interface{}, 0, len(c))
	for _, x := range c {
		cc = append(cc, x)
	}
	bb := &bytes.Buffer{}
	err := json.NewEncoder(bb).Encode(cc)
	if err != nil {
		return bb.Bytes(), err
	}
	b, err := json.Marshal(strings.TrimSpace(bb.String()))

	if err != nil {
		return b, errors.WithStack(err)
	}

	return b, nil
}

var cwp = []byte("\"")
var cws = []byte("\"")

func (c *Coverage) UnmarshalJSON(text []byte) error {
	text = bytes.TrimPrefix(text, cwp)
	text = bytes.TrimSuffix(text, cws)
	cc := make([]NullInt, 0, 1024)

	err := json.Unmarshal(text, &cc)
	if err != nil {
		return err
	}
	*c = cc
	return nil
}
