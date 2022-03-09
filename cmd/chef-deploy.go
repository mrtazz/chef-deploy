package main

import (
//"github.com/docopt/docopt-go"
//"github.com/mrtazz/chef-deploy/pkg/git"
//"github.com/mrtazz/chef-deploy/pkg/version"
//"log"
//"os"
)

var (
	usage = `chef-deploy.

  Usage:
  chef-deploy deploy --from=<from> --to=<to> [options]
  chef-deploy preview --from=<from> --to=<to> [options]
  chef-deploy -h | --help
  chef-deploy --version

  Options:
  --from=<from>                 base SHA of the diff to deploy
  --to=<to>                     head SHA of the diff to deploy
  --knife-executable=<knife>    the knife executable to use
  --subdirectory=<subdirectory> the subdirectory of the repo that contains chef code
  -h --help                     Show this screen.
  --version                     Show version.
`

	isDebug = false
)

func main() {
	//	args, err := docopt.Parse(usage, nil, true,
	//		version.GetDocoptVersionString("chef-deploy"), false)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	if args["--knife-executable"] != nil {
	//		chef.KnifeExecutable = args["--knife-executable"].(string)
	//	}
	//
	//	if args["--subdirectory"] != nil {
	//		chef.Subdirectory = args["--subdirectory"].(string)
	//	}
	//
	//	if args["preview"].(bool) {
	//		err = chef.PreviewChanges(git.Ref(args["--from"].(string)),
	//			git.Ref(args["--to"].(string)))
	//
	//		if err != nil {
	//			log.Println(err.Error())
	//			os.Exit(1)
	//		}
	//	}
	//
	//	if args["deploy"].(bool) {
	//		err = chef.DeployChanges(git.Ref(args["--from"].(string)),
	//			git.Ref(args["--to"].(string)))
	//
	//		if err != nil {
	//			log.Println(err.Error())
	//			os.Exit(1)
	//		}
	//	}

}
