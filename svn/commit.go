package svn

import (
	"encoding/xml"
	"time"
)

type Commit struct {
	XMLName  xml.Name  `xml:"commit"`
	Revision int       `xml:"revision,attr"`
	Author   string    `xml:"author"`
	Date     time.Time `xml:"date"`
}
