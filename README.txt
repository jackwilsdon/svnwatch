svnwatch
  watcher for SVN repositories

usage:
  -config string
    the configuration directory for svnwatch (default "~/.svnwatch")
  -interval int
    how often to check for updates (0 disables this and exists after a single check)

config:
  watches.xml:
    <watches>
      <watch url="svn://example.com">
        <command type="env">./notify_email</command>
        <command type="env">./notify_slack</command>
      </watch>
    </watches>
