# cube (WIP)
kubectl configuration manipulation tools, which fetches `~/.kube/config` from remote K8s cluster and merges it into local one.

## Usage

```terminal
# example
cube --remote_ip=core@172.xxx --local_port=7001 --ssh-via user@jump-server --name_suffix=dev

# help
cube --help

# then use kubens/kubectx to switch cluster
kubectx

```

## Note
1. `cube` leverages `SSH` and `SCP` for transfering files from remote cluster. Make sure SSH correctly configured.

2. Only AWS cluster is supported now.
