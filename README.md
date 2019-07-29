# cube (WIP)
kubectl configuration manipulation tools, which fetches `~/.kube/config` from remote K8s cluster and merges it into local one.

## Usage

```terminal

cube --remote_master_ip 172.xxx --name_suffix dev 

# then use kubens/kubectx to switch cluster
kubectx

```

## Note
1. `cubelet` leverages `SSH` and `SCP` for transfering files from remote cluster. Make sure SSH correctly configured.

2. Only AWS cluster is supported now.
