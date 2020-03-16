package chef

import (
	"fmt"
	"github.com/fatih/set"
	"github.com/mrtazz/chef-deploy/pkg/command"
	"github.com/mrtazz/chef-deploy/pkg/git"
	"path"
	"sort"
	"strings"
)

var (
	// CommandRunner is the interface to run commands
	CommandRunner command.Runner
	// KnifeExecutable holds the full path to the knife command
	KnifeExecutable     = "knife"
	knifeCommandLookups = map[string]map[int]string{
		"cookbook": {
			git.DiffAdded:    "%s cookbook upload %s",
			git.DiffDeleted:  "%s cookbook -y delete %s",
			git.DiffModified: "%s cookbook upload %s",
		},
		"roles": {
			git.DiffAdded:    "%s role from file %s",
			git.DiffDeleted:  "%s role -y delete %s",
			git.DiffModified: "%s role from file %s",
		},
		"data_bags": {
			git.DiffAdded:    "%s data bag from file %s %s",
			git.DiffDeleted:  "%s data bag -y delete %s %s",
			git.DiffModified: "%s data bag from file %s %s",
		},
	}
)

func init() {
	CommandRunner = command.NewEchoingRunner()
}

func generateCommands(from, to git.Ref) (ret []string, err error) {
	knifeCommands := set.New(set.ThreadSafe)

	changes, err := git.GetDiff(from, to)
	if err != nil {
		return ret, err
	}

	for _, c := range changes {
		parts := strings.Split(c.File, "/")
		if len(parts) > 0 && parts[0] == "chef" {
			switch parts[1] {
			case "cookbooks":
				if c.Mode == git.DiffDeleted {
					if strings.Contains(c.File, "metadata") {
						knifeCommands.Add(fmt.Sprintf(knifeCommandLookups["cookbook"][c.Mode],
							KnifeExecutable, parts[2]))
					}
				} else {
					knifeCommands.Add(fmt.Sprintf(knifeCommandLookups["cookbook"][c.Mode],
						KnifeExecutable, parts[2]))
				}
			case "data_bags":
				bagitem := c.File
				if c.Mode == git.DiffDeleted {
					bagitem = path.Base(c.File)
					bagitem = strings.Replace(bagitem, ".json", "", -1)
				}
				knifeCommands.Add(fmt.Sprintf(knifeCommandLookups["data_bags"][c.Mode],
					KnifeExecutable, parts[2], bagitem))
			case "roles":
				roleitem := c.File
				if c.Mode == git.DiffDeleted {
					roleitem = path.Base(c.File)
					roleitem = strings.Replace(roleitem, ".json", "", -1)
					roleitem = strings.Replace(roleitem, ".rb", "", -1)
				}
				knifeCommands.Add(fmt.Sprintf(knifeCommandLookups["roles"][c.Mode],
					KnifeExecutable, roleitem))
			}
		}
	}

	commands := set.StringSlice(knifeCommands)
	sort.Strings(commands)
	return commands, nil

}

// PreviewChanges show changes that would be applied
func PreviewChanges(from, to git.Ref) (err error) {
	commands, err := generateCommands(from, to)
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		fmt.Println(cmd)
	}
	return nil
}

// DeployChanges deploy all change in a local repo from the given sha range
func DeployChanges(from, to git.Ref) (err error) {
	knifeCommandFailed := false

	commands, err := generateCommands(from, to)
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		_, _, err := CommandRunner.Run(cmd)
		if err != nil {
			knifeCommandFailed = true
		}
	}
	if knifeCommandFailed {
		return fmt.Errorf("One or more knife commands failed.")
	}
	return nil
}
