package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

const (
	BlockSize = 1024 * 1024
)

type VolumeDevice struct{
	DeviceName string // name of the device in the system
	SizeInBlock uint64 // size of the volume volume device in blocks
	LogicalBlockSize int //Logical block size of the volume device
}
// Volume Interface
type Volume interface{
	GetVolumeDeviceDetails() error //Get Details of volume device
}


func (this *VolumeDevice) GetVolumeDeviceDetails () (error){
	//Get volume device details
	// Get logical block size of device
	logicalBlockSize,err := os.ReadFile("/sys/block/"+this.DeviceName+"/queue/logical_block_size")
	if(err!=nil){
		fmt.Println(err)
		return err
	}
	value := strings.TrimSpace(string(logicalBlockSize))
	this.LogicalBlockSize, err = strconv.Atoi(value)
	if(err!=nil){
		fmt.Println(err)
		return err
	}

	VolumeSize,err := os.ReadFile("/sys/block/"+this.DeviceName+"/size")
	if(err!=nil){
		fmt.Println(err)
		return err
	}
	value = strings.TrimSpace(string(VolumeSize))
	sizeInBlock, err := strconv.Atoi(value)
	this.SizeInBlock = uint64(sizeInBlock)

	return nil
}

func NewVolumeDevice(name string) *VolumeDevice{
	return &VolumeDevice{
		DeviceName:name,
	}
}

//Agent Interface
type AgentTask interface{
	FindVolumeDevices() error //Scan the machine for volume devices
}

func (this *Agent) FindVolumeDevices () (error){
	//get the list of all volume devices mounted to the machine
	files,err := os.ReadDir("/sys/block")
	if(err!=nil){
		fmt.Println(err) 
		return err
	}
	var devices []VolumeDevice

	for _,file := range files{
		VolumeDeviceName := file.Name()
		if(!strings.Contains(VolumeDeviceName,"loop")){
			device := NewVolumeDevice(VolumeDeviceName)
			device.GetVolumeDeviceDetails()
			devices = append(devices,*device)
		}
	}
	this.Devices = devices
	return nil
}

//creates a new agent instance
func NewAgent(){
	// Create new agent
}

type Agent struct{
	MasterServer string //location of master server i.e IP address/server url with port
	ReplicationNode string //location of replication server i.e IP address/server url with port
	Devices []VolumeDevice // list volume devices available in the machine
}

func main(){
	agent := Agent{}
	agent.FindVolumeDevices()
	fmt.Println(agent.Devices)
}