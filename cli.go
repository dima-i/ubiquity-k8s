package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/IBM/ubiquity-k8s/controller"
	flags "github.com/jessevdk/go-flags"

	"github.com/IBM/ubiquity/resources"
	"github.com/IBM/ubiquity/utils"
)

var configFile = flag.String(
	"configFile",
	"/tmp/ubiquity-client.conf",
	"config file with ubiquity client configuration params",
)

type InitCommand struct {
	Init func() `short:"i" long:"init" description:"Initialize the plugin"`
}

func (i *InitCommand) Execute(args []string) error {
	controller, err := createController(*configFile)
	if err != nil {
		response := resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed tocreate controller %#v", err),
			Device:  "",
		}
		return utils.PrintResponse(response)
	}
	response := controller.Init()
	return utils.PrintResponse(response)
}

type GetVolumeNameCommand struct {
	GetVolumeName func() `short:"g" long:"getvolumename" description:"Get a cluster wide unique volume name for the volume"`
}

func (g *GetVolumeNameCommand) Execute(args []string) error {
	getVolumeNameRequest := make(map[string]string)
	err := json.Unmarshal([]byte(args[0]), &getVolumeNameRequest)
	if err != nil {
		response := resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed to unmarshall request %#v", err),
		}
		return utils.PrintResponse(response)
	}
	controller, err := createController(*configFile)
	if err != nil {
		response := resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed to create controller %#v", err),
		}
		return utils.PrintResponse(response)
	}
	response := controller.GetVolumeName(getVolumeNameRequest)
	return utils.PrintResponse(response)
}

type AttachCommand struct {
	Attach func() `short:"a" long:"attach" description:"Attach a volume"`
}

func (a *AttachCommand) Execute(args []string) error {
	attachRequest := make(map[string]string)
	err := json.Unmarshal([]byte(args[0]), &attachRequest)
	if err != nil {
		response := resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed to attach volume %#v", err),
			Device:  "",
		}
		return utils.PrintResponse(response)
	}
	attachRequest["nodeName"] = args[1]
	controller, err := createController(*configFile)

	if err != nil {
		response := resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed to create controller in attach %#v", err),
			Device:  "",
		}
		return utils.PrintResponse(response)
	}
	attachResponse := controller.Attach(attachRequest)
	return utils.PrintResponse(attachResponse)
}

type DetachCommand struct {
	Detach func() `short:"d" long:"detach" description:"Detach a volume"`
}

func (d *DetachCommand) Execute(args []string) error {
	mountDevice := args[0]
	node := args[1]
	controller, err := createController(*configFile)

	if err != nil {
		response := resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed to create controller in dettach %#v", err),
			Device:  "",
		}
		return utils.PrintResponse(response)
	}

	detachRequest := resources.FlexVolumeDetachRequest{Name: mountDevice, Node: node}
	detachResponse := controller.Detach(detachRequest)
	return utils.PrintResponse(detachResponse)
}

type WaitForAttachCommand struct {
	WaitForAttach func() `short:"w" long:"waitforattach" description:"waits for a volume to get attached"`
}

func (w *WaitForAttachCommand) Execute(args []string) error {
	waitForAttachRequest := make(map[string]string)
	err := json.Unmarshal([]byte(args[1]), &waitForAttachRequest)
	if err != nil {
		response := resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed to unmarshall wait for attach request%#v", err),
			Device:  "",
		}
		return utils.PrintResponse(response)
	}
	waitForAttachRequest["device"] = args[0]
	controller, err := createController(*configFile)

	if err != nil {
		response := resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed to create controller for wait for attach %#v", err),
		}
		return utils.PrintResponse(response)
	}
	waitForAttachResponse := controller.WaitForAttach(waitForAttachRequest)

	return utils.PrintResponse(waitForAttachResponse)
}

type IsAttachedCommand struct {
	Attach func() `short:"i" long:"isattached" description:"Is volume attached"`
}

func (i *IsAttachedCommand) Execute(args []string) error {
	isAttachRequest := make(map[string]string)
	err := json.Unmarshal([]byte(args[0]), &isAttachRequest)
	if err != nil {
		response := resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed to unmarshall isAttached request %#v", err),
			Device:  "",
		}
		return utils.PrintResponse(response)
	}
	isAttachRequest["volumeName"] = args[1]
	controller, err := createController(*configFile)

	if err != nil {
		panic(fmt.Sprintf("backend %s not found", configFile))
	}
	isAttachResponse := controller.IsAttached(isAttachRequest)

	return utils.PrintResponse(isAttachResponse)
}

type MountDeviceCommand struct {
	Mount func() `short:"" long:"mountdevice" description:"Mount a device "`
}

func (m *MountDeviceCommand) Execute(args []string) error {
	targetMountDir := args[0]
	mountDevice := args[1]
	var mountOpts map[string]interface{}

	err := json.Unmarshal([]byte(args[2]), &mountOpts)
	if err != nil {
		mountResponse := resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed to mount device %s to %s due to: %#v", mountDevice, targetMountDir, err),
			Device:  mountDevice,
		}
		return utils.PrintResponse(mountResponse)
	}

	mountRequest := resources.FlexVolumeMountRequest{
		MountPath:   targetMountDir,
		MountDevice: mountDevice,
		Opts:        mountOpts,
	}

	controller, err := createController(*configFile)

	if err != nil {
		panic("backend not found")
	}
	mountResponse := controller.Mount(mountRequest)
	return utils.PrintResponse(mountResponse)
}

