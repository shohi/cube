# cube
kubectl configuration manipulation tool, which fetches `~/.kube/config` from remote K8s cluster and merges it into local one.

## Prerequisite
`cube` depends on `ssh tunnel` for communication with remote cluster. Make sure these files -- `~/.ssh/config` and `/etc/hosts` -- are correctly set.

### `~/.ssh/config`

```terminal
# add SSH dynamic port forwarding, where `SSH_VIA` is in the format of "<user>@<public-ip>"
# alias aws_proxy='ssh -qTfnN -D 127.0.0.1:62222 ${SSH_VIA}'

# Rules for Remote
Host [remote/private/ip/range(e.g. 172.31.*)]
 ProxyCommand /usr/bin/nc -X 4 -x 127.0.0.1:62222 %h %p
 StrictHostKeyChecking no
 UserKnownHostsFile=/dev/null
 user core
 IdentityFile [/path/to/pem/file]
 LogLevel ERROR

```

### `/etc/hosts`

```terminal
# add following line to /etc/hosts
# to access AWS k8s cluster by SSH tunnel
sudo echo "127.0.0.1	kubernetes" >> /etc/hosts
```

## Install

```
go get -u github.com/shohi/cube

```

## Usage

### help

```
$> cube --help
kubectl config manipulation tool

Usage:
  cube [command]

Available Commands:
  add         add remote cluster to kube config
  delete      delete kubectl config for specified cluster
  forward     run local ssh port forwarding for remote cluster
  help        Help about any command
  history     show cube commands history
  list        list all clusters
  show        show local kubectl config
  version     print version info

Flags:
  -h, --help   help for cube

Use "cube [command] --help" for more information about a command.
```

use [kubectx](https://github.com/ahmetb/kubectx) to switch cluster

```
$> kubectx
```

~~### Docker~~

```terminal

docker run --rm -it \
    -v $PWD/.ssh:/root/.ssh \
    -v $PWD/.kube:/root/.kube \
    cube:0.4.1

```

## Note
1. `cube` leverages `SSH` and `SCP` for transfering files from remote cluster. Make sure SSH correctly configured.

2. Only AWS cluster is supported now.
