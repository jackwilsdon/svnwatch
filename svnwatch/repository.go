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

func (r *Repositories) ForURL(url string) *Repository {
	for i, _ := range r.Repositories {
		if url == r.Repositories[i].URL {
			return &r.Repositories[i]
		}
	}

	r.Repositories = append(r.Repositories, Repository{
		Revision: 0,
		URL:      url,
	})

	return &r.Repositories[len(r.Repositories)-1]
}

type Repository struct {
	XMLName  xml.Name `xml:"repository"`
	Revision int      `xml:",chardata"`
	URL      string   `xml:"url,attr"`
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
