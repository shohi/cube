# cubelet (WIP)
kubectl configuration manipulation tools for remote cluster. Only AWS cluster is supported now.

## Usage

```terminal

cubelet --remote_master_ip 172.xxx --name_suffix dev 

```

## Note
1. `cubelet` leverages `SSH` and `SCP` for transfering files from remote cluster. Make sure SSH correctly configured.
