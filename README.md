# cubelet (WIP)
kubectl configuration manipulation tools, which fetches kubeconfig for remote K8s cluster and merges into local one.

Only AWS cluster is supported now.

## Usage

```terminal

cubelet --remote_master_ip 172.xxx --name_suffix dev 

# then use kubens/kubectx to switch cluster
kubectx

```

## Note
1. `cubelet` leverages `SSH` and `SCP` for transfering files from remote cluster. Make sure SSH correctly configured.
