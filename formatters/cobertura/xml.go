package cobertura

import "encoding/xml"

type xmlFile struct {
	XMLName  xml.Name `xml:"coverage"`
	Packages []struct {
		Name    string `xml:"name,attr"`
		Classes []struct {
			Name     string `xml:"name,attr"`
			FileName string `xml:"filename,attr"`
			Lines    []struct {
				Num  int `xml:"number,attr"`
				Hits int `xml:"hits,attr"`
			} `xml:"lines>line"`
		} `xml:"classes>class"`
	} `xml:"packages>package"`
}
