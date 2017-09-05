package cobertura

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Lines struct {
	Num  int `xml:"number,attr"`
	Hits int `xml:"hits,attr"`
}

type Source struct {
	Path string `xml:",chardata"`
}

type xmlClass struct {
	Name     string  `xml:"name,attr"`
	FileName string  `xml:"filename,attr"`
	Lines    []Lines `xml:"lines>line"`
}

type xmlFile struct {
	XMLName  xml.Name `xml:"coverage"`
	Sources  []Source `xml:"sources>source"`
	Packages []struct {
		Name    string     `xml:"name,attr"`
		Classes []xmlClass `xml:"classes>class"`
	} `xml:"packages>package"`
}

// Interface to sort []Lines by line number
type ByLineNum []Lines

func (a ByLineNum) Len() int           { return len(a) }
func (a ByLineNum) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLineNum) Less(i, j int) bool { return a[i].Num < a[j].Num }

func (coberturaFile xmlFile) getFullFilePath(filename string) string {
	fullFilePath := filename

	for _, source := range coberturaFile.Sources {
		filepath := fmt.Sprintf("%s/%s", source.Path, filename)
		if _, err := os.Stat(filepath); err == nil {
			fullFilePath = filepath
			break
		}
	}
	return fullFilePath
}
