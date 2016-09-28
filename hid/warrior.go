package HID

import (
	"github.com/paradoxxl/hid"
	"log"
	"time"
	"github.com/paradoxxl/goroke/Interfaces"
)

type WarriorController struct {
	Connected     bool
	DevicePresent bool

	State *WarriorControllerState

	winfo *hid.DeviceInfo
	wdev  hid.Device

	pollStop chan interface{}

	NetworkController Interfaces.INetworkController
	Audiorecorder Interfaces.IRecorderController
}

type WarriorControllerState struct {
	GeneralVolumeSliderPos byte
	VoiceVolumeSliderPos   byte
	MaleVolumeSliderPos    byte
	FemaleVolumeSliderPos  byte

	PlayPauseBtnPressed bool
	NextBtnPressed      bool
	PrevBtnPressed      bool
	TempoDownPressed    bool
	TempoUpPressed      bool
	PitchDownBtnPressed bool
	PitchUpBtnPressed   bool
	RecordBtnPressed    bool

	PlayPauseState bool
	PitchState     int8
	TempoState     int8
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
	if pos := scaleSliderPos(input[GeneralVolumeSliderBytePos]); self.State.GeneralVolumeSliderPos != pos {
		self.State.GeneralVolumeSliderPos = pos
		log.Printf("General Volume:\t%v", pos)

		self.NetworkController.GeneralVolumeSliderChanged(pos)

	}
	if pos := scaleSliderPos(input[VoiceVolumeSliderBytePos]); self.State.VoiceVolumeSliderPos != pos {
		self.State.VoiceVolumeSliderPos = pos
		log.Printf("Voice Volume:\t%v", pos)

		self.NetworkController.VoiceVolumeSliderPosChanged(pos)
	}
	if pos := scaleSliderPos(input[MaleVolumeSliderBytePos]); self.State.MaleVolumeSliderPos != pos {
		self.State.MaleVolumeSliderPos = pos
		log.Printf("Male Volume:\t%v", pos)

		self.NetworkController.MaleVolumeSliderPosChanged(pos)

	}
	if pos := scaleSliderPos(input[FemaleVolumeSliderBytePos]); self.State.FemaleVolumeSliderPos != pos {
		self.State.FemaleVolumeSliderPos = pos
		log.Printf("Female Volume:\t%v", pos)

		self.NetworkController.FemaleVolumeSliderPosChanged(pos)

	}

	//Check Buttons
	var btnchange bool

	buttonInput := input[ButtonsBytePos]
	if btnchange, self.State.PlayPauseBtnPressed = checkButton(self.State.PlayPauseBtnPressed, getButtonState(PlayPauseBtnMask, buttonInput)); btnchange {
		if self.State.PlayPauseState = !self.State.PlayPauseState; self.State.PlayPauseState {
			self.NetworkController.PlayPressed()
			self.sendCompleteState()
		} else {
			self.NetworkController.PausePressed()
		}
		log.Printf("PlayPause pressed, pause:%v", self.State.PlayPauseState)

	}
	if btnchange, self.State.NextBtnPressed = checkButton(self.State.NextBtnPressed, getButtonState(NextBtnMask, buttonInput)); btnchange {
		log.Printf("Next pressed")

		self.State.PlayPauseState = Pause

		self.NetworkController.NextPressed()
	}
	if btnchange, self.State.PrevBtnPressed = checkButton(self.State.PrevBtnPressed, getButtonState(PrevBtnMask, buttonInput)); btnchange {
		log.Printf("Prev pressed")

		self.NetworkController.PrevPressed()
	}
	if btnchange, self.State.RecordBtnPressed = checkButton(self.State.RecordBtnPressed, getButtonState(RecordBtnMask, buttonInput)); btnchange {
		log.Printf("Record pressed")
		self.Audiorecorder.RecordButtonPressed()
	}
	if btnchange, self.State.TempoUpPressed = checkButton(self.State.TempoUpPressed, getButtonState(TempoDownBtnMask, buttonInput)); btnchange {
		log.Printf("Tempo+ pressed")

		self.State.TempoState++
		self.NetworkController.TempoChanged(self.State.TempoState)
	}
	if btnchange, self.State.TempoDownPressed = checkButton(self.State.TempoDownPressed, getButtonState(TempoUpBtnMask, buttonInput)); btnchange {
		log.Printf("Tempo- pressed")

		self.State.TempoState--
		self.NetworkController.TempoChanged(self.State.TempoState)
	}
	if btnchange, self.State.PitchUpBtnPressed = checkButton(self.State.PitchUpBtnPressed, getButtonState(PitchUpBtnMask, buttonInput)); btnchange {
		log.Printf("Pitch Up pressed")
		if self.State.PitchState < 6{
			self.State.PitchState++
		}else{
			self.State.PitchState = 6
		}
		self.NetworkController.PitchChanged(self.State.PitchState)
	}
	if btnchange, self.State.PitchDownBtnPressed = checkButton(self.State.PitchDownBtnPressed, getButtonState(PitchDownBtnMask, buttonInput)); btnchange {
		log.Printf("Pitch Down pressed")

		if self.State.PitchState > -5{
			self.State.PitchState--
		}else{
			self.State.PitchState = -6
		}
		self.NetworkController.PitchChanged(self.State.PitchState)
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
			time.Sleep(PollingDelay)
		}
	}

}

//returns (bool: true when the button is pressed initially; new value of the state)
func checkButton(oldVal, newVal bool) (bool, bool) {
	if !newVal {
		return false, false
	}
	return oldVal != newVal, true
}
func getButtonState(input byte, mask byte) bool {
	return (input & mask) != 0
}

func (self *WarriorController) sendCompleteState(){
	self.NetworkController.GeneralVolumeSliderChanged(self.State.GeneralVolumeSliderPos)
	self.NetworkController.VoiceVolumeSliderPosChanged(self.State.VoiceVolumeSliderPos)
	self.NetworkController.MaleVolumeSliderPosChanged(self.State.MaleVolumeSliderPos)
	self.NetworkController.FemaleVolumeSliderPosChanged(self.State.FemaleVolumeSliderPos)

	self.NetworkController.TempoChanged(self.State.TempoState)
	self.NetworkController.PitchChanged(self.State.PitchState)


}
