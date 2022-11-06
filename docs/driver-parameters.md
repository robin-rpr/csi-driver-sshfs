## Driver Parameters
> This driver requires existing and already configured SSHFSv3 or SSHFSv4 server, it supports dynamic provisioning of Persistent Volumes via Persistent Volume Claims by creating a new sub directory under SSHFS server.

### storage class usage (dynamic provisioning)
> [`StorageClass` example](../deploy/example/storageclass-sshfs.yaml)

Name | Meaning | Example Value | Mandatory | Default value
--- | --- | --- | --- | ---
server | SSHFS Server address | domain name `sshfs-server.default.svc.cluster.local` <br>or IP address `127.0.0.1` | Yes |
share | SSHFS share path | `/` | Yes |
port | SSHFS Server port | `2222` | Yes |
private-key-name | Secret holding SSHFS Server authorized private key | `nginx-pod-rsa-key` | Yes |
private-key-namespace | Namespace of Secret holding SSHFS Server authorized private key | `default` | No | `default` |
user | SSHFS Server username | `admin` | Yes |
subDir | sub directory under sshfs share |  | No | if sub directory does not exist, this driver would create a new one
mountPermissions | mounted folder permissions. The default is `0777`, if set as `0`, driver will not perform `chmod` after mount |  | No |

### PV/PVC usage (static provisioning)
> [`PersistentVolume` example](../deploy/example/pv-sshfs-csi.yaml)

Name | Meaning | Example Value | Mandatory | Default value
--- | --- | --- | --- | ---
volumeAttributes.server | SSHFS Server address | domain name `sshfs-server.default.svc.cluster.local` <br>or IP address `127.0.0.1` | Yes |
volumeAttributes.share | SSHFS share path | `/` |  Yes  |
volumeAttributes.mountPermissions | mounted folder permissions. The default is `0777` |  | No |

### Tips
#### `subDir` parameter supports following pv/pvc metadata conversion
> if `subDir` value contains following strings, it would be converted into corresponding pv/pvc name or namespace
 - `${pvc.metadata.name}`
 - `${pvc.metadata.namespace}`
 - `${pv.metadata.name}`

#### provide `mountOptions` for `DeleteVolume`
> since `DeleteVolumeRequest` does not provide `mountOptions`, following is the workaround to provide `mountOptions` for `DeleteVolume`
  - create a secret with `mountOptions`
```console
kubectl create secret generic mount-options --from-literal mountOptions="sshfsvers=3.7.3,hard"
```
  - define a storage class with `csi.storage.k8s.io/provisioner-secret-name` and `csi.storage.k8s.io/provisioner-secret-namespace` setting:
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: sshfs-csi
provisioner: sshfs.csi.k8s.io
parameters:
  server: sshfs-server.default.svc.cluster.local
  share: /
  # csi.storage.k8s.io/provisioner-secret is only needed for providing mountOptions in DeleteVolume
  csi.storage.k8s.io/provisioner-secret-name: "mount-options"
  csi.storage.k8s.io/provisioner-secret-namespace: "default"
reclaimPolicy: Delete
volumeBindingMode: Immediate
mountOptions:
  - sshfsvers=3.7.3
```
