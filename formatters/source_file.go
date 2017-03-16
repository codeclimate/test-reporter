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

func (a SourceFile) Merge(b SourceFile) SourceFile {
	prime := a
	second := b
	if len(a.Coverage) < len(b.Coverage) {
		prime = b
		second = a
	}

	// go through the shorter or the two:
	for i, c := range second.Coverage {
		if x, ok := c.(int); ok {
			// the secondary is a number
			if y, ok := prime.Coverage[i].(int); ok {
				// the primary is also a number
				prime.Coverage[i] = x + y
				continue
			}
			// set to the secondary:
			prime.Coverage[i] = x
		}
	}
	prime.CalcLineCounts()
	return prime
}

func (sf *SourceFile) CalcLineCounts() {
	lc := LineCounts{
		Total: len(sf.Coverage),
	}

	for _, c := range sf.Coverage {
		if _, ok := c.(int); ok {
			lc.Covered++
			continue
		}
		lc.Missed++
	}

	sf.LineCounts = lc
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
