package simplecov

import (
	"encoding/json"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

type resultSet struct {
	Coverage map[string]formatters.Coverage `json:"coverage"`
}

func legacyFormat(r Formatter, rep formatters.Report) (formatters.Report, error) {
        logrus.Debugf("Analyzing simplecov json output from legacy format %s", r.Path)
        jf, err := os.Open(r.Path)
        if err != nil {
          return rep, errors.WithStack(errors.Errorf("could not open coverage file %s", r.Path))
        }

        m := map[string]resultSet{}
        err = json.NewDecoder(jf).Decode(&m)

        if err != nil {
          return rep, errors.WithStack(err)
        }

        gitHead, _ := env.GetHead()
        for _, v := range m {
          for n, ls := range v.Coverage {
                  fe, err := formatters.NewSourceFile(n, gitHead)
                  if err != nil {
                          return rep, errors.WithStack(err)
                  }
                  fe.Coverage = ls
                  err = rep.AddSourceFile(fe)
                  if err != nil {
                          return rep, errors.WithStack(err)
                  }
          }
        }

        return rep, nil
}
