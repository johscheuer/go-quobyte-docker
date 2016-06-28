package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/go-plugins-helpers/volume"
	quobyte_api "github.com/quobyte/api"
)

type quobyteDriver struct {
	client *quobyte_api.QuobyteClient
}

func newQuobyteDriver(apiURL string, username string, password string) quobyteDriver {
	driver := quobyteDriver{
		client: quobyte_api.NewQuobyteClient(apiURL, username, password),
	}

	return driver
}

func (driver quobyteDriver) Create(request volume.Request) volume.Response {
	log.Printf("Creating volume %s\n", request.Name)
	//TODO how to get user and group -> request.Opts[user] + request.Opts[group] if null set root?

	if _, err := driver.client.CreateVolume(request.Name, "root", "root"); err != nil {
		return volume.Response{Err: err.Error()}
	}

	return volume.Response{Err: ""}
}

func (driver quobyteDriver) Remove(request volume.Request) volume.Response {
	log.Printf("Removing volume %s\n", request.Name)

	if err := driver.client.DeleteVolumeByName(request.Name); err != nil {
		return volume.Response{Err: err.Error()}
	}

	return volume.Response{Err: ""}
}

func (driver quobyteDriver) Mount(request volume.Request) volume.Response {
	mPoint := driver.mountpoint(request.Name)
	log.Printf("Mounting volume %s on %s\n", request.Name, mPoint)
	if fi, err := os.Lstat(mPoint); err != nil || !fi.IsDir() {
		return volume.Response{Err: fmt.Sprintf("%v not mounted", mPoint)}
	}

	return volume.Response{Err: "", Mountpoint: mPoint}
}

func (driver quobyteDriver) Path(request volume.Request) volume.Response {
	return volume.Response{Mountpoint: driver.mountpoint(request.Name)}
}

func (driver quobyteDriver) Unmount(request volume.Request) volume.Response {
	return volume.Response{}
}

func (driver quobyteDriver) Get(request volume.Request) volume.Response {
	mPoint := driver.mountpoint(request.Name)

	if fi, err := os.Lstat(mPoint); err != nil || !fi.IsDir() {
		return volume.Response{Err: fmt.Sprintf("%v not mounted", mPoint)}
	}

	return volume.Response{Volume: &volume.Volume{Name: request.Name, Mountpoint: mPoint}}
}

func (driver quobyteDriver) List(request volume.Request) volume.Response {
	var vols []*volume.Volume
	files, err := ioutil.ReadDir(mountQuobytePath)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	for _, entry := range files {
		if entry.IsDir() {
			vols = append(vols, &volume.Volume{Name: entry.Name(), Mountpoint: driver.mountpoint(entry.Name())})
		}
	}

	return volume.Response{Volumes: vols}
}

func (driver quobyteDriver) Capabilities(request volume.Request) volume.Response {
	return volume.Response{Capabilities: volume.Capability{Scope: "local"}}
}

func (driver *quobyteDriver) mountpoint(name string) string {
	return filepath.Join(mountQuobytePath, name)
}
