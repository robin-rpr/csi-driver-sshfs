/*
Copyright 2020 The Kubernetes Authors.

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

package testsuites

import (
	"github.com/kubernetes-csi/csi-driver-nfs/test/e2e/driver"
	"github.com/robin-rpr/csi-driver-nfs/pkg/nfs"

	v1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
)

// DynamicallyProvisionedReclaimPolicyTest will provision required PV(s) and PVC(s)
// Testing the correct behavior for different reclaimPolicies
type DynamicallyProvisionedReclaimPolicyTest struct {
	CSIDriver              driver.DynamicPVTestDriver
	Volumes                []VolumeDetails
	ControllerServer       nfs.ControllerServer
	StorageClassParameters map[string]string
}

func (t *DynamicallyProvisionedReclaimPolicyTest) Run(client clientset.Interface, namespace *v1.Namespace) {
	for _, volume := range t.Volumes {
		tpvc, _ := volume.SetupDynamicPersistentVolumeClaim(client, namespace, t.CSIDriver, t.StorageClassParameters)

		// will delete the PVC
		// will also wait for PV to be deleted when reclaimPolicy=Delete
		tpvc.Cleanup()
		// first check PV stills exists, then manually delete it
		if tpvc.ReclaimPolicy() == v1.PersistentVolumeReclaimRetain {
			tpvc.WaitForPersistentVolumePhase(v1.VolumeReleased)
			tpvc.DeleteBoundPersistentVolume()
			// The controller server cannot resolve the nfs server hosting inside the testing k8s cluster, skipping the cleanup step.
			// tpvc.DeleteBackingVolume(&t.ControllerServer)
		}
	}
}
