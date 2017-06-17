package svn

import "encoding/xml"

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
		return nil, err
	}

	return &log, nil
}
