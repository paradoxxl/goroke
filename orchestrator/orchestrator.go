package orchestrator

import (
	"github.com/paradoxxl/goroke/network"
	"github.com/paradoxxl/goroke/hid"
	"net"
)

type Orchestrator struct{
	NetworkController *network.NetworkController
	WarriorController *HID.WarriorController
}

func (self *Orchestrator  )Setup(KaraFunIP net.IP, KaraFunPort uint16){
	self.WarriorController = HID.NewWarrior()
	self.NetworkController = network.NewNetworkController(KaraFunIP,KaraFunPort)

	self.NetworkController.Connect()
	self.WarriorController.NetworkController = self.NetworkController
}

func (self *Orchestrator) Start(){
	self.WarriorController.Connect()
	self.WarriorController.StartPolling()
}
