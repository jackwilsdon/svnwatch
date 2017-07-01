package env

import (
	"fmt"
	"os/exec"

	"github.com/jackwilsdon/svnwatch/svn"
	"github.com/jackwilsdon/svnwatch/svnwatch"
)

func env(cmd *exec.Cmd, repository svnwatch.Repository, revision svn.Revision) error {
	cmd.Env = append(
		cmd.Env,
		fmt.Sprintf("SVN_URL=%s", repository.URL),
		fmt.Sprintf("SVN_REVISION=%d", revision.Revision),
	)
	return nil
}

func init() {
	svnwatch.RegisterCommandType("env", env)
}
