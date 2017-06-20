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
	Message string   `xml:"msg"`
}

func GetLog(address string) (*Log, error) {
	log := Log{}

	if err := Run(&log, "svn", "log", "--xml", address); err != nil {
		return nil, errors.Wrapf(err, "failed to get log for %s", address)
	}

	return &log, nil
}
