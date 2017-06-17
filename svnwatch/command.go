package svnwatch

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
)

type Command struct {
	XMLName xml.Name `xml:"command"`
	Type    string   `xml:"type,attr"`
	Command string   `xml:",chardata"`
}

func (c *Command) Execute(repo *Repository) error {
	cmd := exec.Command(c.Command)
	cmd.Env = os.Environ()

	if c.Type == "revision" {
		cmd.Env = append(
			cmd.Env,
			fmt.Sprintf("SVN_REVISION=%d", repo.Revision),
			fmt.Sprintf("SVN_URL=%s", repo.URL),
		)

		if err, ok := cmd.Run().(*exec.ExitError); ok {
			return err
		}
	} else {
		return fmt.Errorf("invalid type: %s", c.Type)
	}

	return nil
}
