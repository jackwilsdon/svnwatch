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
      <watch url="http://svn.example.com">
        <command>./send_email</command>
        <command argument-type="env">./send_slack</command>
      </watch>
    </watches>