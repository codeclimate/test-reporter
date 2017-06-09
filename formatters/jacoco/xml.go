package jacoco

import "encoding/xml"

type xmlFile struct {
	XMLName  xml.Name `xml:"report"`
	Packages []struct {
		Name       string `xml:"name,attr"`
		SourceFile []struct {
			Name  string `xml:"name,attr"`
			Lines []struct {
				Num  int `xml:"nr,attr"`
				Hits int `xml:"ci,attr"`
			} `xml:"line"`
		} `xml:"sourcefile"`
	} `xml:"package"`
}
