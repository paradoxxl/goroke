package HID

const (
	GeneralVolumeSliderBytePos byte = iota
	VoiceVolumeSliderBytePos
	MaleVolumeSliderBytePos
	FemaleVolumeSliderBytePos
	ButtonsBytePos
)

const (
	PlayPauseBtnMask byte = 1 << iota
	NextBtnMask
	PrevBtnMask
	KeyDownBtnMask
	KeyUpBtnMask
	PitchDownBtnMask
	PitchUpBtnMask
	RecordBtnMask
)

const MaxSliderPosHW float32 = 255.0
const MaxSliderPosSW float32 = 100

const DataLength byte = 6

const (
	productID = 0x1117
	vendorID = 0x07C0
)

