package network

const url = "ws://%v:%v"
const origin = "http://localhost"

const (
	GetStatusCommand = "<action type=\"getStatus\"></action>"
	GetCatalogListCommand = "<action type=\"getCatalogList\"></action>"
	GeneralVolumeCommand = "<action type=\"setVolume\" volume_type=\"general\">%v</action>"
	VoiceVolumeCommand = "<action type=\"setVolume\" volume_type=\"bv\">%v</action>"
	MaleVolumeCommand = "<action type=\"setVolume\" volume_type=\"lead1\">%v</action>"
	FemaleVolumeCommand = "<action type=\"setVolume\" volume_type=\"lead2\">%v</action>"

	PlayCommand = "<action type=\"play\"></action>"
	PauseCommand = "<action type=\"pause\"></action>"
	NextCommand = "<action type=\"next\"></action>"
	SeekCommand = "<action type=\"seek\">%v</action>"
	PitchCommand = "<action type=\"pitch\">%v</action>"
	TempoCommand = "<action type=\"tempo\">%v</action>"


)
