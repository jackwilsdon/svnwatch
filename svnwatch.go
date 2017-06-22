package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func fatalf(format interface{}, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], fmt.Sprintf(fmt.Sprint(format), a...))
	os.Exit(1)
}

func main() {
	configDir := flag.String("config", "~/.svnwatch", "the configuration directory for svnwatch")
	interval := flag.Int("interval", 0, "how often to check for updates (0 disables this and exists after a single check)")

	flag.Parse()

	watcher, err := LoadWatcher(*configDir)

	if *interval < 0 {
		fatalf("%s: invalid interval: %d", os.Args[0], *interval)
	}

	if err != nil {
		fatalf(err)
	}

	for {
		if err := watcher.Update(); err != nil {
			fatalf(err)
		}

		if err := watcher.Save(*configDir); err != nil {
			fatalf(err)
		}

		if *interval > 0 {
			time.Sleep(time.Duration(*interval) * time.Second)
		} else {
			break
		}
	}
}
