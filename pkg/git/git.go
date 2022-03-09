// Package git provides a small set of helper functions to get information
// about the current repo
package git

import (
	"fmt"
	"github.com/mrtazz/chef-deploy/pkg/command"
	"github.com/mrtazz/chef-deploy/pkg/deploy"
	"strings"
)

type (
	// Ref is a git ref
	Ref string
)

// GetRepoRoot provides a simple
func GetRepoRoot() (string, error) {
	stdout, _, err := command.DefaultRunner{}.Run("git rev-parse --show-toplevel")
	return stdout, err
}

// GetOrigin returns the url of the remote origin
func GetOrigin() (string, error) {
	stdout, _, err := command.DefaultRunner{}.Run("git config --get remote.origin.url")
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

// Differ implements the deploy.Differ interface for git
type Differ struct {
	from, to Ref
	cmd      command.Runner
}

// NewDiffer returns a new git differ
func NewDiffer(from, to Ref) *Differ {
	return &Differ{
		to:   to,
		from: from,
		cmd:  command.DefaultRunner{},
	}
}

// WithRunner returns a differ with the configured runner
func (d *Differ) WithRunner(cmd command.Runner) *Differ {
	return &Differ{
		to:   d.to,
		from: d.from,
		cmd:  cmd,
	}
}

// Diff implements the deploy.Differ interface
func (d *Differ) Diff() ([]deploy.Change, error) {
	ret := make([]deploy.Change, 0, 10)
	stdout, stderr, err := d.cmd.Run(fmt.Sprintf("git diff --name-status %s...%s", d.from, d.to))
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
				mode = deploy.ResourceAdded
			case "D":
				mode = deploy.ResourceDeleted
			case "M":
				mode = deploy.ResourceModified
			}

			ret = append(ret,
				deploy.Change{Type: mode, File: filename})
		}
	}

	return ret, nil
}
