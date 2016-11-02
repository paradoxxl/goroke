package network

import (
	"net"
	"golang.org/x/net/websocket"
	"log"
	"fmt"
	"github.com/paradoxxl/goroke/Interfaces"
	"time"
	"sync"
)

type SlidersState struct{
	sync.RWMutex
	Data map[string]*SliderState
}
type SliderState struct{
	Position byte
	Changed bool
}
func NewSlidersState() *SlidersState{
	return & SlidersState{
		Data:make(map[string]*SliderState),
	}
}

type NetworkController struct{
	IP net.IP
	Port uint16
	Connected bool

	Ws *websocket.Conn
	sigkill chan interface{}

	Recorder Interfaces.IRecorderController

	SlidersState *SlidersState

}

func (self *NetworkController) SlidersUpdate(){
	for{

			//Check all sliders and update when necessary

			self.SlidersState.Lock()
			defer self.SlidersState.Unlock()

			if sl,ok:=self.SlidersState.Data[GeneralVolumeSlider];ok {
				self.sendString(fmt.Sprintf(GeneralVolumeCommand,sl.Position))
				sl.Changed = false
			}
			if sl,ok:=self.SlidersState.Data[VoiceVolumeSlider];ok {
				self.sendString(fmt.Sprintf(VoiceVolumeCommand,sl.Position))
				sl.Changed = false
			}
			if sl,ok:=self.SlidersState.Data[MaleVolumeSlider];ok {
				self.sendString(fmt.Sprintf(MaleVolumeCommand,sl.Position))
				sl.Changed = false
			}
			if sl,ok:=self.SlidersState.Data[FemaleVolumeSlider];ok {
				self.sendString(fmt.Sprintf(FemaleVolumeCommand,sl.Position))
				sl.Changed = false

			}

			time.Sleep(SlidersUpdateInterval)
		}

}

func NewNetworkController(KaraFunIP net.IP, KaraFunPort uint16) *NetworkController{
	return &NetworkController{
		IP: KaraFunIP,
		Port:KaraFunPort,
		Connected:false,
		Ws:nil,
		sigkill:make(chan interface{}),
		SlidersState:NewSlidersState(),
	}
}


func (self *NetworkController) GeneralVolumeSliderChanged(pos byte){
	self.SlidersState.Lock()
	defer self.SlidersState.Unlock()

	if self.SlidersState != nil{
		sl,ok:=self.SlidersState.Data[GeneralVolumeSlider];
		if ok{
			if sl.Position != pos {
				sl.Position = pos
				sl.Changed = true
			}
		}else{
			self.SlidersState.Data[GeneralVolumeSlider] = &SliderState{Position:pos,Changed:true}
		}
	}

}
func (self *NetworkController) VoiceVolumeSliderPosChanged(pos byte){
	self.SlidersState.Lock()
	defer self.SlidersState.Unlock()

	if self.SlidersState != nil{
		sl,ok:=self.SlidersState.Data[VoiceVolumeSlider];
		if ok{
			if sl.Position != pos {
				sl.Position = pos
				sl.Changed = true
			}
		}else{
			self.SlidersState.Data[VoiceVolumeSlider] = &SliderState{Position:pos,Changed:true}
		}
	}
}
func (self *NetworkController) MaleVolumeSliderPosChanged(pos byte){
	self.SlidersState.Lock()
	defer self.SlidersState.Unlock()

	if self.SlidersState != nil{
		sl,ok:=self.SlidersState.Data[MaleVolumeSlider];
		if ok{
			if sl.Position != pos {
				sl.Position = pos
				sl.Changed = true
			}
		}else{
			self.SlidersState.Data[MaleVolumeSlider] = &SliderState{Position:pos,Changed:true}
		}
	}
}
func (self *NetworkController) FemaleVolumeSliderPosChanged(pos byte){
	self.SlidersState.Lock()
	defer self.SlidersState.Unlock()

	if self.SlidersState != nil{
		sl,ok:=self.SlidersState.Data[FemaleVolumeSlider];
		if ok{
			if sl.Position != pos {
				sl.Position = pos
				sl.Changed = true
			}
		}else{
			self.SlidersState.Data[FemaleVolumeSlider] = &SliderState{Position:pos,Changed:true}
		}
	}
}
func (self *NetworkController) PitchChanged(pitch int8){
	self.sendString(fmt.Sprintf(PitchCommand,pitch))
}
func (self *NetworkController) TempoChanged(tempo int8){
	self.sendString(fmt.Sprintf(TempoCommand, tempo))
}

func (self *NetworkController) PlayPressed(){
	self.sendString(PlayCommand)
}
func (self *NetworkController) PausePressed(){
	self.sendString(PauseCommand)
}
func (self *NetworkController) NextPressed(){
	self.sendString(NextCommand)
}
func (self *NetworkController) PrevPressed(){
	self.sendString(PrevCommand)
}


func (self *NetworkController) GetStatus(){
	self.sendString(GetStatusCommand)
}

func (self *NetworkController) sendString(s string){
	if self.Connected && self.Ws != nil{
		if err:= websocket.Message.Send(self.Ws,s);err!=nil{
			log.Printf("cannot send: %v", err)
		}
	}
}

func (self *NetworkController) Connect() bool{
	if ip:=self.IP.To4();ip!= nil{
		dst:= fmt.Sprintf(url,self.IP.String(),self.Port)
		ws, err := websocket.Dial(dst, "", origin)
		if err != nil{
			self.Connected = false
			log.Printf("Cannot connect to the server: %v",dst)
			return false
		}

		self.Ws = ws
		log.Printf("Connected to the server: %v",dst)
		self.Connected = true
		go self.receiveMsg()
	}
	return true
}


func (self *NetworkController) receiveMsg(){
	var msg string
	for{
		select{
		case <- self.sigkill:
			return
		default:
			if err := websocket.Message.Receive(self.Ws, &msg); err != nil {
				log.Printf("receive error %v", err)
				return
			}
			//fmt.Printf("Receive: %s\n", msg)
		self.Recorder.XMLInput(msg)
		}
	}
}
