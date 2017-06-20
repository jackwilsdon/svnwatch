package svn

import (
	"encoding/xml"
	"os/exec"

	"github.com/pkg/errors"
)

func Run(v interface{}, name string, arg ...string) error {
	output, err := exec.Command(name, arg...).Output()

	if err != nil {
		return errors.Wrapf(err, "failed to execute %s", name)
	}

	return errors.Wrapf(xml.Unmarshal(output, &v), "failed to parse output from %s", name)
}
