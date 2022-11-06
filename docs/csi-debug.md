## CSI driver debug tips

### case#1: volume create/delete failed
> There could be multiple controller pods (only one pod is the leader), if there are no helpful logs, try to get logs from the leader controller pod.
 - find csi driver controller pod
```console
$ kubectl get pod -o wide -n kube-system | grep csi-sshfs-controller
NAME                                     READY   STATUS    RESTARTS   AGE     IP             NODE
csi-sshfs-controller-56bfddd689-dh5tk      5/5     Running   0          35s     10.240.0.19    k8s-agentpool-22533604-0
csi-sshfs-controller-56bfddd689-sl4ll      5/5     Running   0          35s     10.240.0.23    k8s-agentpool-22533604-1
```
 - get pod description and logs
```console
$ kubectl describe csi-sshfs-controller-56bfddd689-dh5tk -n kube-system > csi-sshfs-controller-description.log
$ kubectl logs csi-sshfs-controller-56bfddd689-dh5tk -c sshfs -n kube-system > csi-sshfs-controller.log
```

### case#2: volume mount/unmount failed
 - locate csi driver pod that does the actual volume mount/unmount

```console
$ kubectl get pod -o wide -n kube-system | grep csi-sshfs-node
NAME                                      READY   STATUS    RESTARTS   AGE     IP             NODE
csi-sshfs-node-cvgbs                        3/3     Running   0          7m4s    10.240.0.35    k8s-agentpool-22533604-1
csi-sshfs-node-dr4s4                        3/3     Running   0          7m4s    10.240.0.4     k8s-agentpool-22533604-0
```

 - get pod description and logs
```console
$ kubectl describe po csi-sshfs-node-cvgbs -n kube-system > csi-sshfs-node-description.log
$ kubectl logs csi-sshfs-node-cvgbs -c sshfs -n kube-system > csi-sshfs-node.log
```

 - check sshfs mount inside driver
```console
kubectl exec -it csi-sshfs-node-cvgbss -n kube-system -c sshfs -- mount | grep sshfs
```

### troubleshooting connection failure on agent node
```console
mkdir /tmp/test
mount -v -t sshfs -o ... sshfs-server:/path /tmp/test
```
