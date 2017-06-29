package svnwatch

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/jackwilsdon/svnwatch/svn"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

type Command struct {
	XMLName      xml.Name `xml:"command"`
	ArgumentType string   `xml:"argument-type,attr,omitempty"`
	Command      string   `xml:",chardata"`
}

func (c Command) Execute(repo Repository, revision svn.Revision) error {
	pieces, err := shellwords.Parse(c.Command)

	if err != nil {
		return errors.Wrapf(err, "failed to parse \"%s\"", c.Command)
	}

	cmd := exec.Command(pieces[0], pieces[1:]...)
	cmd.Env = os.Environ()

	if c.ArgumentType == "normal" {
		cmd.Args = append(cmd.Args, repo.URL, strconv.Itoa(repo.Revision))
	} else if c.ArgumentType == "env" {
		cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("SVN_URL=%s", repo.URL),
			fmt.Sprintf("SVN_REVISION=%d", repo.Revision),
		)
	} else {
		return fmt.Errorf("invalid argument type \"%s\" for \"%s\"", c.ArgumentType, c.Command)
	}

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to execute \"%s\"", c.Command)
	}

	return nil
}

func (c *Command) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	cmd := struct {
		ArgumentType *string `xml:"argument-type,attr,omitempty"`
		Command      string  `xml:"chardata"`
	}{}

	if err := d.DecodeElement(&cmd, &start); err != nil {
		return err
	}

	if cmd.ArgumentType == nil {
		c.ArgumentType = "normal"
	} else if *cmd.ArgumentType == "normal" || *cmd.ArgumentType == "env" {
		c.ArgumentType = *cmd.ArgumentType
	} else {
		return fmt.Errorf("invalid argument type \"%s\" for \"%s\"", *cmd.ArgumentType, cmd.Command)
	}

	c.Command = cmd.Command

	return nil
}
