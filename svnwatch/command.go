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

// Command represents a command that is executed when a change is detected.
type Command struct {
	XMLName xml.Name `xml:"command"`
	Type    string   `xml:"type,attr,omitempty"`
	Command string   `xml:",chardata"`
}

// Execute the command on the specified repository and revision.
func (command Command) Execute(repository Repository, revision svn.Revision) error {
	pieces, err := shellwords.Parse(command.Command)

	if err != nil {
		return errors.Wrapf(err, "failed to parse \"%s\"", command.Command)
	}

	if len(pieces) == 0 {
		return fmt.Errorf("failed to parse \"%s\"", command.Command)
	}

	cmd := exec.Command(pieces[0], pieces[1:]...)
	cmd.Env = os.Environ()

	if command.Type == "normal" {
		cmd.Args = append(cmd.Args, repository.URL, strconv.Itoa(revision.Revision))
	} else if command.Type == "env" {
		cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("SVN_URL=%s", repository.URL),
			fmt.Sprintf("SVN_REVISION=%d", revision.Revision),
		)
	} else {
		return fmt.Errorf("invalid type \"%s\" for \"%s\"", command.Type, command.Command)
	}

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to execute \"%s\"", command.Command)
	}

	return nil
}

// UnmarshalXML unmarshals the command from XML whilst providing some extra
// validation.
func (command *Command) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	cmd := struct {
		Type    *string `xml:"type,attr,omitempty"`
		Command string  `xml:",chardata"`
	}{}

	if err := decoder.DecodeElement(&cmd, &start); err != nil {
		return err
	}

	if cmd.Type == nil {
		command.Type = "normal"
	} else if *cmd.Type == "normal" || *cmd.Type == "env" {
		command.Type = *cmd.Type
	} else {
		return fmt.Errorf("invalid type \"%s\" for \"%s\"", *cmd.Type, cmd.Command)
	}

	command.Type = *cmd.Type
	command.Command = cmd.Command

	return nil
}
