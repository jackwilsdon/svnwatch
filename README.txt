svnwatch
  watcher for SVN repositories

installing:
  go get -u github.com/jackwilsdon/svnwatch

usage:
  -config string
    the configuration directory for svnwatch (default "/etc/svnwatch")
  -interval int
    how often to check for updates (0 disables this and exits after a single check) (default 0)

config:
  watches.xml:
    <watches>
      <watch url="svn://example.com">
        <command type="env">./notify_email</command>
        <command type="env">./notify_slack</command>
      </watch>
    </watches>
