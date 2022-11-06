# Install SSHFS CSI driver master version on a kubernetes cluster

If you have already installed Helm, you can also use it to install this driver. Please check [Installation with Helm](../charts/README.md).

## Install with kubectl
 - Option#1. remote install
```console
curl -skSL https://raw.githubusercontent.com/robin-rpr/csi-driver-sshfs/master/deploy/install-driver.sh | bash -s master --
```

 - Option#2. local install
```console
git clone https://github.com/robin-rpr/csi-driver-sshfs.git
cd csi-driver-sshfs
./deploy/install-driver.sh master local
```

- check pods status:
```console
kubectl -n kube-system get pod -o wide -l app=csi-sshfs-controller
kubectl -n kube-system get pod -o wide -l app=csi-sshfs-node
```

example output:

```console
NAME                                       READY   STATUS    RESTARTS   AGE     IP             NODE
csi-sshfs-controller-56bfddd689-dh5tk       4/4     Running   0          35s     10.240.0.19    k8s-agentpool-22533604-0
csi-sshfs-node-cvgbs                        3/3     Running   0          35s     10.240.0.35    k8s-agentpool-22533604-1
csi-sshfs-node-dr4s4                        3/3     Running   0          35s     10.240.0.4     k8s-agentpool-22533604-0
```

### clean up SSHFS CSI driver
 - Option#1. remote uninstall
```console
curl -skSL https://raw.githubusercontent.com/robin-rpr/csi-driver-sshfs/master/deploy/uninstall-driver.sh | bash -s master --
```

 - Option#2. local uninstall
```console
git clone https://github.com/robin-rpr/csi-driver-sshfs.git
cd csi-driver-sshfs
git checkout master
./deploy/uninstall-driver.sh master local
```
