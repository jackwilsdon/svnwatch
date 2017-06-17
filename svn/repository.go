package svn

import "encoding/xml"

type Repository struct {
	XMLName xml.Name `xml:"repository"`
	Root    string   `xml:"root"`
	UUID    string   `xml:"uuid"`
}
