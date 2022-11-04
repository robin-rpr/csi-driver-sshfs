**Status:** Experimental (only for development, expect crashes)
<img height="151px" align="right" alt="csi-driver-sshfs" src="https://raw.githubusercontent.com/robin-rpr/csi-driver-sshfs/master/csi-driver-sshfs.svg" title="csi-driver-sshfs"/>

# SSHFS CSI driver for Kubernetes
![build status](https://github.com/robin-rpr/csi-driver-sshfs/actions/workflows/linux.yaml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/robin-rpr/csi-driver-sshfs/badge.svg?branch=master)](https://coveralls.io/github/robin-rpr/csi-driver-sshfs?branch=master)

### Overview

This is a repository for [SSHFS](https://en.wikipedia.org/wiki/SSHFS) [CSI](https://kubernetes-csi.github.io/docs/) driver, csi plugin name: `sshfs.csi.k8s.io`. This driver requires existing and already configured SSHFS server, it supports dynamic provisioning of Persistent Volumes via Persistent Volume Claims by creating a new sub directory under SSHFS server.

### Project status: pre-alpha [Learn more ...](https://en.wikipedia.org/wiki/Software_release_life_cycle)

### Container Images & Kubernetes Compatibility:
|driver version  | supported k8s version | status    |
|----------------|-----------------------|-----------|
|master branch   | 1.20+                 | pre-alpha |

### Install driver on a Kubernetes cluster
 - install by [kubectl](./docs/install-sshfs-csi-driver.md)
 - install by [helm charts](./charts)

### Driver parameters
Please refer to [`sshfs.csi.k8s.io` driver parameters](./docs/driver-parameters.md)

### Examples
 - [Basic usage](./deploy/example/README.md)
 - [fsGroupPolicy](./deploy/example/fsgroup)

### Troubleshooting
 - [CSI driver troubleshooting guide](./docs/csi-debug.md) 

## Kubernetes Development
Please refer to [development guide](./docs/csi-dev.md)

### Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [Slack channel](https://kubernetes.slack.com/messages/sig-storage)
- [Mailing list](https://groups.google.com/forum/#!forum/kubernetes-sig-storage)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

[owners]: https://git.k8s.io/community/contributors/guide/owners.md
[Creative Commons 4.0]: https://git.k8s.io/website/LICENSE
