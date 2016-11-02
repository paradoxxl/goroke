package main


import (
	"github.com/paradoxxl/goroke/audiorecorder"
	"flag"
	"fmt"
	"github.com/paradoxxl/goroke/orchestrator"
	"net"
)

var ListAudioDevices  = flag.Bool("l",false,"list audio devices for recording")
var RecordingDevice = flag.String("d","","Audio device which will be used for recording. Use -l for displaying suitable devices. If the devicename has spaces, it needs to be in quotes \"\"")

func main(){
	flag.Parse()

	if *ListAudioDevices{
		rec := audiorecorder.RecorderController{}
		rec.ListDevices()
	}else{
		if *RecordingDevice == ""{
			fmt.Println("You did not set an recording device, hence recording will be unsupported")
		}
			orchestrator := orchestrator.Orchestrator{}
			orchestrator.Setup(net.ParseIP("127.0.0.1"),57570)

			orchestrator.Start()

			cstop := make(chan interface{})
			<- cstop

	}


}
