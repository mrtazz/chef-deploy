package knife

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
	M chef/cookbooks/chef-deploy/attributes/default.rb
	D chef/cookbooks/deleted/metadata.rb
	M chef/cookbooks/homedirs/recipes/mrtazz.rb
	M chef/cookbooks/homedirs/attributes/mrtazz.rb
	M chef/cookbooks/homedirs/recipes/mrtazz.rb
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

func TestSomething(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]struct {
		subdirectory     string
		knifeExecutable  string
		expectedCommands []string
	}{
		"DeployChanges": {
			expectedCommands: []string{
				"knife cookbook -ay delete deleted",
				"knife cookbook upload chef-deploy",
				"knife cookbook upload homedirs",
				"knife cookbook upload jenkins",
				"knife cookbook upload nginx",
				"knife data bag -y delete chef-keys bla",
				"knife data bag -y delete chef-keys blubb",
				"knife data bag from file chef-keys chef/data_bags/chef-keys/friday.json",
				"knife data bag from file chef-keys chef/data_bags/chef-keys/friday_keys.json",
				"knife role -y delete jail",
				"knife role from file chef/roles/ci.rb",
			},
		},
		"DeployChangesWithSubdirectory": {
			subdirectory: "chef",
			expectedCommands: []string{
				"knife cookbook -ay delete deleted",
				"knife cookbook upload chef-deploy",
				"knife cookbook upload homedirs",
				"knife cookbook upload jenkins",
				"knife cookbook upload nginx",
				"knife data bag -y delete chef-keys bla",
				"knife data bag -y delete chef-keys blubb",
				"knife data bag from file chef-keys data_bags/chef-keys/friday.json",
				"knife data bag from file chef-keys data_bags/chef-keys/friday_keys.json",
				"knife role -y delete jail",
				"knife role from file roles/ci.rb",
			},
		},
		"DeployChangesWithSubdirectoryWithSlash": {
			subdirectory: "chef/",
			expectedCommands: []string{
				"knife cookbook -ay delete deleted",
				"knife cookbook upload chef-deploy",
				"knife cookbook upload homedirs",
				"knife cookbook upload jenkins",
				"knife cookbook upload nginx",
				"knife data bag -y delete chef-keys bla",
				"knife data bag -y delete chef-keys blubb",
				"knife data bag from file chef-keys data_bags/chef-keys/friday.json",
				"knife data bag from file chef-keys data_bags/chef-keys/friday_keys.json",
				"knife role -y delete jail",
				"knife role from file roles/ci.rb",
			},
		},
		"DeployChangesWithDifferentKnifeExecutable": {
			knifeExecutable: "/opt/chef/bin/knife",
			expectedCommands: []string{
				"/opt/chef/bin/knife cookbook -ay delete deleted",
				"/opt/chef/bin/knife cookbook upload chef-deploy",
				"/opt/chef/bin/knife cookbook upload homedirs",
				"/opt/chef/bin/knife cookbook upload jenkins",
				"/opt/chef/bin/knife cookbook upload nginx",
				"/opt/chef/bin/knife data bag -y delete chef-keys bla",
				"/opt/chef/bin/knife data bag -y delete chef-keys blubb",
				"/opt/chef/bin/knife data bag from file chef-keys chef/data_bags/chef-keys/friday.json",
				"/opt/chef/bin/knife data bag from file chef-keys chef/data_bags/chef-keys/friday_keys.json",
				"/opt/chef/bin/knife role -y delete jail",
				"/opt/chef/bin/knife role from file chef/roles/ci.rb",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// execute test logic here
			differ := git.NewDiffer("HEAD~5", "HEAD").WithRunner(fakeGitRunner{})
			RanCommands = make([]string, 0, 10)
			d := NewDeployer().WithRunner(fakeTestRunner)
			if tc.subdirectory != "" {
				d = d.WithSubdirectory(tc.subdirectory)
			}
			if tc.knifeExecutable != "" {
				d = d.WithKnifeExecutable(tc.knifeExecutable)
			}
			changes, err := differ.Diff()
			assert.Nil(err)
			d.DeployChanges(changes)
			sort.Strings(RanCommands)

			assert.Equal(len(tc.expectedCommands), len(RanCommands))
			assert.Equal(tc.expectedCommands, RanCommands)
		})
	}
}
