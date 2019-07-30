# cube (WIP)
kubectl configuration manipulation tools, which fetches `~/.kube/config` from remote K8s cluster and merges it into local one.

## Usage

```terminal
# example
cube \
    --remote-addr=core@172.xxx \
    --local-port=7001 \
    --ssh-via user@jump-server \
    --name-suffix=dev \
    --dry-run=false

# help
cube --help

# then use kubens/kubectx to switch cluster
kubectx

```

## Note
1. `cube` leverages `SSH` and `SCP` for transfering files from remote cluster. Make sure SSH correctly configured.

2. Only AWS cluster is supported now.
