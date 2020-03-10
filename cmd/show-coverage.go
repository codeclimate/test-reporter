package cmd

import (
  "fmt"
  "strings"
  "io/ioutil"
  "encoding/json"

  "github.com/pkg/errors"
  "github.com/spf13/cobra"
)

var showCoverageCmd = &cobra.Command{
  Use:   "show-coverage",
  Short: "Show coverage results in standard output",
  RunE: func(cmd *cobra.Command, args []string) error {
    if len(args) == 0 {
      return errors.New("you must pass in one file with the coverage results")
    }

    dat, err := ioutil.ReadFile(args[0])
    if err != nil {
      return errors.New("could not open input file")
    }

    var result map[string]interface{}
    json.Unmarshal([]byte(string(dat)), &result)
    line_counts := result["line_counts"].(map[string]interface{})
    header := "Coverage: %.2f%% (%d/%d lines covered, %d missing)"
    fmt.Println(fmt.Sprintf(header, result["covered_percent"], int(line_counts["covered"].(float64)), int(line_counts["total"].(float64)), int(line_counts["missed"].(float64))))
    if int(line_counts["missed"].(float64)) > 0 {
      fmt.Println("Uncovered lines by file:")
      files := result["source_files"].([]interface{})
      for _, file_obj := range files {
        file := file_obj.(map[string]interface{})
        if file["covered_percent"].(float64) < 100 {
          var uncovered_lines []int
          var values []interface{}
          json.Unmarshal([]byte(file["coverage"].(string)), &values)
          for i, value := range values {
            if value != nil && int(value.(float64)) == 0 {
              uncovered_lines = append(uncovered_lines, i + 1)
            }
          }
          uncovered_lines_str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(uncovered_lines)), ", "), "[]")
          fmt.Println(fmt.Sprintf("%s: %s", file["name"], uncovered_lines_str))
        }
      }
    }

    return nil
  },
}

func init() {
  RootCmd.AddCommand(showCoverageCmd)
}
