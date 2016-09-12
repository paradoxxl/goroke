package HID

import (
	"github.com/paradoxxl/hid"
	"log"
	"github.com/paradoxxl/goroke/network"
)

type WarriorController struct {
	Connected     bool
	DevicePresent bool

	State *WarriorControllerState

	winfo *hid.DeviceInfo
	wdev  hid.Device

	pollStop chan interface{}

	NetworkController *network.NetworkController
}

type WarriorControllerState struct {
	GeneralVolumeSliderPos byte
	VoiceVolumeSliderPos   byte
	MaleVolumeSliderPos    byte
	FemaleVolumeSliderPos  byte

	PlayPauseBtnPressed bool
	NextBtnPressed      bool
	PrevBtnPressed      bool
	KeyDownBtnPressed   bool
	KeyUpBtnPressed     bool
	PitchDownBtnPressed bool
	PitchUpBtnPressed   bool
	RecordBtnPressed    bool

	PlayPauseState	bool
}

func scaleSliderPos(pos byte) byte {
	return byte((float32(pos) / MaxSliderPosHW) * MaxSliderPosSW)
}

func NewWarrior() *WarriorController {
	return &WarriorController{
		Connected: false,
		State:     &WarriorControllerState{},
	}
}

func (self *WarriorController) findWarrior() bool {
	wchan := hid.FindDevices(vendorID, productID)
	self.winfo = <-wchan

	if self.winfo == nil {
		self.DevicePresent = false
		log.Println("Warrior not plugged in")
	}
	log.Printf("Warrior found: PID: %#04X \t VID: %#04X\n", self.winfo.ProductId, self.winfo.VendorId)

	self.DevicePresent = true
	return self.DevicePresent
}

func (self *WarriorController) Connect() {
	if !self.findWarrior() {
		// Device not present!
		return
	}

	wdev, err := self.winfo.Open()
	if err != nil {
		log.Printf("Cannot connect to warrior. Msg: \t%v", err)
		return
	}
	self.wdev = wdev
	log.Println("Connected to warrior")
}

func (self *WarriorController) StartPolling() {
	go self.poll()
}

func (self *WarriorController) StopPolling() {
	self.pollStop <- true
}

func (self *WarriorController) detectChanges(input []byte) {
	// Check Volume Sliders
	if pos:= scaleSliderPos(input[GeneralVolumeSliderBytePos]); self.State.GeneralVolumeSliderPos != pos {
		self.State.GeneralVolumeSliderPos = pos
		log.Printf("General Volume:\t%v",pos)

		self.NetworkController.GeneralVolumeSliderChanged(pos)
	}
	if pos:= scaleSliderPos(input[VoiceVolumeSliderBytePos]); self.State.VoiceVolumeSliderPos != pos {
		self.State.VoiceVolumeSliderPos = pos
		log.Printf("Voice Volume:\t%v",pos)
	}
	if pos:= scaleSliderPos(input[MaleVolumeSliderBytePos]); self.State.MaleVolumeSliderPos != pos {
		self.State.MaleVolumeSliderPos = pos
		log.Printf("Male Volume:\t%v",pos)
	}
	if pos:= scaleSliderPos(input[FemaleVolumeSliderBytePos]); self.State.FemaleVolumeSliderPos != pos {
		self.State.FemaleVolumeSliderPos = pos
		log.Printf("Female Volume:\t%v",pos)
	}

	//Check Buttons
	var btnchange bool

	buttonInput := input[ButtonsBytePos]
	if btnchange, self.State.PlayPauseBtnPressed = checkButton(self.State.PlayPauseBtnPressed,getButtonState(PlayPauseBtnMask,buttonInput));btnchange {
		self.State.PlayPauseState = !self.State.PlayPauseState
		log.Printf("PlayPause pressed, pause:%v",self.State.PlayPauseState)
	}
	if btnchange, self.State.NextBtnPressed = checkButton(self.State.NextBtnPressed,getButtonState(NextBtnMask,buttonInput));btnchange {
		log.Printf("Next pressed")
	}
	if btnchange, self.State.PrevBtnPressed = checkButton(self.State.PrevBtnPressed,getButtonState(PrevBtnMask,buttonInput));btnchange {
		log.Printf("Prev pressed")
	}
	if btnchange, self.State.RecordBtnPressed = checkButton(self.State.RecordBtnPressed,getButtonState(RecordBtnMask,buttonInput));btnchange {
		log.Printf("Record pressed")
	}
	if btnchange, self.State.KeyUpBtnPressed = checkButton(self.State.KeyUpBtnPressed,getButtonState(KeyUpBtnMask,buttonInput));btnchange {
		log.Printf("Up pressed")
	}
	if btnchange, self.State.KeyDownBtnPressed = checkButton(self.State.KeyDownBtnPressed,getButtonState(KeyDownBtnMask,buttonInput));btnchange {
		log.Printf("Down pressed")
	}
	if btnchange, self.State.PitchUpBtnPressed = checkButton(self.State.PitchUpBtnPressed,getButtonState(PitchUpBtnMask,buttonInput));btnchange {
		log.Printf("Pitch Up pressed")
	}
	if btnchange, self.State.PitchDownBtnPressed = checkButton(self.State.PitchDownBtnPressed,getButtonState(PitchDownBtnMask,buttonInput));btnchange {
		log.Printf("Pitch Down pressed")
	}

}

func (self *WarriorController) poll() {
	log.Println("Start Polling")

	for {
		select {
		case <-self.pollStop:
			log.Println("Stop Polling")

			return

		default:
			winputchan := self.wdev.ReadCh()
			self.detectChanges(<-winputchan)
		}
	}

}


//returns (bool: true when the button is pressed initially; new value of the state)
func checkButton(oldVal,newVal bool) (bool,bool){
	if !newVal {
		return false,false
	}
	return oldVal != newVal,true
}
func getButtonState(input byte,mask byte) bool{
	return (input & mask) != 0
}
