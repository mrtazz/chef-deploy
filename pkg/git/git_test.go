package git

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var flagtests = []struct {
	in  string
	out string
}{
	{"git@github.com:mrtazz/chef-deploy.git", "mrtazz/chef-deploy"},
	{"https://github.com/mrtazz/chef-deploy.git", "mrtazz/chef-deploy"},
	{"git://github.com/mrtazz/chef-deploy.git", "mrtazz/chef-deploy"},
	{"git@github.com:mrtazz/chef-deploy", "mrtazz/chef-deploy"},
	{"https://github.com/mrtazz/chef-deploy", "mrtazz/chef-deploy"},
	{"git://github.com/mrtazz/chef-deploy", "mrtazz/chef-deploy"},
}

func TestGetConfigKeyFromSection(t *testing.T) {
	for _, tt := range flagtests {
		res, err := GetSlug(tt.in)
		assert.Equal(t, res, tt.out)
		assert.Equal(t, err, nil)
	}
}

type fakeRunner struct{}

func (f fakeRunner) Run(cmd string) (string, string, error) {
	return `M command/command.go
  M command/command_test.go
  M git/git.go
  M infra/agent/checks/freebsdkernel/freebsdkernel.go
  M infra/agent/checks/freebsdversion/freebsdversion.go
  M infra/agent/sysinfo/sysinfo.go
  M infra/agent/sysinfo/sysinfo_test.go
  M mail/mail.go
  A foo/bla.go
  D bla/foo.go`, "", nil
}

func TestGetDiff(t *testing.T) {
	CommandRunner = fakeRunner{}

	diff, err := GetDiff("HEAD~5", "HEAD")

	assert.Equal(t, nil, err)
	assert.Equal(t, 10, len(diff))
	assert.Equal(t, "command/command.go", diff[0].File)
	assert.Equal(t, DiffModified, diff[0].Mode)
	assert.Equal(t, "foo/bla.go", diff[8].File)
	assert.Equal(t, DiffAdded, diff[8].Mode)
	assert.Equal(t, "bla/foo.go", diff[9].File)
	assert.Equal(t, DiffDeleted, diff[9].Mode)

}
