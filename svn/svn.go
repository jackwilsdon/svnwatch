package svn

import (
	"encoding/xml"
	"os/exec"

	"github.com/pkg/errors"
)

// Execute svn with the specified arguments and unmarshal it's output into the
// provided interface. --xml is automatically passed to svn as it's first
// argument, meaning it is not required to be passed into this function.
func Execute(v interface{}, arg ...string) error {
	output, err := exec.Command("svn", append([]string{"--xml"}, arg...)...).Output()

	if err != nil {
		return errors.Wrap(err, "failed to execute svn")
	}

	return errors.Wrap(xml.Unmarshal(output, &v), "failed to parse output from svn")
}
