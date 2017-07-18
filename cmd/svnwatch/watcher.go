package main

import (
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/jackwilsdon/svnwatch"
)

type watcher struct {
	Repositories svnwatch.Repositories
	Watches      svnwatch.Watches
}

func (w watcher) save(directory string) error {
	if err := save(path.Join(directory, "repositories.xml"), w.Repositories); err != nil {
		return errors.Wrap(err, "failed to save repositories")
	}

	return nil
}

func (w *watcher) update() error {
	for _, watch := range w.Watches.Watches {
		if err := watch.Update(&w.Repositories); err != nil {
			return err
		}
	}

	return nil
}

func loadWatcher(directory string) (*watcher, error) {
	watcher := watcher{}

	repositories := path.Join(directory, "repositories.xml")

	_, err := os.Stat(repositories)

	if err == nil {
		if err := load(repositories, &watcher.Repositories); err != nil {
			return nil, errors.Wrap(err, "failed to load repositories")
		}
	}

	if err := load(path.Join(directory, "watches.xml"), &watcher.Watches); err != nil {
		return nil, errors.Wrap(err, "failed to load watches")
	}

	return &watcher, nil
}
