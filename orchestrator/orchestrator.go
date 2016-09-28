package orchestrator

import (
	"github.com/paradoxxl/goroke/network"
	"github.com/paradoxxl/goroke/hid"
	"net"
	"github.com/paradoxxl/goroke/audiorecorder"
)

type Orchestrator struct{
	NetworkController *network.NetworkController
	WarriorController *HID.WarriorController
	Recorder *audiorecorder.RecorderController
}

func (self *Orchestrator  )Setup(KaraFunIP net.IP, KaraFunPort uint16){
	self.WarriorController = HID.NewWarrior()
	self.NetworkController = network.NewNetworkController(KaraFunIP,KaraFunPort)
	self.Recorder = &audiorecorder.RecorderController{}

	self.Recorder.NetworkController = self
	self.WarriorController.Audiorecorder = self
	self.WarriorController.NetworkController = self
	self.NetworkController.Recorder = self
	self.NetworkController.Connect()
}

func (self *Orchestrator) Start(){
	self.WarriorController.Connect()
	self.WarriorController.StartPolling()
}


//Implementing Network Interface for redirecting requests
func (self *Orchestrator)GeneralVolumeSliderChanged(pos byte){
	self.NetworkController.GeneralVolumeSliderChanged(pos)
}
func (self *Orchestrator)VoiceVolumeSliderPosChanged(pos byte){
	self.NetworkController.VoiceVolumeSliderPosChanged(pos)
}
func (self *Orchestrator)MaleVolumeSliderPosChanged(pos byte){
	self.NetworkController.MaleVolumeSliderPosChanged(pos)
}
func (self *Orchestrator)FemaleVolumeSliderPosChanged(pos byte){
	self.NetworkController.FemaleVolumeSliderPosChanged(pos)
}
func (self *Orchestrator)PitchChanged(val int8){
	self.NetworkController.PitchChanged(val)
}
func (self *Orchestrator)TempoChanged(val int8){
	self.NetworkController.TempoChanged(val)
}
func (self *Orchestrator)PlayPressed(){
	self.NetworkController.PlayPressed()
}
func (self *Orchestrator)PausePressed(){
	self.NetworkController.PausePressed()
}
func (self *Orchestrator)NextPressed(){
	self.NetworkController.NextPressed()
}
func (self *Orchestrator)PrevPressed(){
	self.NetworkController.PrevPressed()
}
func (self *Orchestrator)GetStatus(){
	self.NetworkController.GetStatus()
}


//Implementing Recorder Interface for redirecting requests
func (self *Orchestrator) RecordButtonPressed(){
	self.Recorder.RecordButtonPressed()
}
func (self *Orchestrator) XMLInput(in string){
	self.Recorder.XMLInput(in)
}