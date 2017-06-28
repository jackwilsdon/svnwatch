package svn

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type log struct {
	XMLName   xml.Name   `xml:"log"`
	Revisions []Revision `xml:"logentry"`
}

type Revision struct {
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

func GetLog(address string) ([]Revision, error) {
	log := log{}

	if err := Execute(&log, "log", "--xml", "--verbose", address); err != nil {
		return nil, errors.Wrapf(err, "failed to get log for %s", address)
	}

	return log.Revisions, nil
}

func GetLogRange(address string, start int, end *int) ([]Revision, error) {
	log := log{}

	revision := strconv.Itoa(start) + ":"

	if end == nil {
		revision += "HEAD"
	} else {
		revision += strconv.Itoa(*end)
	}

	if err := Execute(&log, "log", "--xml", "--verbose", "--revision", revision, address); err != nil {
		return nil, errors.Wrapf(err, "failed to get log for %s (revision %d)", address, revision)
	}

	return log.Revisions, nil
}

func GetRevision(address string, revision int) (*Revision, error) {
	log := log{}

	if err := Execute(&log, "log", "--xml", "--verbose", "--revision", strconv.Itoa(revision), address); err != nil {
		return nil, errors.Wrapf(err, "failed to get log for %s (revision %d)", address, revision)
	}

	if len(log.Revisions) != 1 {
		return nil, fmt.Errorf("found %d log entries but expected 1 for %s (revison %d)", len(log.Revisions), address, revision)
	}

	return &log.Revisions[0], nil
}

func GetLatestRevision(address string) (*Revision, error) {
	log := log{}

	if err := Execute(&log, "log", "--xml", "--verbose", "--limit", "1", address); err != nil {
		return nil, errors.Wrapf(err, "failed to get log for %s", address)
	}

	if len(log.Revisions) != 1 {
		return nil, fmt.Errorf("found %d log entries but expected 1 for %s", len(log.Revisions), address)
	}

	return &log.Revisions[0], nil
}
