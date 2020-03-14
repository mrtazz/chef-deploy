# chef-deploy

A tool to deploy to chef based on a git diff


## Usage

```
chef-deploy.

  Usage:
  chef-deploy deploy --from=<from> --to=<to> [options]
  chef-deploy -h | --help
  chef-deploy --version

  Options:
  --from=<from>                 base SHA of the diff to deploy
  --to=<to>                     head SHA of the diff to deploy
  --knife-executable=<knife>    the knife executable to use
  -h --help                     Show this screen.
  --version                     Show version.
```
