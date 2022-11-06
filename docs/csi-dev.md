# SSHFS CSI driver development guide

## How to build this project
 - Clone repo
```console
$ mkdir -p $GOPATH/src/sigs.k8s.io/
$ git clone https://github.com/robin-rpr/csi-driver-sshfs $GOPATH/src/github.com/robin-rpr/csi-driver-sshfs
```

 - Build CSI driver
```console
$ cd $GOPATH/src/github.com/robin-rpr/csi-driver-sshfs
$ make
```

 - Run verification test before submitting code
```console
$ make verify
```

 - If there is config file changed under `charts` directory, run following command to update chart file
```console
helm package charts/latest/csi-driver-sshfs -d charts/latest/
```

## How to test CSI driver in local environment

Install `csc` tool according to https://github.com/rexray/gocsi/tree/master/csc
```console
$ mkdir -p $GOPATH/src/github.com
$ cd $GOPATH/src/github.com
$ git clone https://github.com/rexray/gocsi.git
$ cd rexray/gocsi/csc
$ make build
```

#### Start CSI driver locally
```console
$ cd $GOPATH/src/github.com/robin-rpr/csi-driver-sshfs
$ ./bin/sshfsplugin --endpoint unix:///tmp/csi.sock --nodeid CSINode -v=5 &
```

#### 0. Set environment variables
```console
$ cap="1,mount,"
$ volname="test-$(date +%s)"
$ volsize="2147483648"
$ endpoint="unix:///tmp/csi.sock"
$ target_path="/tmp/targetpath"
$ params="server=127.0.0.1,share=/"
```

#### 1. Get plugin info
```console
$ csc identity plugin-info --endpoint "$endpoint"
"sshfs.csi.k8s.io"    "v2.0.0"
```

#### 2. Create a new sshfs volume
```console
$ value="$(csc controller new --endpoint "$endpoint" --cap "$cap" "$volname" --req-bytes "$volsize" --params "$params")"
$ sleep 15
$ volumeid="$(echo "$value" | awk '{print $1}' | sed 's/"//g')"
$ echo "Got volume id: $volumeid"
```

#### 3. Publish a sshfs volume
```
$ csc node publish --endpoint "$endpoint" --cap "$cap" --vol-context "$params" --target-path "$target_path" "$volumeid"
```

#### 4. Unpublish a sshfs volume
```console
$ csc node unpublish --endpoint "$endpoint" --target-path "$target_path" "$volumeid"
```

#### 6. Validate volume capabilities
```console
$ csc controller validate-volume-capabilities --endpoint "$endpoint" --cap "$cap" "$volumeid"
```

#### 7. Delete the sshfs volume
```console
$ csc controller del --endpoint "$endpoint" "$volumeid" --timeout 10m
```

#### 8. Get NodeID
```console
$ csc node get-info --endpoint "$endpoint"
CSINode
```

## How to test CSI driver in a Kubernetes cluster
- Set environment variable
```console
export REGISTRY=<dockerhub-alias>
export IMAGE_VERSION=latest
```

- Build continer image and push image to dockerhub
```console
# run `docker login` first
# build docker image
make container
# push the docker image
make push
```

- Deploy a Kubernetes cluster and make sure `kubectl get nodes` works on your dev box.

- Run E2E test on the Kubernetes cluster.

```console
# install SSHFS CSI Driver on the Kubernetes cluster
make e2e-bootstrap

# run the E2E test
make e2e-test
```
