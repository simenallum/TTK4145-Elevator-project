package SingleElevatorModule

import (
	"fmt"
	"time"

	"github.com/TTK4145-Students-2021/project-gruppe48/Elevator"
	"github.com/TTK4145-Students-2021/project-gruppe48/Elevio"
	"github.com/TTK4145-Students-2021/project-gruppe48/SingleElevator/Fsm"
	"github.com/TTK4145-Students-2021/project-gruppe48/SingleElevator/Timer"
)

func SingleElevatorModule(
	numFloors int,
	elevatorID int,

	elevatorStatesCh chan<- Elevator.ElevatorStates,

	elevatorOrdersCh <-chan Elevator.ElevatorOrders) {

	bufsize := 1

	atFloorCh := make(chan int, bufsize)
	go Elevio.PollFloorSensor(atFloorCh)

	obstructedCh := make(chan bool, bufsize)
	go Elevio.PollObstructionSwitch(obstructedCh)

	doorTimeOutCh := make(chan bool, bufsize)
	stuckTimeOutCh := make(chan bool, bufsize)
	go Timer.TimeOut(stuckTimeOutCh, doorTimeOutCh)

	Fsm.Init(numFloors)

	Fsm.SendElevatorStateUpdate(elevatorStatesCh)

	SendElevatorStatesTicker := time.NewTicker(20 * time.Millisecond)

	fmt.Println("Single elevator started..")
	for {
		select {
		//Prioritized Channels
		case order := <-elevatorOrdersCh:
			Fsm.OnRecievedNewOrders(numFloors, order)
		//Less Prioritized Channels
		default:
			select {

			case floor := <-atFloorCh:
				Fsm.OnFloorArrival(numFloors, floor)

			case <-doorTimeOutCh:
				Fsm.OnDoorTimeout(numFloors)

			case <-obstructedCh:
				Fsm.OnObstruction()

			case <-stuckTimeOutCh:
				Fsm.OnStuck()
				fmt.Println("Single elevator stuck. ID = ", elevatorID)

			case <-SendElevatorStatesTicker.C:

				Fsm.SendElevatorStateUpdate(elevatorStatesCh)

			}
		}
	}

}
