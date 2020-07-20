# CoGo
CoGo is a cli utility for managing cognito users and groups written in Go.
This is mainly a side project to learn go and is my first go project.

## Installation
I'll do the whole release thing soon but hopefully you just curl your version:
```
curl TODO
```

Alternatively you can use the Docker images now:
```
docker run -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} -e AWS_DEFAULT_REGION=us-west-2 unstableunicorn/cogo -p <poolid> list users
```
or if you have a .aws folder with credentials you can mount the folder and pass just the parameters you need
```
docker run --rm -e AWS_PROFILE=cogo-dev -v ${HOME}/.aws:/.aws unstableunicorn/cogo -p <poolid> list users
```

## Usage
Must pass the poolid of the pool to manage and then can create, update, list and delete users and groups. Get started with the help:
```
>cogo --help
Usage: cogo [OPTIONS] [COMMAND]
  Cogo (short for Cognito Go)  is a cli written in Go that allows
  you to create, update, list and delete cognito users and groups including
  filtering and providing the ability to bulk update users.
  
  Examples:
  To list users:
  >cogo -p <poolid> list users
  
  To list groups:
  >cogo -p <poolid> list groups
  
  To list users and only show certain attributes:
  >cogo -p <poolid> list users --attr username email status custom:somecustomattribute
  
  To create a user with sane defaults and add to existing groups:
  >cogo -p <poolid> add user first.last@organisation.com --groups grp1 grp2

  Shortcuts! to make life easier you can use the following aliases:
  list|ls
  users|user|usr|u
  groups|group|grp|g
  e.g. cogo -p <poolid> list users -> cogo ls u

  You can also enter the poolid anywhere on the command yay:
  >cogo create user username -email -p <poolid>
  >cogo create user -p <poolid> username -email   
  >cogo create -p <poolid> user username -email

Usage:
  cogo [command]

Available Commands:
  create      Create cognito users or groups
  delete      Delete cognito users or groups
  help        Help about any command
  list        list cognito users or groups
  update      update cognito users or groups

Flags:
  -h, --help            help for cogo
  -p, --poolid string   AWS Cognito User PoolID (required)
  -v, --version         version for cogo

Use "cogo [command] --help" for more information about a command.
```

Lets leave it there for now