package main


import (
	"github.com/paradoxxl/goroke/orchestrator"
	"net"
)

/*
const productID = 0x1117
const vendorID = 0x07C0

func main() {

	wchan := hid.FindDevices(vendorID, productID)
	winfo := <-wchan
	if winfo == nil {
		return
	}
	fmt.Printf("Warrior found: PID: %#04X \t VID: %#04X\n", winfo.ProductId, winfo.VendorId)
	fmt.Println(winfo.Path)

	wdev, err := winfo.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Warrior opened: PID: %#04X \t VID: %#04X\n", winfo.ProductId, winfo.VendorId)

	defer wdev.Close()
	cstop := make(chan interface{})

	go readInput(wdev)

	<- cstop

}

func readInput(wdev hid.Device){
	for{
		winputchan := wdev.ReadCh()
		fmt.Println(<-winputchan)
	}

}
/*

 */
func main(){
	orchestrator := orchestrator.Orchestrator{}
	orchestrator.Setup(net.ParseIP("127.0.0.1"),57570)

	orchestrator.Start()

	cstop := make(chan interface{})
	<- cstop

}
