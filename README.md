# cube
kubectl configuration manipulation tools, which fetches `~/.kube/config` from remote K8s cluster and merges it into local one.

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
# add following line into host
# use to access AWS k8s cluster by SSH tunnel
sudo echo "127.0.0.1	kubernetes" >> /etc/hosts
```

## Install

```
go get -u github.com/shohi/cube

```

## Usage

### Binary
```terminal
# `merge` example
cube add \
    --remote-user=core \
    --remote-ip=172.xxx \
    --ssh-via user@jump-server \
    --local-port=7001 \
    --name-suffix=dev \
    --dry-run=false

# `purge` example
cube del \
    --remote-user=core \
    --remote-ip=172.xxx \
    --ssh-via user@jump-server \
    --dry-run=false

# help
cube --help

# then use kubens/kubectx to switch cluster
kubectx

```

### Docker

```terminal

docker run --rm -it \
    -v $PWD/.ssh:/root/.ssh \
    -v $PWD/.kube:/root/.kube \
    cube:0.1.0

```

## Note
1. `cube` leverages `SSH` and `SCP` for transfering files from remote cluster. Make sure SSH correctly configured.

2. Only AWS cluster is supported now.
