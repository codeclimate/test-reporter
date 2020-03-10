package cmd

import (
  "fmt"
  "strings"
  "io/ioutil"
  "encoding/json"

  "github.com/pkg/errors"
  "github.com/spf13/cobra"
)

func getLineCount(result map[string]interface{}, key string) int {
  line_counts := result["line_counts"].(map[string]interface{})
  return int(line_counts[key].(float64))
}

func printHeader(result map[string]interface{}) {
  header := "Coverage: %.2f%% (%d/%d lines covered, %d missing)"
  fmt.Println(fmt.Sprintf(header, result["covered_percent"], getLineCount(result, "covered"), getLineCount(result, "total"), getLineCount(result, "missed")))
}

func printUncoveredLines(result map[string]interface{}) {
  if getLineCount(result, "missed") > 0 {

    fmt.Println("Uncovered lines by file:")
    files := result["source_files"].([]interface{})

    for _, file_obj := range files {
      file := file_obj.(map[string]interface{})
      if file["covered_percent"].(float64) < 100 {
        printUncoveredLinesFromFile(file)
      }
    }
  }
}

func printUncoveredLinesFromFile(file map[string]interface{}) {
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

    // Parse JSON from coverage file
    var result map[string]interface{}
    json.Unmarshal([]byte(string(dat)), &result)

    printHeader(result)

    // If there are missed lines, print which are them, by file
    printUncoveredLines(result)

    return nil
  },
}

func init() {
  RootCmd.AddCommand(showCoverageCmd)
}
