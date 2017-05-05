package clover

import "encoding/xml"

type xmlFile struct {
	XMLName  xml.Name `xml:"coverage"`
	Packages []struct {
		Name  string `xml:"name,attr"`
		Files []struct {
			Name  string `xml:"name,attr"`
			Lines []struct {
				Num   int `xml:"num,attr"`
				Count int `xml:"count,attr"`
			} `xml:"line"`
		} `xml:"file"`
	} `xml:"project>package"`
}
