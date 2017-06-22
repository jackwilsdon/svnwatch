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
	URL      string   `xml:"url,attr"`
	Revision int      `xml:",chardata"`
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

func (r *Repository) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	repo := struct {
		URL *string `xml:"url,attr"`
	}{}

	if err := d.DecodeElement(&repo, &start); err != nil {
		return err
	}

	if repo.URL == nil {
		return errors.New("missing URL from watch")
	}

	r.URL = *repo.URL

	return nil
}
