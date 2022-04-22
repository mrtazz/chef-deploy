package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/mrtazz/chef-deploy/pkg/deploy/knife"
	"github.com/mrtazz/chef-deploy/pkg/git"
	"os"
)

var (
	version   = "unknown"
	goversion = "unknown"
	cli       struct {
		Subdirectory    string `help:"subdirectory where chef code is located"`
		KnifeExecutable string `help:"the knife executable to use"`
		Deploy          struct {
			From string `required:"" help:"start ref for generating changes diff"`
			To   string `required:"" help:"end ref for generating changes diff"`
		} `cmd:"" help:"deploy changes."`
		Preview struct {
			From string `required:"" help:"start ref for generating changes diff"`
			To   string `required:"" help:"end ref for generating changes diff"`
		} `cmd:"" help:"preview deploy changes but don't deploy."`
		Version struct {
		} `cmd:"" help:"print version and exit."`
	}
)

func main() {
	ctx := kong.Parse(&cli)
	switch ctx.Command() {
	case "deploy":
		differ := git.NewDiffer(git.Ref(cli.Deploy.From), git.Ref(cli.Deploy.To))
		// it's fine to always try to add the subdirectory here because empty
		// string is the default case anyways
		d := knife.NewDeployer().
			WithSubdirectory(cli.Subdirectory)
		// we don't want to accidentally unset the executable, so only do this if
		// the option was passed
		if cli.KnifeExecutable != "" {
			d = d.WithKnifeExecutable(cli.KnifeExecutable)
		}
		changes, err := differ.Diff()
		if err != nil {
			fmt.Printf("Failed to generate changes from diff: '%s'\n", err.Error())
		}
		if err = d.DeployChanges(changes); err != nil {
			fmt.Printf("Failed to deploy changes from diff: '%s'\n", err.Error())
		}
	case "preview":
		differ := git.NewDiffer(git.Ref(cli.Preview.From), git.Ref(cli.Preview.To))
		// it's fine to always try to add the subdirectory here because empty
		// string is the default case anyways
		d := knife.NewDeployer().
			WithSubdirectory(cli.Subdirectory)
		// we don't want to accidentally unset the executable, so only do this if
		// the option was passed
		if cli.KnifeExecutable != "" {
			d = d.WithKnifeExecutable(cli.KnifeExecutable)
		}
		changes, err := differ.Diff()
		if err != nil {
			fmt.Printf("Failed to generate changes from diff: '%s'\n", err.Error())
		}
		if err = d.PreviewChanges(changes); err != nil {
			fmt.Printf("Failed to preview changes from diff: '%s'\n", err.Error())
		}
	case "version":
		fmt.Printf("chef-deploy %s built for %s\n", version, goversion)
	default:
		fmt.Println("Unknown command: " + ctx.Command())
		os.Exit(1)
	}
}
