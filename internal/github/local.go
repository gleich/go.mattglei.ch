package github

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"go.mattglei.ch/timber"
)

var (
	CLONE_DIRECTORY = "repositories"
	cloneLock       sync.Mutex
)

func SetupCloneFolder() error {
	err := os.MkdirAll(CLONE_DIRECTORY, 0755)
	if err != nil {
		return fmt.Errorf("creating directory %s: %w", CLONE_DIRECTORY, err)
	}
	return nil
}

func (r Repository) Clone() error {
	start := time.Now()

	cloneLock.Lock()
	defer cloneLock.Unlock()

	destination := filepath.Join(CLONE_DIRECTORY, r.Name)
	_, err := os.Stat(destination)
	alreadyCloned := !errors.Is(err, fs.ErrNotExist)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if alreadyCloned {
		out, err := exec.CommandContext(ctx, "git", "-C", destination, "fetch", "origin").
			CombinedOutput()
		if err != nil {
			timber.Debug(string(out))
			return fmt.Errorf("running git fetch: %w", err)
		}
		out, err = exec.CommandContext(ctx, "git", "-C", destination, "reset", "--hard", "origin/HEAD").
			CombinedOutput()
		if err != nil {
			timber.Debug(string(out))
			return fmt.Errorf("running git reset: %w", err)
		}
		timber.DoneSince(start, "updated", r.Name)
	} else {
		repoURL, err := url.JoinPath("https://github.com/gleich", r.Name+".git")
		if err != nil {
			return fmt.Errorf("creating url: %w", err)
		}
		out, err := exec.CommandContext(ctx, "git", "clone", repoURL, destination).
			CombinedOutput()
		if err != nil {
			timber.Debug(string(out))
			return fmt.Errorf("running git clone: %w", err)
		}
		timber.DoneSince(start, "cloned", r.Name)
	}

	return nil
}

func (r Repository) EnsurePath(loc string) bool {
	destination := filepath.Join(CLONE_DIRECTORY, loc)
	_, err := os.Stat(destination)
	return !errors.Is(err, fs.ErrNotExist)
}
