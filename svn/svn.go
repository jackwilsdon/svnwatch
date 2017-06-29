package svn

import (
	"encoding/xml"
	"os/exec"

	"github.com/pkg/errors"
)

func Execute(v interface{}, arg ...string) error {
	output, err := exec.Command("svn", append([]string{"--xml"}, arg...)...).Output()

	if err != nil {
		return errors.Wrap(err, "failed to execute svn")
	}

	return errors.Wrap(xml.Unmarshal(output, &v), "failed to parse output from svn")
}
