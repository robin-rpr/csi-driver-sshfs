#!/bin/bash

# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eo pipefail

function cleanup {
  echo 'pkill -f sshfsplugin'
  pkill -f sshfsplugin
  echo 'Deleting CSI sanity test binary'
  rm -rf csi-test
  echo 'Uninstalling SSHFS server on localhost'
  docker rm sshfs -f
}
trap cleanup EXIT

function install_csi_sanity_bin {
  echo 'Installing CSI sanity test binary...'
  mkdir -p $GOPATH/src/github.com/kubernetes-csi
  pushd $GOPATH/src/github.com/kubernetes-csi
  export GO111MODULE=off
  git clone https://github.com/kubernetes-csi/csi-test.git -b v4.2.0
  pushd csi-test/cmd/csi-sanity
  make install
  popd
  popd
}

function provision_sshfs_server {
  echo 'Installing SSHFS server on localhost'
  apt-get update -y
  apt-get install -y sshfs-common
  docker run -d --name sshfs --privileged -p 2049:2049 -v "$(pwd)"/sshfsshare:/sshfsshare -e SHARED_DIRECTORY=/sshfsshare itsthenetwork/sshfs-server-alpine:latest
}

provision_sshfs_server
# wait for sshfs-server running complete
sleep 10
install_csi_sanity_bin

readonly endpoint='unix:///tmp/csi.sock'
nodeid='CSINode'
if [[ "$#" -gt 0 ]] && [[ -n "$1" ]]; then
  nodeid="$1"
fi

bin/sshfsplugin --endpoint "$endpoint" --nodeid "$nodeid" -v=5 &

echo 'Begin to run sanity test...'
readonly CSI_SANITY_BIN='csi-sanity'
"$CSI_SANITY_BIN" --ginkgo.v --csi.testvolumeparameters="$(pwd)/test/sanity/params.yaml" --csi.endpoint="$endpoint" --ginkgo.skip="should not fail when requesting to create a volume with already existing name and same capacity|should fail when requesting to create a volume with already existing name and different capacity|should work|should fail when the requested volume does not exist|should return appropriate capabilities"
