package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	configDir := flag.String("config", "~/.svnwatch", "the configuration directory for svnwatch")
	interval := flag.Int("interval", 0, "how often to check for updates (0 disables this and exists after a single check)")

	flag.Parse()

	watcher, err := LoadWatcher(*configDir)

	if *interval < 0 {
		fmt.Fprintf(os.Stderr, "%s: invalid interval: %d", os.Args[0], *interval)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	spew.Dump(watcher.Repositories)

	for {
		if err := watcher.Update(); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		watcher.Save("./config")

		if *interval > 0 {
			time.Sleep(time.Duration(*interval) * time.Second)
		} else {
			break
		}
	}
}
