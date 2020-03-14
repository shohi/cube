# cube
kubectl configuration manipulation tool, which fetches `~/.kube/config` from remote K8s cluster and merges it into local one.

## Prerequisite
`cube` depends on `ssh tunnel` for communication with remote cluster. Make sure these files -- `~/.ssh/config` and `/etc/hosts` -- are correctly set.

### `~/.ssh/config`

```terminal
# add SSH dynamic port forwarding
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
  version     print version info

Flags:
  -h, --help   help for cube

```

use [kubectx](https://github.com/ahmetb/kubectx) to switch cluster

```
$> kubectx
```

### Add

```
$> cube add
add remote cluster to kube config

Usage:
  cube add [flags]

Flags:
      --dry-run                dry-run mode. validate config and then exit
      --force                  merge configuration forcedly. Only take effect when cluster name is unique
  -h, --help                   help for add
      --local-port int         local forwarding port
      --name-suffix string     cluster name suffix
      --print-ssh-forwarding   print ssh forwarding command and exit
      --remote-ip string       remote master private ip
      --remote-user string     remote user (default "core")
      --ssh-via string         ssh jump server, e.g. user@jump. If not set, SSH_VIA env will be used
```

### Delete

```
$> cube delete
delete kubectl config for specified cluster

Usage:
  cube delete [flags]

Aliases:
  delete, del

Flags:
      --all           delete all matched cluster.
      --dry-run       dry-run mode. print modified config and exit
  -h, --help          help for delete
      --name string   cluster name to delete
```

~~### Docker~~

```terminal

docker run --rm -it \
    -v $PWD/.ssh:/root/.ssh \
    -v $PWD/.kube:/root/.kube \
    cube:0.1.0

```

## Note
1. `cube` leverages `SSH` and `SCP` for transfering files from remote cluster. Make sure SSH correctly configured.

2. Only AWS cluster is supported now.
