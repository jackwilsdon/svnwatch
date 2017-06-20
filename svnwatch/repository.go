package svnwatch

import (
	"encoding/xml"

	"github.com/jackwilsdon/svnwatch/svn"
	"github.com/pkg/errors"
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
		return false, errors.Wrapf(err, "failed to update %s", r.URL)
	}

	if len(info.Entries) == 0 {
		return false, errors.New("no entries in info")
	}

	revision := info.Entries[0].Revision

	if revision > r.Revision {
		r.Revision = revision
		return true, nil
	}

	return false, nil
}
