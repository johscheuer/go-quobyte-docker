package main

import (
	"fmt"
	"log"
	"os"

	"github.com/docker/go-plugins-helpers/volume"
)

const quobyteID = "quobyte"

// Mandatory configuration
var qmgmtUser string
var qmgmtPassword string
var quobyteAPIURL string
var quobyteRegistry string

// Optional configuration
var mountQuobytePath string
var mountQuobyteOptions string
var qmgmtPath string
var defaultVolumeConfiguration string

// Constants
const pluginDirectory string = "/run/docker/plugins/"
const pluginSocket string = "/run/docker/plugins/quobyte.sock"
const mountDirectory string = "/run/docker/quobyte/mnt"

func main() {
	readMandatoryConfig()
	readOptionalConfig()

	if err := os.MkdirAll(mountDirectory, 0555); err != nil {
		log.Println(err.Error())
	}

	if err := os.MkdirAll(pluginDirectory, 055); err != nil {
		log.Println(err.Error())
	}

	if !isMounted(mountDirectory) {
		log.Printf("Mounting Quobyte namespace in %s", mountDirectory)
		mountAll()
	}

	qDriver := newQuobyteDriver(quobyteAPIURL, qmgmtUser, qmgmtPassword)
	fmt.Println(volume.NewHandler(qDriver).ServeUnix("root", quobyteID))
}
