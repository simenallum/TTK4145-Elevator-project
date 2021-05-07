package Elevator

import (
	"time"
)

const (
	N_BUTTONS = 3
)

type ElevatorBehaviour string

const (
	EB_Idle     = "idle"
	EB_Moving   = "moving"
	EB_DoorOpen = "doorOpen"
)

type ElevatorDirection string

const (
	ED_Up   = "up"
	ED_Down = "down"
	ED_Stop = "stop"
)

type TimeOfLastVisit []time.Time

func InitializeTimeOfLastVisit(setNumberOfFloors int) TimeOfLastVisit {
	tolv := make(TimeOfLastVisit, setNumberOfFloors)
	return tolv
}

type ElevatorStates struct {
	Stuck     bool
	Behaviour ElevatorBehaviour
	Floor     int
	Direction ElevatorDirection
	LastVisit TimeOfLastVisit
}

func InitializeElevatorStates(setNumberOfFloors int) ElevatorStates {
	es := ElevatorStates{
		Stuck:     false,
		Behaviour: EB_Idle,
		Floor:     0,
		Direction: ED_Stop,
		LastVisit: InitializeTimeOfLastVisit(setNumberOfFloors),
	}
	return es
}

func CompareElevatorStates(states1 ElevatorStates, states2 ElevatorStates) bool {
	if states1.Stuck == states2.Stuck && states1.Behaviour == states2.Behaviour && states1.Floor == states2.Floor && states1.Direction == states2.Direction {
		equal := true
		for floor := range states1.LastVisit {
			equal = equal && states1.LastVisit[floor].Equal(states2.LastVisit[floor])
		}

		return equal
	}
	return false
}

type ElevatorOrders [][]bool

func InitializeElevatorOrders(setNumberOfFloors int, setNumberOfButtons int) ElevatorOrders {
	eo := make(ElevatorOrders, setNumberOfFloors)
	for floor := range eo {
		eo[floor] = make([]bool, setNumberOfButtons)
	}
	return eo
}
