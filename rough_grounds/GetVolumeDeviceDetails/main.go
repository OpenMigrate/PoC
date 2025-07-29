package main

import(
	"fmt"
	"os"
	"strings"
)

func main(){
	files,err := os.ReadDir("/sys/block")
	if(err!=nil){
		fmt.Println(err)
	}
	for _,file:= range files{
		VolumeDeviceName := file.Name()
		if(!strings.Contains(VolumeDeviceName,"loop")){
			fmt.Println(VolumeDeviceName)
		}
	}	
}