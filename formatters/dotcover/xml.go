package dotcover

import "encoding/xml"

type xmlDotCover struct {
	XMLName  xml.Name `xml:"Root"`
	Files []struct {
		Path string `xml:"Name,attr"`
		Index int `xml:"Index,attr"`
	} `xml:"FileIndices>File"`
	Statements []struct {
		FileIndex  int `xml:"FileIndex,attr"`
		Covered bool `xml:"Covered,attr"`
	} `xml:"Assembly>Namespace>Type>Method>Statement"`
}
