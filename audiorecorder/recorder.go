package audiorecorder

import (
	"time"
	"os/exec"
	"fmt"
	"os"
	"github.com/paradoxxl/goroke/xmlparser"
	"github.com/paradoxxl/goroke/Interfaces"
)

type RecorderController struct{
	Singer string
	SongTitle string
	Date time.Time
	Recording bool

	Cmd *exec.Cmd
	Filename string

	sigkill chan interface{}

	NetworkController Interfaces.INetworkController

}


func (self *RecorderController) RecordButtonPressed(){
	if self.Recording{
		self.Stop()
	}else{
		self.Start()
	}
}
func (self *RecorderController) XMLInput(input string){
	parsed,err := xmlparser.ParseXML(input)
	if err != nil{
		fmt.Println(err)
		return
	}
	singer := parsed.Queue.Items[0].Singer
	title := parsed.Queue.Items[0].Title

	if self.Recording && parsed.State != "playing"{
		self.Stop()
	}

	self.Setup(singer,title)
}

func (self *RecorderController) Setup(Singer string,SongTitle string){
	self.Singer = Singer
	self.SongTitle = SongTitle
	self.Date = time.Now()

	y,m,d:=self.Date.Date()
	self.Filename = fmt.Sprintf("%v-%v-%v_%v_%v.%v",self.Singer,self.SongTitle,y,m,d,FileEnding)

	self.Cmd = exec.Command(FFmpegPath, "-f","dshow", "-i","audio=" + RecordingDevice, self.Filename)
	self.Cmd.Stdout = os.Stdout
	self.Cmd.Stderr = os.Stderr
	self.Cmd.Stdin = os.Stdin

	self.sigkill = make(chan interface{})
}

func (self *RecorderController) Start(){
	if self.SongTitle != "" {
		self.NetworkController.GetStatus()
	}
	self.Recording = true
	go self.record()
}

func (self *RecorderController) Stop(){
	if self.Recording {
		self.Recording = false
		self.sigkill<-true
	}
}

func (self *RecorderController) record(){
	err := self.Cmd.Start()
	if err != nil {
		println(err.Error())
	}


	<-self.sigkill
	self.Cmd.Process.Kill()

	err = self.Cmd.Wait()
	if err != nil {
		println(err.Error())
	}
}
