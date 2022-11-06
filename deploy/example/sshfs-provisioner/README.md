# Set up a SSHFS Server on a Kubernetes cluster

After the SSHFS CSI Driver is deployed in your cluster, you can follow this documentation to quickly deploy some example applications. You can use SSHFS CSI Driver to provision Persistent Volumes statically or dynamically. Please read Kubernetes Persistent Volumes for more information about Static and Dynamic provisioning.

There are multiple different SSHFS servers you can use for testing of 
the plugin, the major versions of the protocol v2, v3 and v4 should be supported
by the current implementation. This page will show you how to set up a SSHFS Server deployment on a Kubernetes cluster.

- Your SSHFS provisioner needs a Key Pair from NGINX webserver for Authentication, generate one below. 
> Default all SSH Keys are supported: RSA, DSA, ECDSA, or EdDSA.

```bash
ssh-keygen -f ./id_rsa_nginx_pod -t rsa -b 4096
kubectl create secret generic nginx-pod-rsa-key --from-file=id_rsa=id_rsa_nginx_pod
kubectl create secret generic nginx-pod-rsa-key-public --from-file=id_rsa.pub=id_rsa_nginx_pod.pub

# To delete the keys from your local device, run the command below.
rm id_rsa_nginx_pod id_rsa_nginx_pod.pub
```

- To create a SSHFS provisioner on your Kubernetes cluster, run the following command.

```bash
kubectl create -f https://raw.githubusercontent.com/robin-rpr/csi-driver-sshfs/master/deploy/example/sshfs-provisioner/sshfs-server.yaml
```

- During the deployment, a new service `sshfs-server` will be created which exposes the SSHFS server endpoint `sshfs-server.default.svc.cluster.local` and the share path `/`. You can specify `PersistentVolume` or `StorageClass` using these information.

- Deploy the SSHFS CSI driver, please refer to [install SSHFS CSI driver](../../../docs/install-sshfs-csi-driver.md).

- To check if the SSHFS server is working, we can statically create a PersistentVolume and a PersistentVolumeClaim, and mount it onto a sample pod:

```bash
kubectl create -f https://raw.githubusercontent.com/robin-rpr/csi-driver-sshfs/master/deploy/example/sshfs-provisioner/nginx-pod.yaml
```

 - Verify if the SSHFS server is functional, you can check the mount point from the example pod.

 ```bash
kubectl exec nginx-sshfs-example -- bash -c "findmnt /var/www -o TARGET,SOURCE,FSTYPE"
```

 - The output should look like the following:

 ```bash
TARGET   SOURCE                                   FSTYPE
/var/www sshfs-server.default.svc.cluster.local:/ sshfs
```
