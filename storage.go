package main

import (
	"encoding/xml"
	"io/ioutil"

	"github.com/pkg/errors"
)

func load(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	return errors.Wrapf(xml.Unmarshal(data, v), "failed to parse %s", filename)
}

func save(filename string, v interface{}) error {
	data, err := xml.MarshalIndent(v, "", "    ")

	if err != nil {
		return errors.Wrapf(err, "failed to convert %s", filename)
	}

	return errors.Wrap(ioutil.WriteFile(filename, data, 0666), "failed to write file")
}
