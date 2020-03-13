package versionize

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Git runs git with the specified arguments, and returns stdout as a string, with whitespace trimmed.
func Git(arg... string) (string, error) {
	c := exec.Command("git", arg...)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		return "", fmt.Errorf("git error: %v",  stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

// MustGitRevision returns the current HEAD revision, suitable for adding to the Revision field of a Metadata struct.
// If Git can't be called or fails, it panics.
func MustGitRevision() string {
	rev, err := GitRevision()
	if err != nil {
		panic(err)
	}
	return rev
}

// GitRevision returns the current HEAD revision, suitable for adding to the Revision field of a Metadata struct.
func GitRevision() (string, error) {
	rev, err := Git("rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("can't get revision info: %w", err)
	}
	return rev, nil
}

// MustGitVersion returns the current code version, suitable for adding to the Version field of a Metadata struct.
// See GitVersion for details. If Git can't be called or fails, it panics.
func MustGitVersion() string {
	rev, err := GitVersion()
	if err != nil {
		panic(err)
	}
	return rev
}

// GitVersion returns the current code version, suitable for adding to the Version field of a Metadata struct.
// It will look like the output of git describe, but with -DEV appended if the build tree is unclean
// (has changes not reflected in Git).
// Examples, in increasing order of dirtiness:
//  - v1.0.2 // clean tagged release code
//  - v1.0.2-3-3efcf4 // tagged release with 3 subsequent commits, last of which is 3efcf4
//  - v1.0.2-3.3efc4-DEV // as above but with uncommitted changes
//  - b834f9 // built from a commit, no releases ever tagged yet
//  - b834f9-DEV // as above, but with uncommitted changes
func GitVersion() (string, error) {
	ver, err := Git("describe", "--always")
	if err != nil {
		return "", fmt.Errorf("can't get code version: %v", err)
	}
	s, err := GitStatus()
	if err != nil {
		return "", fmt.Errorf("can't check code is clean: %v", err)
	}
	if s.Unclean() {
		ver += "-DEV"
	}
	return ver, nil
}

type Status struct {
	Added bool
	Deleted bool
	Modified bool
	Untracked bool
	Renamed bool
	Copied bool
	Updated bool
}

// Unclean returns true if the work tree status indicates an unclean build;
// i.e. if any of the status flags are true, meaning the tree doesn't reflect a Git commit.
func (s *Status) Unclean() bool {
	return s.Added || s.Deleted || s.Modified || s.Untracked || s.Renamed || s.Copied || s.Updated
}

// GitStatus returns the Git status of the work tree, i.e. the status of the files used to build the code,
// as a Status object indicating the presence of files which are Added, Deleted, Modified, Untracked and so on.
func GitStatus() (Status, error) {
	s := Status{}
	stats, err := Git("status", "--porcelain")
	if err != nil {
		return s, err
	}
	lines := strings.Split(stats, "\n")
	for _, l := range lines {
		if len(l) > 3 {
			// x := rune(l[0])
			y := l[1]
			switch y {
			case 'A':
				s.Added = true
			case 'D':
				s.Deleted = true
			case 'M':
				s.Modified = true
			case '?':
				s.Untracked = true
			case 'R':
				s.Renamed = true
			case 'C':
				s.Copied = true
			case 'U':
				s.Updated = true
			}
		}
	}
	return s, nil
}

