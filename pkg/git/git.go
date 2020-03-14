// Package git provides a small set of helper functions to get information
// about the current repo
package git

import (
	"fmt"
	"github.com/mrtazz/chef-deploy/pkg/command"
	"strings"
)

type (
	// Ref is a git ref
	Ref string

	// Diff is a single file change
	Diff struct {
		Mode int
		File string
	}
)

const (
	// DiffAdded represents the status of a file being added in a git diff
	DiffAdded = iota
	// DiffDeleted represents the status of a file being deleted in a diff
	DiffDeleted
	// DiffModified represents the status of a file being modified in a diff
	DiffModified
)

var (
	// CommandRunner is the helper interface to run commands off of
	CommandRunner command.Runner
)

func init() {
	CommandRunner = command.DefaultRunner{}
}

// GetRepoRoot provides a simple
func GetRepoRoot() (string, error) {
	stdout, _, err := CommandRunner.Run("git rev-parse --show-toplevel")
	return stdout, err
}

// GetOrigin returns the url of the remote origin
func GetOrigin() (string, error) {
	stdout, _, err := CommandRunner.Run("git config --get remote.origin.url")
	return stdout, err
}

// GetSlug returns the repo slug
func GetSlug(repoURL string) (string, error) {
	parts := strings.FieldsFunc(repoURL, func(r rune) bool {
		switch r {
		case '/', ':', '.':
			return true
		}
		return false
	})

	var ret string
	// figure out if the URL ends in .git or not and set up the slug accordingly
	if parts[len(parts)-1] == "git" {
		repo := parts[len(parts)-2]
		user := parts[len(parts)-3]
		ret = fmt.Sprintf("%s/%s", user, repo)
	} else {
		repo := parts[len(parts)-1]
		user := parts[len(parts)-2]
		ret = fmt.Sprintf("%s/%s", user, repo)
	}
	return ret, nil
}

// GetSlugForRepo try to get the slug for the current repo
func GetSlugForRepo() (string, error) {
	origin, err := GetOrigin()
	if err != nil {
		return "", err
	}
	return GetSlug(origin)
}

// GetDiff returns a list of files in a diff
func GetDiff(from, to Ref) ([]Diff, error) {
	ret := make([]Diff, 0, 10)
	stdout, stderr, err := CommandRunner.Run(fmt.Sprintf("git diff --name-status %s...%s", from, to))
	if err != nil {
		return ret, fmt.Errorf("%s: %s", err.Error(), stderr)
	}

	lines := strings.Split(stdout, "\n")

	for _, line := range lines {
		lineParts := strings.Fields(line)

		if len(lineParts) == 2 {
			filename := strings.TrimSpace(lineParts[1])
			mode := -1
			switch strings.TrimSpace(lineParts[0]) {
			case "A":
				mode = DiffAdded
			case "D":
				mode = DiffDeleted
			case "M":
				mode = DiffModified
			}

			ret = append(ret,
				Diff{Mode: mode, File: filename})
		}
	}

	return ret, nil
}
