package chef

import (
	"github.com/mrtazz/chef-deploy/pkg/git"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

type fakeRunner struct{}

func (f fakeRunner) Run(cmd string) (string, string, error) {
	RanCommands = append(RanCommands, cmd)
	return "", "", nil
}

type fakeGitRunner struct{}

func (f fakeGitRunner) Run(cmd string) (string, string, error) {
	return `M Jenkinsfile
	M chef/Jenkinsfile
	M chef/cookbooks/homedirs/recipes/mrtazz.rb
	M chef/cookbooks/homedirs/attributes/mrtazz.rb
	M chef/cookbooks/homedirs/recipes/mrtazz.rb
	M chef/cookbooks/jarvis/attributes/default.rb
	M chef/cookbooks/jenkins/recipes/chef_keys.rb
	M chef/cookbooks/nginx/recipes/pkgng.rb
	M chef/data_bags/chef-keys/friday.json
	M chef/data_bags/chef-keys/friday_keys.json
	M chef/roles/ci.rb
	D chef/data_bags/chef-keys/bla.json
	D chef/data_bags/chef-keys/blubb.json
	D chef/roles/jail.rb`, "", nil
}

var (
	fakeTestRunner = fakeRunner{}
	RanCommands    []string
)

func TestDeployChanges(t *testing.T) {
	git.CommandRunner = fakeGitRunner{}
	RanCommands = make([]string, 0, 10)
	CommandRunner = fakeTestRunner

	DeployChanges("HEAD~5", "HEAD")

	sort.Strings(RanCommands)

	assert.Equal(t, len(RanCommands), 10)
	expectedRanCommands := []string{
		"knife cookbook upload homedirs",
		"knife cookbook upload jarvis",
		"knife cookbook upload jenkins",
		"knife cookbook upload nginx",
		"knife data bag -y delete chef-keys bla",
		"knife data bag -y delete chef-keys blubb",
		"knife data bag from file chef-keys chef/data_bags/chef-keys/friday.json",
		"knife data bag from file chef-keys chef/data_bags/chef-keys/friday_keys.json",
		"knife role -y delete jail",
		"knife role from file chef/roles/ci.rb",
	}

	for i := range RanCommands {
		assert.Equal(t, expectedRanCommands[i], RanCommands[i])
	}

}

func TestDeployChangesWithDifferentKnifeExecutable(t *testing.T) {
	git.CommandRunner = fakeGitRunner{}
	RanCommands = make([]string, 0, 10)
	CommandRunner = fakeTestRunner

	KnifeExecutable = "/opt/chef/bin/knife"

	DeployChanges("HEAD~5", "HEAD")

	sort.Strings(RanCommands)

	assert.Equal(t, len(RanCommands), 10)

	expectedRanCommands := []string{
		"/opt/chef/bin/knife cookbook upload homedirs",
		"/opt/chef/bin/knife cookbook upload jarvis",
		"/opt/chef/bin/knife cookbook upload jenkins",
		"/opt/chef/bin/knife cookbook upload nginx",
		"/opt/chef/bin/knife data bag -y delete chef-keys bla",
		"/opt/chef/bin/knife data bag -y delete chef-keys blubb",
		"/opt/chef/bin/knife data bag from file chef-keys chef/data_bags/chef-keys/friday.json",
		"/opt/chef/bin/knife data bag from file chef-keys chef/data_bags/chef-keys/friday_keys.json",
		"/opt/chef/bin/knife role -y delete jail",
		"/opt/chef/bin/knife role from file chef/roles/ci.rb",
	}

	for i := range RanCommands {
		assert.Equal(t, expectedRanCommands[i], RanCommands[i])
	}

}
