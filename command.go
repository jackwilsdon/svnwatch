package svnwatch

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"strings"

	shellwords "github.com/mattn/go-shellwords"
	"github.com/pkg/errors"

	"github.com/jackwilsdon/svnwatch/svn"
)

var commandTypes = make(map[string]CommandType)

// A CommandType is a method of passing data about a revision into a command.
type CommandType func(cmd *exec.Cmd, repository Repository, revision svn.Revision) error

// RegisterCommandType registers the specified command type under the name provided.
func RegisterCommandType(name string, commandType CommandType) {
	if commandType == nil {
		panic(fmt.Sprintf("type \"%s\" is nil", name))
	}

	if _, registered := commandTypes[name]; registered {
		panic(fmt.Sprintf("type \"%s\" is already registered", name))
	}

	commandTypes[name] = commandType
}

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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	commandType, registered := commandTypes[command.Type]

	if !registered {
		keys := make([]string, 0, len(commandTypes))

		for key := range commandTypes {
			keys = append(keys, key)
		}

		return fmt.Errorf("invalid type \"%s\" for \"%s\" (supported types: %s)", command.Type, command.Command, strings.Join(keys, ", "))
	}

	if err := commandType(cmd, repository, revision); err != nil {
		return errors.Wrapf(err, "failed to use command type \"%s\"", command.Type)
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
		return fmt.Errorf("missing type for \"%s\"", cmd.Command)
	}

	if _, registered := commandTypes[*cmd.Type]; !registered {
		keys := make([]string, 0, len(commandTypes))

		for key := range commandTypes {
			keys = append(keys, key)
		}

		return fmt.Errorf("invalid type \"%s\" for \"%s\" (supported types: %s)", *cmd.Type, cmd.Command, strings.Join(keys, ", "))
	}

	command.Type = *cmd.Type
	command.Command = cmd.Command

	return nil
}
