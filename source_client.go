package main
import (
	"fmt"
	"io"
	"os"
	"time"
	"github.com/u-root/u-root/pkg/mount/block"
	"net"
)


const (
	BlockSize = 1024 * 1024
)

type Device struct{
	Details *block.BlockDev
	MountPoint *string
	DeviceBlockSize uint64 //Size of the block device in bytes
	DeviceName string
}

type DeviceReplicator struct {
	DeviceDetails Device
	Blocksize int
}


func (device *Device) GetDetails() error{
	// Get device details
	// Check for the device
	var err error
	device.Details, err = block.Device(device.DeviceName)
	if err != nil{
		fmt.Println("Device not found")
		return err
	}
	fmt.Println("Device Found Details: ", device.Details)
	
	// Get Device Path
	devicePath := device.Details.DevicePath()
	fmt.Println("Device Path: ", devicePath)

	// Get Device Mount Point
	device.MountPoint, err= block.GetMountpointByDevice(devicePath)
	if(err != nil){
		fmt.Println("Mount point ot found for device: ",device.DeviceName)
		return err
	}
	fmt.Println("Device Mount : ", *device.MountPoint)

	device.DeviceBlockSize, err = device.Details.Size()
	if(err!=nil){
		fmt.Println("Error in finding out size: ", err)
		return err
	}
	fmt.Println("Device Size in bytes: ", device.DeviceBlockSize)
	return nil
}

func (dr *DeviceReplicator) ReplicateToDestination() error{
	//Open Source Device for reading
	source, err := os.OpenFile(dr.DeviceDetails.DeviceName,os.O_RDONLY,0)
	if (err != nil){
		return err
	}
	defer source.Close()


	/// Connect to server
	conn,err := net.Dial("tcp","localhost:4000")
	if err != nil{
		return err
	}
	defer conn.Close()

	fmt.Printf("Connected to the server")

	fmt.Printf("Replicating Device %s (%d bytes) ...\n",dr.DeviceDetails.Details,dr.DeviceDetails.DeviceBlockSize)


	/// Read the source and write to network
	clientWriter := &ClientWriter{
		Writer: conn,
		Total: dr.DeviceDetails.DeviceBlockSize,
		StartTime: time.Now(),
		BlockSize: dr.Blocksize,
	}

	buffer := make([]byte,dr.Blocksize)

	fmt.Println("Starting Replication")
	
	// start writing to network
	copied, err := io.CopyBuffer(clientWriter,source,buffer)

	if err != nil {
		fmt.Printf("Error during replication: %v\n", err)
		return err
	}

	fmt.Println("Replication Completed")
	fmt.Printf("Total copied: %d bytes (%.2f GB)\n", copied, float64(copied)/(1024*1024*1024))
	
	return nil
}

// ClientWriter to write the data from soirce to netweork
type ClientWriter struct{
	Writer     io.Writer
	Total      uint64
	Current    uint64
	StartTime  time.Time
	lastReport uint64
	BlockSize  uint64
}


// Client Writer that will write to network
func (cw *ClientWriter) Write(p []byte) (int,error){
	n,err := cw.Writer.Write(p)
	cw.Current += uint64(n)

	reportInterval := uint64(100*1024*1024)
	blockInterval := cw.BlockSize * 1000

	if(blockInterval<reportInterval){
		reportInterval=blockInterval
	}

	if(cw.Current-cw.lastReport>=reportInterval || err==io.EOF){
		percentage := float64(cw.Current)/float64(cw.Total) * 100
		elapsed := time.Since(cw.StartTime)
		speed := float64(cw.Current)/elapsed.Seconds() / (1024*1024) //MB/S
		blocksProcessed := cw.Current/cw.BlockSize

		fmt.Printf("\rProgress: %.1f%% (%d/%d bytes, %d blocks) - Speed: %.1f MB/s", 
			percentage, cw.Current, cw.Total, blocksProcessed, speed)
		cw.lastReport=cw.Current
	}
	return n,err
}

func main(){
	// Get a device details
	//Source device
	device := Device{
		DeviceName: "/dev/nvme0n1",
	}
	
	if err:=device.GetDetails(); err!=nil{
		fmt.Println("Error",err)
	}

	dr := DeviceReplicator{
		DeviceDetails: device,
		BlockSize: 1024,
	}
	if err:=dr.ReplicateToDestination(); err!=nil{
		fmt.Println("Replication Function Error: ",err)
	}
}