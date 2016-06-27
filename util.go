package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func readOptionalConfig() {
	mountQuobytePath = os.Getenv("MOUNT_QUOBYTE_PATH")
	qmgmtPath = os.Getenv("QMGMT_PATH")
	mountQuobyteOptions = os.Getenv("MOUNT_QUOBYTE_OPTIONS")
	if len(mountQuobyteOptions) == 0 {
		mountQuobyteOptions = "-o user_xattr"
	}
	defaultVolumeConfiguration = os.Getenv("DEFAULT_VOLUME_CONFIGURATION")
	if len(defaultVolumeConfiguration) == 0 {
		defaultVolumeConfiguration = "BASE"
	}
}

func readMandatoryConfig() {
	qmgmtUser = getMandatoryEnv("QUOBYTE_API_USER")
	qmgmtPassword = getMandatoryEnv("QUOBYTE_API_PASSWORD")
	quobyteAPIURL = getMandatoryEnv("QUOBYTE_API_URL")
	// host[:port][,host:port] or SRV record name
	quobyteRegistry = getMandatoryEnv("QUOBYTE_REGISTRY")
}

func getMandatoryEnv(name string) string {
	env := os.Getenv(name)
	if len(env) < 0 {
		log.Fatalf("Please set %s in environment", name)
	}

	return env
}

func isMounted(mountPath string) bool {
	content, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		log.Println(err)
	}
	for _, mount := range strings.Split(string(content), "\n") {
		if strings.Split(mount, " ")[1] == mountPath {
			return true
		}
	}

	return false
}

func mountAll() {
	binary := path.Join(mountQuobytePath, "mount.quobyte")
	if err := exec.Command(binary, mountQuobyteOptions, fmt.Sprintf("%s/", quobyteRegistry), mountDirectory).Run(); err != nil {
		log.Fatal(err)
	}
}
