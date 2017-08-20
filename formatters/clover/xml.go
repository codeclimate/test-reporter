package clover

import "encoding/xml"

type xmlFile struct {
	Name  string `xml:"name,attr"`
	Lines []struct {
		Num   int `xml:"num,attr"`
		Count int `xml:"count,attr"`
	} `xml:"line"`
}

type xmlClover struct {
	XMLName  xml.Name `xml:"coverage"`
	Packages []struct {
		Name  string    `xml:"name,attr"`
		Files []xmlFile `xml:"file"`
	} `xml:"project>package"`
	Files []xmlFile `xml:"project>file"`
}
