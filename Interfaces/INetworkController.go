package Interfaces

type INetworkController interface {
	GeneralVolumeSliderChanged(byte)
	VoiceVolumeSliderPosChanged(byte)
	MaleVolumeSliderPosChanged(byte)
	FemaleVolumeSliderPosChanged(byte)
	PitchChanged(int8)
	TempoChanged(int8)
	PlayPressed()
	PausePressed()
	NextPressed()
	PrevPressed()
	GetStatus()

}
