package svn

import (
	"encoding/xml"

	"github.com/pkg/errors"
)

type Info struct {
	XMLName xml.Name    `xml:"info"`
	Entries []InfoEntry `xml:"entry"`
}

type InfoEntry struct {
	XMLName     xml.Name   `xml:"entry"`
	Path        string     `xml:"path,attr"`
	Revision    int        `xml:"revision,attr"`
	Kind        string     `xml:"kind,attr"`
	URL         string     `xml:"url"`
	RelativeURL string     `xml:"relative-url"`
	Repository  Repository `xml:"repository"`
	Commit      Commit     `xml:"commit"`
}

func GetInfo(address string) (*Info, error) {
	info := Info{}

	if err := Run(&info, "svn", "info", "--xml", address); err != nil {
		return nil, errors.Wrapf(err, "failed to get info for %s", address)
	}

	return &info, nil
}
