package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type CoverageSummer struct {
	Output string
	Parts  int
}

var summerOptions = CoverageSummer{}

var sumCoverageCmd = &cobra.Command{
	Use:   "sum-coverage",
	Short: "Combine (sum) multiple pre-formatted coverage payloads into one.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// defer profile.Start(profile.MemProfile).Stop()
		if len(args) == 0 {
			return errors.New("you must pass in one or more files to be summarized")
		}
		if summerOptions.Parts != 0 && len(args) != summerOptions.Parts {
			return errors.Errorf("expected %d parts, received %d parts", summerOptions.Parts, len(args))
		}

		rep, err := wip4(args)
		if err != nil {
			return errors.WithStack(err)
		}

		out, err := writer(summerOptions.Output)
		if err != nil {
			return errors.WithStack(err)
		}

		err = rep.Save(out)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	},
}

// 1,16 minutes
func wip4(files []string) (formatters.Report, error) {
	start := time.Now()
	logrus.Debugf("processing %d files", len(files))
	defer func() { logrus.Debugf("finished processing %d files [%s]", len(files), time.Since(start)) }()

	rep := formatters.Report{
		SourceFiles: formatters.SourceFiles{},
	}

	f, err := os.Open(files[0])
	if err != nil {
		return rep, errors.WithStack(err)
	}

	err = json.NewDecoder(f).Decode(&rep)
	if err != nil {
		return rep, errors.WithStack(err)
	}

	files = files[1:]
	reports := make([]*formatters.Report, len(files))
	wg := &errgroup.Group{}
	fx := map[string]bytes.Buffer{}
	for i, n := range files {
		err := func(i int, n string) error {
			st := time.Now()
			logrus.Debugf("processing %s", n)
			defer func() { logrus.Debugf("finished processing %s [%s]", n, time.Since(st)) }()
			f, err := os.Open(n)
			defer f.Close()
			if err != nil {
				return errors.WithStack(err)
			}

			bb := bytes.Buffer{}
			_, err = io.Copy(&bb, f)
			if err != nil {
				return errors.WithStack(err)
			}
			fx[n] = bb

			wg.Go(func() error {
				st := time.Now()
				logrus.Debugf("unmarshaling %s", n)
				defer func() { logrus.Debugf("finished unmarshaling %s [%s]", n, time.Since(st)) }()
				rr := &formatters.Report{
					SourceFiles: formatters.SourceFiles{},
				}

				bb, ok := fx[n]
				if !ok {
					return errors.WithStack(errors.Errorf("could not find json for %s in map", n))
				}
				err = json.Unmarshal(bb.Bytes(), rr)
				if err != nil {
					return errors.WithStack(err)
				}
				reports[i] = rr
				return nil
			})

			return nil
		}(i, n)
		if err != nil {
			return rep, errors.WithStack(err)
		}
	}

	err = wg.Wait()
	if err != nil {
		return rep, errors.WithStack(err)
	}

	st := time.Now()
	logrus.Debugf("merging %d reports", len(reports)+1)
	err = rep.Merge(reports...)
	if err != nil {
		return rep, errors.WithStack(err)
	}
	logrus.Debugf("completed merging %d reports [%s]", len(reports)+1, time.Since(st))
	return rep, nil
}

func init() {
	sumCoverageCmd.Flags().IntVarP(&summerOptions.Parts, "parts", "p", 0, "total number of parts to sum")
	sumCoverageCmd.Flags().StringVarP(&summerOptions.Output, "output", "o", ccDefaultCoveragePath, "output path")
	RootCmd.AddCommand(sumCoverageCmd)
}
