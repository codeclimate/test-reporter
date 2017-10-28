package coveragepy

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Source struct {
	Path string `xml:",chardata"`
}

type xmlFile struct {
	XMLName  xml.Name `xml:"coverage"`
	Sources  []Source `xml:"sources>source"`
	Packages []struct {
		Name    string `xml:"name,attr"`
		Classes []struct {
			FileName string `xml:"filename,attr"`
			Lines    []struct {
				Hits   int `xml:"hits,attr"`
				Number int `xml:"number,attr"`
			} `xml:"lines>line"`
		} `xml:"classes>class"`
	} `xml:"packages>package"`
}

func (covpyFile xmlFile) getFullFilePath(filename string) string {
	fullFilePath := filename

	for _, source := range covpyFile.Sources {
		filepath := fmt.Sprintf("%s/%s", source.Path, filename)
		if _, err := os.Stat(filepath); err == nil {
			fullFilePath = filepath
			break
		}
	}
	return fullFilePath
}
