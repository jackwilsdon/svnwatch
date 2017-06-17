package svn

import (
	"encoding/xml"
	"os/exec"
)

func Run(v interface{}, name string, arg ...string) error {
	output, err := exec.Command(name, arg...).Output()

	if err != nil {
		return err
	}

	if err := xml.Unmarshal(output, &v); err != nil {
		return err
	}

	return nil
}
