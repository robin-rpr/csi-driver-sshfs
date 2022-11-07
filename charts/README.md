# Install CSI driver with Helm 3

## Prerequisites
 - [install Helm](https://helm.sh/docs/intro/quickstart/#install-helm)

### Tips
 - run controller on control plane node: `--set controller.runOnControlPlane=true`
 - set replica of controller as `2`: `--set controller.replicas=2`

### install a specific version
```console
helm repo add csi-driver-sshfs https://raw.githubusercontent.com/robin-rpr/csi-driver-sshfs/master/charts
helm install csi-driver-sshfs csi-driver-sshfs/csi-driver-sshfs --namespace kube-system --version v4.1.0
```

### install driver with customized driver name, deployment name
> only supported from `v3.1.0`+
 - following example would install a driver with name `sshfs2`
```console
helm install csi-driver-sshfs2 csi-driver-sshfs/csi-driver-sshfs --namespace kube-system --set driver.name="sshfs2.csi.k8s.io" --set controller.name="csi-sshfs2-controller" --set rbac.name=sshfs2 --set serviceAccount.controller=csi-sshfs2-controller-sa --set serviceAccount.node=csi-sshfs2-node-sa --set node.name=csi-sshfs2-node --set node.livenessProbe.healthPort=39653
```

### search for all available chart versions
```console
helm search repo -l csi-driver-sshfs
```

## uninstall CSI driver
```console
helm uninstall csi-driver-sshfs -n kube-system
```

## latest chart configuration

The following table lists the configurable parameters of the latest SSHFS CSI Driver chart and default values.

| Parameter                                         | Description                                                | Default                                                           |
|---------------------------------------------------|------------------------------------------------------------|-------------------------------------------------------------------|
| `customLabels`                                    | optional extra labels to k8s resources deployed by chart   | `{}`                                                              |
| `driver.name`                                     | alternative driver name                                    | `sshfs.csi.k8s.io` |
| `driver.mountPermissions`                         | mounted folder permissions name                            | `0777`
| `feature.enableFSGroupPolicy`                     | enable `fsGroupPolicy` on a k8s 1.20+ cluster              | `true`                      |
| `feature.enableInlineVolume`                      | enable inline volume                     | `false`                      |
| `kubeletDir`                                      | alternative kubelet directory                              | `/var/lib/kubelet`                                                  |
| `image.sshfs.repository`                            | csi-driver-sshfs image                                       | `registry.k8s.io/sig-storage/sshfsplugin`                          |
| `image.sshfs.tag`                                   | csi-driver-sshfs image tag                                   | `latest`                                                |
| `image.sshfs.pullPolicy`                            | csi-driver-sshfs image pull policy                           | `IfNotPresent`                                                      |
| `image.csiProvisioner.repository`                 | csi-provisioner docker image                               | `registry.k8s.io/sig-storage/csi-provisioner`                            |
| `image.csiProvisioner.tag`                        | csi-provisioner docker image tag                           | `v3.2.0`                                                            |
| `image.csiProvisioner.pullPolicy`                 | csi-provisioner image pull policy                          | `IfNotPresent`                                                      |
| `image.livenessProbe.repository`                  | liveness-probe docker image                                | `registry.k8s.io/sig-storage/livenessprobe`                              |
| `image.livenessProbe.tag`                         | liveness-probe docker image tag                            | `v2.7.0`                                                            |
| `image.livenessProbe.pullPolicy`                  | liveness-probe image pull policy                           | `IfNotPresent`                                                      |
| `image.nodeDriverRegistrar.repository`            | csi-node-driver-registrar docker image                     | `registry.k8s.io/sig-storage/csi-node-driver-registrar`                  |
| `image.nodeDriverRegistrar.tag`                   | csi-node-driver-registrar docker image tag                 | `v2.5.1`                                                            |
| `image.nodeDriverRegistrar.pullPolicy`            | csi-node-driver-registrar image pull policy                | `IfNotPresent`                                                      |
| `imagePullSecrets`                                | Specify docker-registry secret names as an array           | [] (does not add image pull secrets to deployed pods)                                                           |
| `serviceAccount.create`                           | whether create service account of csi-sshfs-controller       | `true`                                                              |
| `rbac.create`                                     | whether create rbac of csi-sshfs-controller                  | `true`                                                              |
| `controller.replicas`                             | replica number of csi-sshfs-controller                         | `1`                                                                 |
| `controller.runOnMaster`                          | run controller on master node(deprecated on k8s 1.25+)                                                          |`false`                                                           |
| `controller.runOnControlPlane`                    | run controller on control plane node                                                          |`false`                                                           |
| `controller.dnsPolicy`                            | dnsPolicy of controller driver, available values: `Default`, `ClusterFirstWithHostNet`, `ClusterFirst`                              | `Default`                                                             |
| `controller.logLevel`                             | controller driver log level                                                          |`5`                                                           |
| `controller.workingMountDir`                      | working directory for provisioner to mount sshfs shares temporarily                  | `/tmp`                                                             |
| `controller.affinity`                                 | controller pod affinity                               | `{}`                                                             |
| `controller.nodeSelector`                             | controller pod node selector                          | `{}`                                                             |
| `controller.tolerations`                              | controller pod tolerations                            |                                                              |
| `controller.resources.csiProvisioner.limits.memory`   | csi-provisioner memory limits                         | 100Mi                                                          |
| `controller.resources.csiProvisioner.requests.cpu`    | csi-provisioner cpu requests limits                   | 10m                                                            |
| `controller.resources.csiProvisioner.requests.memory` | csi-provisioner memory requests limits                | 20Mi                                                           |
| `controller.resources.livenessProbe.limits.memory`    | liveness-probe memory limits                          | 100Mi                                                          |
| `controller.resources.livenessProbe.requests.cpu`     | liveness-probe cpu requests limits                    | 10m                                                            |
| `controller.resources.livenessProbe.requests.memory`  | liveness-probe memory requests limits                 | 20Mi                                                           |
| `controller.resources.sshfs.limits.memory`              | csi-driver-sshfs memory limits                         | 200Mi                                                          |
| `controller.resources.sshfs.requests.cpu`               | csi-driver-sshfs cpu requests limits                   | 10m                                                            |
| `controller.resources.sshfs.requests.memory`            | csi-driver-sshfs memory requests limits                | 20Mi                                                           |
| `node.name`                                           | driver node daemonset name                            | `csi-sshfs-node`
| `node.dnsPolicy`                                      | dnsPolicy of driver node daemonset, available values: `Default`, `ClusterFirstWithHostNet`, `ClusterFirst`                              |
| `node.maxUnavailable`                             | `maxUnavailable` value of driver node daemonset                            | `1`
| `node.logLevel`                                   | node driver log level                                                          |`5`                                                           |
| `node.livenessProbe.healthPort `                  | the health check port for liveness probe                    |`29653`                                                           |
| `node.affinity`                                      | node pod affinity                                     | {}                                                             |
| `node.nodeSelector`                                   | node pod node selector                                | `{}`                                                             |
| `node.tolerations`                              | node pod tolerations                            |                                                              |
| `node.resources.livenessProbe.limits.memory`          | liveness-probe memory limits                          | 100Mi                                                          |
| `node.resources.livenessProbe.requests.cpu`           | liveness-probe cpu requests limits                    | 10m                                                            |
| `node.resources.livenessProbe.requests.memory`        | liveness-probe memory requests limits                 | 20Mi                                                           |
| `node.resources.nodeDriverRegistrar.limits.memory`    | csi-node-driver-registrar memory limits               | 100Mi                                                          |
| `node.resources.nodeDriverRegistrar.requests.cpu`     | csi-node-driver-registrar cpu requests limits         | 10m                                                            |
| `node.resources.nodeDriverRegistrar.requests.memory`  | csi-node-driver-registrar memory requests limits      | 20Mi                                                           |
| `node.resources.sshfs.limits.memory`                    | csi-driver-sshfs memory limits                         | 300Mi                                                         |
| `node.resources.sshfs.requests.cpu`                     | csi-driver-sshfs cpu requests limits                   | 10m                                                            |
| `node.resources.sshfs.requests.memory`                  | csi-driver-sshfs memory requests limits                | 20Mi                                                           |

## troubleshooting
 - Add `--wait -v=5 --debug` in `helm install` command to get detailed error
 - Use `kubectl describe` to acquire more info
