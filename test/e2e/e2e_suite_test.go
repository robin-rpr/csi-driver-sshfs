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

package e2e

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/pborman/uuid"
	"github.com/robin-rpr/csi-driver-sshfs/pkg/sshfs"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/framework/config"
)

const (
	kubeconfigEnvVar   = "KUBECONFIG"
	testWindowsEnvVar  = "TEST_WINDOWS"
	sshfsServerAddress = "sshfs-server.default.svc.cluster.local"
	sshfsShare         = "/"
)

var (
	nodeID                        = os.Getenv("NODE_ID")
	sshfsDriver                   *sshfs.Driver
	isWindowsCluster              = os.Getenv(testWindowsEnvVar) != ""
	defaultStorageClassParameters = map[string]string{
		"server": sshfsServerAddress,
		"share":  sshfsShare,
		"csi.storage.k8s.io/provisioner-secret-name":      "mount-options",
		"csi.storage.k8s.io/provisioner-secret-namespace": "default",
		"mountPermissions": "0755",
	}
	storageClassParametersWithZeroMountPermisssions = map[string]string{
		"server": sshfsServerAddress,
		"share":  sshfsShare,
		"csi.storage.k8s.io/provisioner-secret-name":      "mount-options",
		"csi.storage.k8s.io/provisioner-secret-namespace": "default",
		"mountPermissions": "0",
	}
	subDirStorageClassParameters = map[string]string{
		"server": sshfsServerAddress,
		"share":  sshfsShare,
		"subDir": "subDirectory-${pvc.metadata.namespace}",
		"csi.storage.k8s.io/provisioner-secret-name":      "mount-options",
		"csi.storage.k8s.io/provisioner-secret-namespace": "default",
		"mountPermissions": "0755",
	}
	controllerServer *sshfs.ControllerServer
)

type testCmd struct {
	command  string
	args     []string
	startLog string
	endLog   string
}

var _ = ginkgo.BeforeSuite(func() {
	// k8s.io/kubernetes/test/e2e/framework requires env KUBECONFIG to be set
	// it does not fall back to defaults
	if os.Getenv(kubeconfigEnvVar) == "" {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		os.Setenv(kubeconfigEnvVar, kubeconfig)
	}
	handleFlags()
	framework.AfterReadingAllFlags(&framework.TestContext)

	options := sshfs.DriverOptions{
		NodeID:     nodeID,
		DriverName: sshfs.DefaultDriverName,
		Endpoint:   fmt.Sprintf("unix:///tmp/csi-%s.sock", uuid.NewUUID().String()),
	}
	sshfsDriver = sshfs.NewDriver(&options)
	controllerServer = sshfs.NewControllerServer(sshfsDriver)

	// install sshfs server
	installSSHFSServer := testCmd{
		command:  "make",
		args:     []string{"install-sshfs-server"},
		startLog: "Installing SSHFS Server...",
		endLog:   "SSHFS Server successfully installed",
	}

	e2eBootstrap := testCmd{
		command:  "make",
		args:     []string{"e2e-bootstrap"},
		startLog: "Installing SSHFS CSI Driver...",
		endLog:   "SSHFS CSI Driver Installed",
	}
	// todo: Install metrics server once added to this driver

	execTestCmd([]testCmd{installSSHFSServer, e2eBootstrap})
	go func() {
		sshfsDriver.Run(false)
	}()

})

var _ = ginkgo.AfterSuite(func() {
	createExampleDeployment := testCmd{
		command:  "bash",
		args:     []string{"hack/verify-examples.sh"},
		startLog: "create example deployments",
		endLog:   "example deployments created",
	}
	execTestCmd([]testCmd{createExampleDeployment})

	sshfsLog := testCmd{
		command:  "bash",
		args:     []string{"test/utils/sshfs_log.sh"},
		startLog: "===================sshfs log===================",
		endLog:   "==================================================",
	}

	e2eTeardown := testCmd{
		command:  "make",
		args:     []string{"e2e-teardown"},
		startLog: "Uninstalling SSHFS CSI Driver...",
		endLog:   "SSHFS Driver uninstalled",
	}
	execTestCmd([]testCmd{sshfsLog, e2eTeardown})

	// install/uninstall CSI Driver deployment scripts test
	installDriver := testCmd{
		command:  "bash",
		args:     []string{"deploy/install-driver.sh", "master", "local"},
		startLog: "===================install CSI Driver deployment scripts test===================",
		endLog:   "===================================================",
	}
	uninstallDriver := testCmd{
		command:  "bash",
		args:     []string{"deploy/uninstall-driver.sh", "master", "local"},
		startLog: "===================uninstall CSI Driver deployment scripts test===================",
		endLog:   "===================================================",
	}
	execTestCmd([]testCmd{installDriver, uninstallDriver})
})

// handleFlags sets up all flags and parses the command line.
func handleFlags() {
	config.CopyFlags(config.Flags, flag.CommandLine)
	framework.RegisterCommonFlags(flag.CommandLine)
	framework.RegisterClusterFlags(flag.CommandLine)
	flag.Parse()
}

func execTestCmd(cmds []testCmd) {
	err := os.Chdir("../..")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	defer func() {
		err := os.Chdir("test/e2e")
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	}()

	projectRoot, err := os.Getwd()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(strings.HasSuffix(projectRoot, "csi-driver-sshfs")).To(gomega.Equal(true))

	for _, cmd := range cmds {
		log.Println(cmd.startLog)
		cmdSh := exec.Command(cmd.command, cmd.args...)
		cmdSh.Dir = projectRoot
		cmdSh.Stdout = os.Stdout
		cmdSh.Stderr = os.Stderr
		err = cmdSh.Run()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		log.Println(cmd.endLog)
	}
}

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "E2E Suite")
}

func convertToPowershellCommandIfNecessary(command string) string {
	if !isWindowsCluster {
		return command
	}

	switch command {
	case "echo 'hello world' > /mnt/test-1/data && grep 'hello world' /mnt/test-1/data":
		return "echo 'hello world' | Out-File -FilePath C:\\mnt\\test-1\\data.txt; Get-Content C:\\mnt\\test-1\\data.txt | findstr 'hello world'"
	case "touch /mnt/test-1/data":
		return "echo $null >> C:\\mnt\\test-1\\data"
	case "while true; do echo $(date -u) >> /mnt/test-1/data; sleep 100; done":
		return "while (1) { Add-Content -Encoding Unicode C:\\mnt\\test-1\\data.txt $(Get-Date -Format u); sleep 1 }"
	case "echo 'hello world' >> /mnt/test-1/data && while true; do sleep 100; done":
		return "Add-Content -Encoding Unicode C:\\mnt\\test-1\\data.txt 'hello world'; while (1) { sleep 1 }"
	case "echo 'hello world' >> /mnt/test-1/data && while true; do sleep 3600; done":
		return "Add-Content -Encoding Unicode C:\\mnt\\test-1\\data.txt 'hello world'; while (1) { sleep 1 }"
	}

	return command
}
