package formatters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/codeclimate/test-reporter/env"
)

type SourceFile struct {
	BlobID          string     `json:"blob_id"`
	Coverage        Coverage   `json:"coverage"`
	CoveredPercent  float64    `json:"covered_percent"`
	CoveredStrength int        `json:"covered_strength"`
	LineCounts      LineCounts `json:"line_counts"`
	Name            string     `json:"name"`
}

func NewSourceFile(name string) SourceFile {
	if pwd, err := os.Getwd(); err == nil {
		pwd := fmt.Sprintf("%s%s", pwd, string(os.PathSeparator))
		name = strings.TrimPrefix(name, pwd)
	}

	sf := SourceFile{Name: name}
	sf.BlobID, _ = env.GitSHA(name)
	return sf
}

type Coverage []interface{}

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
	return json.Marshal(strings.TrimSpace(bb.String()))
}
