package svnwatch

import (
	"encoding/xml"

	"github.com/pkg/errors"
)

type Watches struct {
	XMLName xml.Name `xml:"watches"`
	Watches []Watch  `xml:"watch"`
}

type Watch struct {
	XMLName  xml.Name  `xml:"watch"`
	URL      string    `xml:"url,attr"`
	Commands []Command `xml:"command"`
}

func (w Watch) Update(repositories *Repositories) error {
	repo := repositories.ForURL(w.URL)

	revisions, err := repo.Update()

	if err != nil {
		return errors.Wrapf(err, "failed to update repository for %s", w.URL)
	}

	for _, revision := range revisions {
		for _, cmd := range w.Commands {
			if err := cmd.Execute(*repo, revision); err != nil {
				return errors.Wrapf(err, "failed to execute command")
			}
		}
	}

	return nil
}

func (w *Watch) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	watch := struct {
		URL      *string   `xml:"url,attr"`
		Commands []Command `xml:"command"`
	}{}

	if err := d.DecodeElement(&watch, &start); err != nil {
		return err
	}

	if watch.URL == nil {
		return errors.New("missing URL from watch")
	}

	w.URL = *watch.URL
	w.Commands = watch.Commands

	return nil
}
