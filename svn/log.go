package svn

import (
	"encoding/xml"

	"github.com/pkg/errors"
)

type Log struct {
	XMLName xml.Name   `xml:"log"`
	Entries []LogEntry `xml:"logentry"`
}

type LogEntry struct {
	*Commit
	XMLName xml.Name `xml:"logentry"`
	Paths   []Path   `xml:"paths>path"`
	Message string   `xml:"msg"`
}

type Path struct {
	XMLName               xml.Name `xml:"path"`
	TextModifications     bool     `xml:"text-mods,attr"`
	Kind                  string   `xml:"kind,attr"`
	CopyFromPath          *string  `xml:"copyfrom-path,attr,omitempty"`
	CopyFromRevision      *int     `xml:"copyfrom-rev,attr,omitempty"`
	Action                string   `xml:"action,attr"`
	PropertyModifications bool     `xml:"prop-mods,attr"`
	Name                  string   `xml:",chardata"`
}

func GetLog(address string) (*Log, error) {
	log := Log{}

	if err := Execute(&log, "log", "--xml", "--verbose", address); err != nil {
		return nil, errors.Wrapf(err, "failed to get log for %s", address)
	}

	return &log, nil
}
