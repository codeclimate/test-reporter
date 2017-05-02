package coveragepy

import "encoding/xml"

type xmlFile struct {
	XMLName  xml.Name `xml:"coverage"`
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
