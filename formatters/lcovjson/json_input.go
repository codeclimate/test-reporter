package lcovjson

import (
	"encoding/json"
	"errors"
)

type segment struct {
	Line int
	Column int
	Count int
	HasCount bool
	IsRegionEntry bool
}

func (segment *segment) UnmarshalJSON(data []byte) error {
	var array []interface{}
	if err := json.Unmarshal(data, &array); err != nil {
		return err
	}

	if n, ok := array[0].(float64); ok {
		segment.Line = int(n)
	} else {
		return errors.New("invalid Line")
	}

	if n, ok := array[1].(float64); ok {
		segment.Column = int(n)
	} else {
		return errors.New("invalid Column")
	}

	if n, ok := array[2].(float64); ok {
		segment.Count = int(n)
	} else {
		return errors.New("invalid Count")
	}

	if b, ok := array[3].(bool); ok {
		segment.HasCount = b
	} else {
		return errors.New("invalid HasCount")
	}

	if b, ok := array[4].(bool); ok {
		segment.IsRegionEntry = b
	} else {
		return errors.New("invalid IsRegionEntry")
	}

	return nil
}

type coverage struct {
	Count int `json:"count"`
	Covered int `json:"covered"`
	Percent float64 `json:"percent"`
}

type summary struct {
	Functions coverage `json:"functions"`
	Instantiations coverage `json:"instantiations"`
	Lines coverage `json:"lines"`
	Regions coverage `json:"regions"`
}

type sourceFile struct {
	Filename string `json:"filename"`
	Segments []segment `json:"segments"`
	Summary summary `json:"summary"`
}

type region struct {
	LineStart int
	ColumnStart int
	LineEnd int
	ColumnEnd int
	ExecutionCount int
	FileID int
	ExpandedFileID int
	Kind int
}

func (region *region) UnmarshalJSON(data []byte) error {
	var array []interface{}
	if err := json.Unmarshal(data, &array); err != nil {
		return err
	}

	if n, ok := array[0].(float64); ok {
		region.LineStart = int(n)
	} else {
		return errors.New("invalid LineStart")
	}

	if n, ok := array[1].(float64); ok {
		region.ColumnStart = int(n)
	} else {
		return errors.New("invalid ColumnStart")
	}

	if n, ok := array[2].(float64); ok {
		region.LineEnd = int(n)
	} else {
		return errors.New("invalid LineEnd")
	}

	if n, ok := array[3].(float64); ok {
		region.ColumnEnd = int(n)
	} else {
		return errors.New("invalid ColumnEnd")
	}

	if n, ok := array[4].(float64); ok {
		region.ExecutionCount = int(n)
	} else {
		return errors.New("invalid ExecutionCount")
	}

	if n, ok := array[5].(float64); ok {
		region.FileID = int(n)
	} else {
		return errors.New("invalid FileID")
	}

	if n, ok := array[6].(float64); ok {
		region.ExpandedFileID = int(n)
	} else {
		return errors.New("invalid ExpandedFileID")
	}

	if n, ok := array[7].(float64); ok {
		region.Kind = int(n)
	} else {
		return errors.New("invalid Kind")
	}

	return nil
}

type function struct {
	Count int `json:"count"`
	Filenames []string `json:"filenames"`
	Name string `json:"name"`
	Regions []region `json:"regions"`
}

type lcovJsonFile struct {
	Data []struct {
		Files []sourceFile `json:"files"`
		Functions []function `json:"functions"`
		Totals summary `json:"totals"`
	} `json:"data"`

	Type string `json:"type"`
	Version string `json:"version"`
}
