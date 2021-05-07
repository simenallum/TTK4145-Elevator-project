package Fsm

import (
	"time"

	"github.com/TTK4145-Students-2021/project-gruppe48/Elevator"
	"github.com/TTK4145-Students-2021/project-gruppe48/Elevio"
	"github.com/TTK4145-Students-2021/project-gruppe48/SingleElevator/Request"
	"github.com/TTK4145-Students-2021/project-gruppe48/SingleElevator/Timer"
)

var elevator Elevator.ElevatorStates
var previousElevatorStates Elevator.ElevatorStates
var elevatorOrders Elevator.ElevatorOrders

func Init(numFloors int) {
	previousElevatorStates = Elevator.InitializeElevatorStates(numFloors)
	elevator = Elevator.InitializeElevatorStates(numFloors)
	elevatorOrders = Elevator.InitializeElevatorOrders(numFloors, Elevator.N_BUTTONS)

	if !onLegalFloor() {
		setElevatorDirection(Elevator.ED_Up)
		elevator.Direction = Elevator.ED_Up
		elevator.Behaviour = Elevator.EB_Moving
		for !onLegalFloor() {
		}
		setElevatorDirection(Elevator.ED_Stop)
	}
	elevator.Direction = Elevator.ED_Stop
	elevator.Behaviour = Elevator.EB_Idle
	closeDoor()
	elevator.Floor = Elevio.GetFloor()
}

func OnRecievedNewOrders(numFloors int, newOrders Elevator.ElevatorOrders) {
	elevatorOrders = newOrders
	switch elevator.Behaviour {
	case Elevator.EB_DoorOpen:
		if Request.OrdersAtCurrentFloor(elevator, elevatorOrders) {
			Timer.ResetDoorTimer()
			elevator = updateLastVisit()
			Timer.ResetStuckTimer()
			elevator = updateStuckStatus(elevator)
			elevatorOrders = Request.ClearAtCurrentFloor(elevator, elevatorOrders)
		}
	case Elevator.EB_Idle:
		if Request.OrdersAtCurrentFloor(elevator, elevatorOrders) {

			openDoor()
			Timer.ResetDoorTimer()
			elevator = updateLastVisit()
			elevator.Behaviour = Elevator.EB_DoorOpen
			Timer.ResetStuckTimer()
			elevator = updateStuckStatus(elevator)
			elevatorOrders = Request.ClearAtCurrentFloor(elevator, elevatorOrders)

		} else {
			elevator.Direction = Request.ChooseDirection(numFloors, elevator, elevatorOrders)
			if elevator.Direction != Elevator.ED_Stop {
				setElevatorDirection(elevator.Direction)
				elevator.Behaviour = Elevator.EB_Moving
				Timer.ResetStuckTimer()
				elevator = updateStuckStatus(elevator)
			}
		}
	}
}

func OnFloorArrival(numFloors int, arrivedFloor int) {
	elevator.Floor = arrivedFloor
	elevator = updateStuckStatus(elevator)
	Elevio.SetFloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case Elevator.EB_Moving:
		if Request.ShouldStop(numFloors, elevator, elevatorOrders) {
			setElevatorDirection(Elevator.ED_Stop)
			if Request.OrdersAtCurrentFloor(elevator, elevatorOrders) {
				openDoor()
				elevatorOrders = Request.ClearAtCurrentFloor(elevator, elevatorOrders)
				Timer.ResetDoorTimer()
				elevator = updateLastVisit()
				elevator.Behaviour = Elevator.EB_DoorOpen
				Timer.ResetStuckTimer()
			} else {
				elevator.Behaviour = Elevator.EB_Idle
				Timer.StopStuckTimer()
				elevator.Direction = Request.ChooseDirection(numFloors, elevator, elevatorOrders)
				if elevator.Direction != Elevator.ED_Stop {
					setElevatorDirection(elevator.Direction)
					elevator.Behaviour = Elevator.EB_Moving
					Timer.ResetStuckTimer()
				}
			}

		}
	}
}

func OnDoorTimeout(numFloors int) {

	switch elevator.Behaviour {
	case Elevator.EB_DoorOpen:
		elevator.Direction = Request.ChooseDirection(numFloors, elevator, elevatorOrders)
		closeDoor()
		setElevatorDirection(elevator.Direction)
		elevator = updateStuckStatus(elevator)
		if elevator.Direction == Elevator.ED_Stop {
			elevator.Behaviour = Elevator.EB_Idle
			Timer.StopStuckTimer()
		} else {
			elevator.Behaviour = Elevator.EB_Moving
			Timer.ResetStuckTimer()

		}
	}
}

func OnObstruction() {
	switch elevator.Behaviour {
	case Elevator.EB_DoorOpen:
		Timer.ResetDoorTimer()
		elevator = updateLastVisit()
	}

}

func OnStuck() {
	elevator.Stuck = true
}

func SendElevatorStateUpdate(elevatorStatesCh chan<- Elevator.ElevatorStates) {
	if !Elevator.CompareElevatorStates(elevator, previousElevatorStates) {
		elevatorStatesCh <- elevator
	}
	previousElevatorStates = elevator
}

func setElevatorDirection(dir Elevator.ElevatorDirection) {
	switch dir {
	case Elevator.ED_Down:
		Elevio.SetMotorDirection(Elevio.MD_Down)
	case Elevator.ED_Up:
		Elevio.SetMotorDirection(Elevio.MD_Up)
	case Elevator.ED_Stop:
		Elevio.SetMotorDirection(Elevio.MD_Stop)
	}
}

func openDoor() {
	Elevio.SetDoorOpenLamp(true)
}

func closeDoor() {
	Elevio.SetDoorOpenLamp(false)
}

func updateStuckStatus(elevator Elevator.ElevatorStates) Elevator.ElevatorStates {
	if elevator.Stuck {
		elevator.Stuck = false
	}
	return elevator
}

func updateLastVisit() Elevator.ElevatorStates {
	elevator.LastVisit[elevator.Floor] = time.Now()
	return elevator
}

func onLegalFloor() bool {
	return Elevio.GetFloor() != -1
}
