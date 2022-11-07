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

if [[ -z "$(command -v csc)" ]]; then
  GO111MODULE=off go get github.com/rexray/gocsi/csc
fi

function cleanup {
  echo 'pkill -f sshfsplugin'
  pkill -f sshfsplugin
  echo 'Uninstalling SSHFS server on localhost'
  docker rm sshfs -f
}
trap cleanup EXIT

function provision_sshfs_server {
  echo 'Installing SSHFS server on localhost'
  apt-get update -y
  apt-get install -y sshfs-common
  docker run -d --name sshfs --privileged -p 2049:2049 -v "$(pwd)"/sshfsshare:/sshfsshare -e SHARED_DIRECTORY=/sshfsshare itsthenetwork/sshfs-server-alpine:latest
}

provision_sshfs_server

readonly CSC_BIN="$GOBIN/csc"
readonly cap="1,mount,"
volname="citest-$(date +%s)"
readonly volsize="2147483648"
readonly endpoint="unix:///tmp/csi.sock"
readonly target_path="/tmp/targetpath"
readonly params="server=127.0.0.1,share=/"

nodeid='CSINode'
if [[ "$#" -gt 0 ]] && [[ -n "$1" ]]; then
  nodeid="$1"
fi

# Run CSI driver as a background service
bin/sshfsplugin --endpoint "$endpoint" --nodeid "$nodeid" -v=5 &
sleep 5

echo 'Begin to run integration test...'

# Begin to run CSI functions one by one
echo "Create volume test:"
value="$("$CSC_BIN" controller new --endpoint "$endpoint" --cap "$cap" "$volname" --req-bytes "$volsize" --params "$params")"
sleep 15

volumeid="$(echo "$value" | awk '{print $1}' | sed 's/"//g')"
echo "Got volume id: $volumeid"

"$CSC_BIN" controller validate-volume-capabilities --endpoint "$endpoint" --cap "$cap" "$volumeid"

echo "publish volume test:"
"$CSC_BIN" node publish --endpoint "$endpoint" --cap "$cap" --vol-context "$params" --target-path "$target_path" "$volumeid"
sleep 2

declare staging_target_path
echo "node stats test:"
csc node stats --endpoint "$endpoint" "$volumeid:$target_path:$staging_target_path"
sleep 2

echo "unpublish volume test:"
"$CSC_BIN" node unpublish --endpoint "$endpoint" --target-path "$target_path" "$volumeid"
sleep 2

echo "Delete volume test:"
"$CSC_BIN" controller del --endpoint "$endpoint" "$volumeid" --timeout 10m
sleep 15

"$CSC_BIN" identity plugin-info --endpoint "$endpoint"
"$CSC_BIN" node get-info --endpoint "$endpoint"

echo "Integration test is completed."