type UnmountDeviceCommand struct {
	Mount func() `short:"umd" long:"unmountdevice" description:"Unmount a device"`
}

func (m *UnmountDeviceCommand) Execute(args []string) error {
	targetMountDir := args[0]
	mountDevice := args[1]
	var mountOpts map[string]interface{}

	err := json.Unmarshal([]byte(args[2]), &mountOpts)
	if err != nil {
		mountResponse := resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed to mount device %s to %s due to: %#v", mountDevice, targetMountDir, err),
			Device:  mountDevice,
		}
		return utils.PrintResponse(mountResponse)
	}

	mountRequest := resources.FlexVolumeMountRequest{
		MountPath:   targetMountDir,
		MountDevice: mountDevice,
		Opts:        mountOpts,
	}

	controller, err := createController(*configFile)

	if err != nil {
		panic("backend not found")
	}
	mountResponse := controller.Mount(mountRequest)
	return utils.PrintResponse(mountResponse)
}

type MountCommand struct {
	Mount func() `short:"m" long:"mount" description:"Mount a volume Id to a path"`
}

func (m *MountCommand) Execute(args []string) error {
	targetMountDir := args[0]
	mountDevice := args[1]
	var mountOpts map[string]interface{}

	err := json.Unmarshal([]byte(args[2]), &mountOpts)
	if err != nil {
		mountResponse := resources.FlexVolumeResponse{
			Status:  "Failure",
			Message: fmt.Sprintf("Failed to mount device %s to %s due to: %#v", mountDevice, targetMountDir, err),
			Device:  mountDevice,
		}
		return utils.PrintResponse(mountResponse)
	}

	mountRequest := resources.FlexVolumeMountRequest{
		MountPath:   targetMountDir,
		MountDevice: mountDevice,
		Opts:        mountOpts,
	}

	controller, err := createController(*configFile)

	if err != nil {
		panic("backend not found")
	}
	mountResponse := controller.Mount(mountRequest)
	return utils.PrintResponse(mountResponse)
}

type UnmountCommand struct {
	UnMount func() `short:"u" long:"unmount" description:"UnMount a volume Id to a path"`
}

func (u *UnmountCommand) Execute(args []string) error {
	mountDir := args[0]
	controller, err := createController(*configFile)

	if err != nil {
		panic("backend not found")
	}

	unmountRequest := resources.FlexVolumeUnmountRequest{
		MountPath: mountDir,
	}
	unmountResponse := controller.Unmount(unmountRequest)
	return utils.PrintResponse(unmountResponse)
}

type Options struct{}

func main() {
	var mountCommand MountCommand
	var unmountCommand UnmountCommand
	var attachCommand AttachCommand
	var detachCommand DetachCommand
	var initCommand InitCommand
	var getVolumeNameCommand GetVolumeNameCommand
	var isAttachedCommand IsAttachedCommand
	var waitForAttachCommand WaitForAttachCommand

	var options Options
	var parser = flags.NewParser(&options, flags.Default)

	parser.AddCommand("init",
		"Init the plugin",
		"The info command print the driver name and version.",
		&initCommand)

	parser.AddCommand("getvolumename",
		"GetVolumeName",
		"Get a cluster wide unique volume name for the volume.",
		&getVolumeNameCommand)
	parser.AddCommand("isattached",
		"IsAttached",
		"Checks if volume is attached.",
		&isAttachedCommand)
	parser.AddCommand("waitforattach",
		"Waits for attach",
		"Wait for a volume to get attached.",
		&waitForAttachCommand)
	parser.AddCommand("mount",
		"Mount Volume",
		"Mount a volume Id to a path - returning the path.",
		&mountCommand)
	parser.AddCommand("unmount",
		"Unmount Volume",
		"UnMount given a mount dir",
		&unmountCommand)
	parser.AddCommand("attach",
		"Attach Volume",
		"Attach Volume",
		&attachCommand)
	parser.AddCommand("detach",
		"Detach Volume",
		"Detach a Volume",
		&detachCommand)

	_, err := parser.Parse()
	if err != nil {

		logger, _ := setupLogger("/tmp")

		logger.Printf("Error parsing %#v", err)
		os.Exit(1)
	}
}

func createController(configFile string) (*controller.Controller, error) {
	var config resources.UbiquityPluginConfig
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		fmt.Printf("error decoding config file", err)
		return nil, err

	}

	logger, _ := setupLogger(config.LogPath)
	//defer closeLogs(logFile)

	storageApiURL := fmt.Sprintf("http://%s:%d/ubiquity_storage", config.UbiquityServer.Address, config.UbiquityServer.Port)
	controller, err := controller.NewController(logger, storageApiURL, config)
	return controller, err
}

func setupLogger(logPath string) (*log.Logger, *os.File) {
	logFile, err := os.OpenFile(path.Join(logPath, "ubiquity-flexvolume.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		fmt.Printf("Failed to setup logger: %s\n", err.Error())
		return nil, nil
	}
	log.SetOutput(logFile)
	logger := log.New(io.MultiWriter(logFile), "ubiquity-flexvolume: ", log.Lshortfile|log.LstdFlags)
	return logger, logFile
}
