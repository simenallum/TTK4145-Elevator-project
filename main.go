package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/TTK4145-Students-2021/project-gruppe48/Elevator"
	"github.com/TTK4145-Students-2021/project-gruppe48/Elevio"
	"github.com/TTK4145-Students-2021/project-gruppe48/NetworkModule"
	"github.com/TTK4145-Students-2021/project-gruppe48/OrdersAndStates"
	"github.com/TTK4145-Students-2021/project-gruppe48/OrdersModule"
	SingleElevatorModule "github.com/TTK4145-Students-2021/project-gruppe48/SingleElevator"
)

func main() {
	var orderTimeOut int
	var numFloors int
	var numElevators int
	var elevatorID int
	var simulatorPort string

	flag.IntVar(&orderTimeOut, "orderTimeOut", 30, "Choose OrderTimeOut (30 is standard)")
	flag.IntVar(&numFloors, "numFloors", 4, "Choose Number of floors (4 is standard)")
	flag.IntVar(&numElevators, "numElevators", 3, "Choose number of elevators (3 is standard)")
	flag.IntVar(&elevatorID, "id", 0, "Choose ElevatorID (0 is standard)")
	flag.StringVar(&simulatorPort, "port", "15660", "Choose port (15660 is standard)")
	flag.Parse()

	dir, _ := filepath.Abs(".")
	os.Chmod(dir+"/DistributorModule/hall_request_assigner", 0111) //Execute permission for the hall_request_assigner

	bufsize := 1

	// channels
	elevatorOrdersCh := make(chan Elevator.ElevatorOrders, bufsize)
	elevatorStateCh := make(chan Elevator.ElevatorStates, bufsize)
	fromNeworkCh := make(chan OrdersAndStates.NetworkMessage, numElevators*bufsize)
	toNeworkCh := make(chan OrdersAndStates.NetworkMessage, numElevators*bufsize)
	connectedNodesCh := make(chan []bool, bufsize)
	stopButtonCh := make(chan bool, bufsize)
	buttonPressCh := make(chan Elevio.ButtonEvent, bufsize)
	obstructedCh := make(chan bool, bufsize)
	atFloorCh := make(chan int, bufsize)

	//Init functions
	Elevio.Init("localhost:"+simulatorPort, numFloors)
	NetworkModule.Init(elevatorID, numElevators, fromNeworkCh, connectedNodesCh, toNeworkCh)

	// Goroutines
	go Elevio.PollStopButton(stopButtonCh)
	go Elevio.PollButtons(buttonPressCh)
	go Elevio.PollObstructionSwitch(obstructedCh)
	go Elevio.PollFloorSensor(atFloorCh)

	go SingleElevatorModule.SingleElevatorModule(
		numFloors,
		elevatorID,

		elevatorStateCh,
		elevatorOrdersCh)

	go OrdersModule.OrdersModule(
		elevatorID,
		numElevators,
		numFloors,
		orderTimeOut,

		buttonPressCh,
		elevatorStateCh,
		connectedNodesCh,
		fromNeworkCh,
		elevatorOrdersCh,
		toNeworkCh)

	// For select for keeping program alive.
	for {
		select {
		case <-stopButtonCh:
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
