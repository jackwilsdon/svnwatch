package svn

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type log struct {
	XMLName   xml.Name   `xml:"log"`
	Revisions []Revision `xml:"logentry"`
}

// A Revision is a revision in SVN.
type Revision struct {
	Revision int       `xml:"revision,attr"`
	Author   string    `xml:"author"`
	Date     time.Time `xml:"date"`
	XMLName  xml.Name  `xml:"logentry"`
	Paths    []Path    `xml:"paths>path"`
	Message  string    `xml:"msg"`
}

// A Path is a record of a change made to a file in a revision.
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

// GetLog returns all revisions for the specified address.
func GetLog(address string) ([]Revision, error) {
	log := log{}

	if err := Execute(&log, "log", "--verbose", address); err != nil {
		return nil, errors.Wrapf(err, "failed to get log for %s", address)
	}

	return log.Revisions, nil
}

// GetLogRange returns all revisions within the specified range. If the end
// revision is nil then all revisions after the specified start revision
// will be returned (including the start revision itself).
func GetLogRange(address string, start int, end *int) ([]Revision, error) {
	log := log{}

	revision := strconv.Itoa(start) + ":"

	if end == nil {
		revision += "HEAD"
	} else {
		revision += strconv.Itoa(*end)
	}

	if err := Execute(&log, "log", "--verbose", "--revision", revision, address); err != nil {
		return nil, errors.Wrapf(err, "failed to get log for %s (revision %s)", address, revision)
	}

	return log.Revisions, nil
}

// GetRevision returns the specified revision for the specified address.
func GetRevision(address string, revision int) (*Revision, error) {
	log := log{}

	if err := Execute(&log, "log", "--verbose", "--revision", strconv.Itoa(revision), address); err != nil {
		return nil, errors.Wrapf(err, "failed to get log for %s (revision %d)", address, revision)
	}

	if len(log.Revisions) != 1 {
		return nil, fmt.Errorf("found %d log entries but expected 1 for %s (revison %d)", len(log.Revisions), address, revision)
	}

	return &log.Revisions[0], nil
}

// GetLatestRevision returns the latest revision for the specified address.
func GetLatestRevision(address string) (*Revision, error) {
	log := log{}

	if err := Execute(&log, "log", "--verbose", "--limit", "1", address); err != nil {
		return nil, errors.Wrapf(err, "failed to get log for %s", address)
	}

	if len(log.Revisions) != 1 {
		return nil, fmt.Errorf("found %d log entries but expected 1 for %s", len(log.Revisions), address)
	}

	return &log.Revisions[0], nil
}
