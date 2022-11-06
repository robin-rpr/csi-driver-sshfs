/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"

	nfs "github.com/robin-rpr/csi-driver-sshfs/pkg/sshfs"

	"k8s.io/klog/v2"
)

var (
	endpoint         = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	nodeID           = flag.String("nodeid", "", "node id")
	mountPermissions = flag.Uint64("mount-permissions", 0777, "mounted folder permissions")
	driverName       = flag.String("drivername", nfs.DefaultDriverName, "name of the driver")
	workingMountDir  = flag.String("working-mount-dir", "/tmp", "working directory for provisioner to mount nfs shares temporarily")
)

func init() {
	_ = flag.Set("logtostderr", "true")
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	if *nodeID == "" {
		klog.Warning("nodeid is empty")
	}

	handle()
	os.Exit(0)
}

func handle() {
	driverOptions := nfs.DriverOptions{
		NodeID:           *nodeID,
		DriverName:       *driverName,
		Endpoint:         *endpoint,
		MountPermissions: *mountPermissions,
		WorkingMountDir:  *workingMountDir,
	}
	d := nfs.NewDriver(&driverOptions)
	d.Run(false)
}
