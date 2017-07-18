package svnwatch

import (
	"encoding/xml"

	"github.com/pkg/errors"

	"github.com/jackwilsdon/svnwatch/svn"
)

// Repositories represents a collection of Repository objects.
type Repositories struct {
	XMLName      xml.Name     `xml:"repositories"`
	Repositories []Repository `xml:"repository"`
}

// ForURL returns the repository for the specified URL or creates one if it
// did not already exist.
func (r *Repositories) ForURL(url string) *Repository {
	for i := range r.Repositories {
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

// Repository represents the last known state of a repository.
type Repository struct {
	XMLName  xml.Name `xml:"repository"`
	URL      string   `xml:"url,attr"`
	Revision int      `xml:",chardata"`
	Updated  bool     `xml:"updated,attr"`
}

// Update the repository and return any new revisions.
func (r *Repository) Update() ([]svn.Revision, error) {
	revisions, err := svn.GetLogRange(r.URL, r.Revision, nil)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get log range for %s (range %d:HEAD)", r.URL, r.Revision)
	}

	originalRevision := r.Revision
	originalUpdated := r.Updated

	for _, revision := range revisions {
		if revision.Revision > r.Revision {
			r.Revision = revision.Revision
		}
	}

	r.Updated = true

	// If the revision hasn't changed or we hadn't updated this repository before, ignore the changes
	if r.Revision == originalRevision || !originalUpdated {
		return nil, nil
	}

	if originalRevision > 0 {
		return revisions[1:], nil
	}

	return revisions, nil
}

// UnmarshalXML unmarshals the repository from XML whilst providing some extra
// validation.
func (r *Repository) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	repo := struct {
		URL      *string `xml:"url,attr"`
		Revision *int    `xml:",chardata"`
		Updated  bool    `xml:"updated,attr"`
	}{}

	if err := decoder.DecodeElement(&repo, &start); err != nil {
		return err
	}

	if repo.URL == nil {
		return errors.New("missing URL from repository")
	}

	if repo.Revision == nil {
		return errors.New("missing revision from repository")
	}

	r.URL = *repo.URL
	r.Revision = *repo.Revision
	r.Updated = repo.Updated

	return nil
}
