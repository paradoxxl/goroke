package audiorecorder

import (
	"time"
	"os/exec"
	"fmt"
	"os"
)

type Recording struct{
	Singer string
	SongTitle string
	Date time.Time

	Cmd exec.Cmd
	Filename string

	sigkill chan interface{}

}

func (self *Recording) Setup(Singer string,SongTitle string){
	self.Singer = Singer
	self.SongTitle = SongTitle
	self.Date = time.Now()

	y,m,d:=self.Date.Date()
	self.Filename = fmt.Sprintf("%v-%v-%v_%v_%v.%v",self.Singer,self.SongTitle,y,m,d,FileEnding)

	self.Cmd = exec.Command(FFmpegPath, "-f","dshow", "-i","audio=" + RecordingDevice, self.Filename)
	//cmd := exec.Command("C:/Dev/ffmpeg/bin/ffmpeg.exe")
	self.Cmd.Stdout = os.Stdout
	self.Cmd.Stderr = os.Stderr
	self.Cmd.Stdin = os.Stdin

	self.sigkill = make(chan interface{})
}

func (self *Recording) Start(){
	go self.record()
}

func (self *Recording) Stop(){
	self.sigkill<-true
}

func (self *Recording) record(){
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
