# CSI driver example

You can use SSHFS CSI Driver to provision Persistent Volumes statically or dynamically. Please read [Kubernetes Persistent Volumes documentation](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) for more information about Static and Dynamic provisioning.

Please refer to [driver parameters](../../docs/driver-parameters.md) for more detailed usage.

## Prerequisite

- [Set up a SSHFS Server on a Kubernetes cluster](./sshfs-provisioner/README.md) as an example
- [Install SSHFS CSI Driver](../../docs/install-sshfs-csi-driver.md)

## Storage Class Usage (Dynamic Provisioning)

 -  Create a storage class
 > change `server`, `share`, `privateKey` with your existing SSHFS server address, share name, and client private key
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: sshfs-csi
provisioner: sshfs.csi.k8s.io
parameters:
  server: sshfs-server.default.svc.cluster.local
  share: /
  private-key-name: "nginx-pod-rsa-key"
  private-key-namespace: "default"
  port: "2222"
  user: "admin"
  # csi.storage.k8s.io/provisioner-secret is only needed for providing mountOptions in DeleteVolume
  # csi.storage.k8s.io/provisioner-secret-name: "mount-options"
  # csi.storage.k8s.io/provisioner-secret-namespace: "default"
reclaimPolicy: Delete
volumeBindingMode: Immediate
mountOptions:
  - sshfsvers=3.7.3
```

 - create PVC
```console
kubectl create -f https://raw.githubusercontent.com/robin-rpr/csi-driver-sshfs/master/deploy/example/pvc-sshfs-csi-dynamic.yaml
```

## PV/PVC Usage (Static Provisioning)

- Follow the following command to create `PersistentVolume` and `PersistentVolumeClaim` statically.

```bash
# create PV
kubectl create -f https://raw.githubusercontent.com/robin-rpr/csi-driver-sshfs/master/deploy/example/pv-sshfs-csi.yaml

# create PVC
kubectl create -f https://raw.githubusercontent.com/robin-rpr/csi-driver-sshfs/master/deploy/example/pvc-sshfs-csi-static.yaml
```

## Create a deployment
```console
kubectl create -f https://raw.githubusercontent.com/robin-rpr/csi-driver-sshfs/master/deploy/example/deployment.yaml
```
