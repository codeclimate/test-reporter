package formatters

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
)

type SourceFile struct {
	BlobID          string     `json:"blob_id"`
	Coverage        Coverage   `json:"coverage"`
	CoveredPercent  float64    `json:"covered_percent"`
	CoveredStrength float64    `json:"covered_strength"`
	LineCounts      LineCounts `json:"line_counts"`
	Name            string     `json:"name"`
}

func (a SourceFile) Merge(b SourceFile) (SourceFile, error) {
	if a.BlobID != b.BlobID {
		return a, errors.Errorf("failed to merge coverage for source file %s: BlobID mismatch", a.Name)
	}
	lenA := len(a.Coverage)
	lenB := len(b.Coverage)

	if lenA > lenB {
		b.Coverage = b.Coverage.AppendNulls(lenA - lenB)
	} else if lenA < lenB {
		a.Coverage = a.Coverage.AppendNulls(lenB - lenA)
	}

	for i, bc := range b.Coverage {
		ac := a.Coverage[i]
		if ac.Valid && bc.Valid {
			// they're both valid numbers so add them:
			a.Coverage[i] = NewNullInt(ac.Int + bc.Int)
			continue
		}

		// ac is null and bc isn't so ac takes precendence, we do nothing
		if !ac.Valid {
			continue
		}

		// bc is null and ac isn't so bc takes precendence
		if !bc.Valid {
			a.Coverage[i] = bc
		}

	}
	a.CalcLineCounts()
	return a, nil
}

func (sf *SourceFile) CalcLineCounts() {
	lc := LineCounts{}

	for _, c := range sf.Coverage {
		if !c.Valid {
			continue
		}
		lc.Total++
		lc.Strength += c.Int
		if c.Int == 0 {
			lc.Missed++
			continue
		}
		lc.Covered++
	}
	sf.LineCounts = lc
	sf.CoveredPercent = lc.CoveredPercent()
	sf.CoveredStrength = lc.CoveredStrength()
}

func NewSourceFile(name string, commit *object.Commit) (SourceFile, error) {
	if prefix, err := envy.MustGet("PREFIX"); err == nil {
		if strings.HasSuffix(prefix, string(os.PathSeparator)) {
			name = strings.TrimPrefix(name, prefix)
			logrus.Printf("triming with prefix %s\n", prefix)
		} else {
			prefix = fmt.Sprintf("%s%s", prefix, string(os.PathSeparator))
			logrus.Printf("triming with prefix %s\n", prefix)
			name = strings.TrimPrefix(name, prefix)
		}
	}

	sf := SourceFile{
		Name:     name,
		Coverage: Coverage{},
	}

	var err error
	sf.BlobID, err = env.GitBlob(name, commit)

	if err != nil {
		return sf, errors.WithStack(err)
	}

	if addPrefix, err := envy.MustGet("ADD_PREFIX"); err == nil {
		if addPrefix != "" {
			if strings.HasSuffix(addPrefix, string(os.PathSeparator)) {
				sf.Name = addPrefix + sf.Name
			} else {
				sf.Name = addPrefix + string(os.PathSeparator) + sf.Name
			}
		}
	}

	return sf, nil
}

type SourceFiles map[string]SourceFile

func (sf SourceFiles) MarshalJSON() ([]byte, error) {
	files := []SourceFile{}
	for _, s := range sf {
		s.CalcLineCounts()
		files = append(files, s)
	}
	b, err := json.Marshal(files)
	if err != nil {
		logrus.Errorf("error marshalling source files: %+v\n", sf)
	}
	return b, err
}

func (sf SourceFiles) UnmarshalJSON(text []byte) error {
	files := []SourceFile{}
	err := json.Unmarshal(text, &files)
	if err != nil {
		return err
	}
	for _, f := range files {
		if ff, ok := sf[f.Name]; ok {
			f, err = ff.Merge(f)
			if err != nil {
				return err
			}
		}
		sf[f.Name] = f
	}
	return nil
}
