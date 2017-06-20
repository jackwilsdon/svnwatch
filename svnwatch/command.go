package svnwatch

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"

	shellwords "github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

type Command struct {
	XMLName xml.Name `xml:"command"`
	Command string   `xml:",chardata"`
}

func (c Command) Execute(repo Repository) error {
	pieces, err := shellwords.Parse(c.Command)

	if err != nil {
		return errors.Wrapf(err, "failed to parse %s", c.Command)
	}

	cmd := exec.Command(pieces[0], pieces[1:]...)
	cmd.Env = os.Environ()

	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("SVN_URL=%s", repo.URL),
		fmt.Sprintf("SVN_REVISION=%d", repo.Revision),
	)

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to execute %s", c.Command)
	}

	return nil
}
