package knife

import (
	"fmt"
	"github.com/fatih/set"
	"github.com/mrtazz/chef-deploy/pkg/command"
	"github.com/mrtazz/chef-deploy/pkg/deploy"
	"path"
	"sort"
	"strings"
)

var (
	defaultKnifeExecutable = "knife"

	knifeCommandLookups = map[string]map[int]string{
		"cookbook": {
			deploy.ResourceAdded:    "%s cookbook upload %s",
			deploy.ResourceDeleted:  "%s cookbook -ay delete %s",
			deploy.ResourceModified: "%s cookbook upload %s",
		},
		"roles": {
			deploy.ResourceAdded:    "%s role from file %s",
			deploy.ResourceDeleted:  "%s role -y delete %s",
			deploy.ResourceModified: "%s role from file %s",
		},
		"data_bags": {
			deploy.ResourceAdded:    "%s data bag from file %s %s",
			deploy.ResourceDeleted:  "%s data bag -y delete %s %s",
			deploy.ResourceModified: "%s data bag from file %s %s",
		},
	}
)

func init() {
}

func (d *Deploy) generateCommands(changes []deploy.Change) (ret []string, err error) {
	knifeCommands := set.New(set.ThreadSafe)

	for _, c := range changes {
		parts := strings.Split(c.File, "/")
		if len(parts) > 0 && parts[0] == "chef" {
			switch parts[1] {
			case "cookbooks":
				if c.Type == deploy.ResourceDeleted {
					if strings.Contains(c.File, "metadata") {
						knifeCommands.Add(fmt.Sprintf(knifeCommandLookups["cookbook"][c.Type],
							d.knifePath, parts[2]))
					}
				} else {
					knifeCommands.Add(fmt.Sprintf(knifeCommandLookups["cookbook"][c.Type],
						d.knifePath, parts[2]))
				}
			case "data_bags":
				bagitem := c.File
				if c.Type == deploy.ResourceDeleted {
					bagitem = path.Base(c.File)
					bagitem = strings.Replace(bagitem, ".json", "", -1)
				}
				knifeCommands.Add(fmt.Sprintf(knifeCommandLookups["data_bags"][c.Type],
					d.knifePath, parts[2], d.stripPrefixMaybe(bagitem)))
			case "roles":
				roleitem := c.File
				if c.Type == deploy.ResourceDeleted {
					roleitem = path.Base(c.File)
					roleitem = strings.Replace(roleitem, ".json", "", -1)
					roleitem = strings.Replace(roleitem, ".rb", "", -1)
				}
				knifeCommands.Add(fmt.Sprintf(knifeCommandLookups["roles"][c.Type],
					d.knifePath, d.stripPrefixMaybe(roleitem)))
			}
		}
	}

	commands := set.StringSlice(knifeCommands)
	sort.Strings(commands)
	return commands, nil

}

// Deploy implements the deployer interface for knife commands
type Deploy struct {
	knifePath    string
	cmd          command.Runner
	subdirectory string
}

// NewDeployer returns a new deployer
func NewDeployer() *Deploy {
	return &Deploy{
		knifePath: defaultKnifeExecutable,
		cmd:       command.NewEchoingRunner(),
	}
}

// WithKnifeExecutable returns a new Deploy struct configured with the given
// knifePath
func (d *Deploy) WithKnifeExecutable(knifePath string) *Deploy {
	newD := d
	newD.knifePath = knifePath
	return newD
}

// WithRunner returns a new Deploy struct configured with the given
// command runner
func (d *Deploy) WithRunner(cmd command.Runner) *Deploy {
	newD := d
	newD.cmd = cmd
	return newD
}

// WithSubdirectory returns a new Deploy struct configured with the given
// command runner
func (d *Deploy) WithSubdirectory(dir string) *Deploy {
	newD := d
	newD.subdirectory = dir
	return newD
}

// PreviewChanges show changes that would be applied
func (d *Deploy) PreviewChanges(changes []deploy.Change) (err error) {
	commands, err := d.generateCommands(changes)
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		fmt.Println(cmd)
	}
	return nil
}

// DeployChanges deploy all change in a local repo from the given sha range
func (d *Deploy) DeployChanges(changes []deploy.Change) (err error) {
	knifeCommandFailed := false

	commands, err := d.generateCommands(changes)
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		_, _, err := d.cmd.Run(cmd)
		if err != nil {
			knifeCommandFailed = true
		}
	}
	if knifeCommandFailed {
		return fmt.Errorf("one or more knife commands failed")
	}
	return nil
}

func (d *Deploy) stripPrefixMaybe(filepath string) string {
	if d.subdirectory == "" {
		return filepath
	}

	if !strings.HasSuffix(d.subdirectory, "/") {
		d.subdirectory = fmt.Sprintf("%s/", d.subdirectory)
	}

	if strings.HasPrefix(filepath, d.subdirectory) {
		return strings.Replace(filepath, d.subdirectory, "", 1)
	}

	return filepath
}
