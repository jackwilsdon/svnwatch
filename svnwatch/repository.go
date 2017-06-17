package svnwatch

import (
	"encoding/xml"
	"errors"

	"github.com/jackwilsdon/svnwatch/svn"
)

type Repositories struct {
	XMLName      xml.Name     `xml:"repositories"`
	Repositories []Repository `xml:"repository"`
}

type Repository struct {
	XMLName  xml.Name `xml:"repository"`
	Revision int      `xml:"revision,attr"`
	URL      string   `xml:",chardata"`
}

func (r *Repository) Update() (bool, error) {
	info, err := svn.GetInfo(r.URL)

	if err != nil {
		return false, err
	}

	if len(info.Entries) == 0 {
		return false, errors.New("no entries in info")
	}

	return info.Entries[0].Revision > r.Revision, nil
}
