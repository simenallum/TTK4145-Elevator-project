package Request

import (
	"github.com/TTK4145-Students-2021/project-gruppe48/Elevator"
	"github.com/TTK4145-Students-2021/project-gruppe48/Elevio"
)

func ChooseDirection(numFloors int, e Elevator.ElevatorStates, o Elevator.ElevatorOrders) Elevator.ElevatorDirection {
	switch e.Direction {
	case Elevator.ED_Up:
		if above(numFloors, e, o) {
			return Elevator.ED_Up
		} else if below(e, o) {
			return Elevator.ED_Down
		} else {
			return Elevator.ED_Stop
		}
	case Elevator.ED_Down:
		fallthrough
	case Elevator.ED_Stop:
		if below(e, o) {
			return Elevator.ED_Down
		} else if above(numFloors, e, o) {
			return Elevator.ED_Up
		} else {
			return Elevator.ED_Stop
		}
	default:
		return Elevator.ED_Stop
	}
}

func ShouldStop(numFloors int, e Elevator.ElevatorStates, o Elevator.ElevatorOrders) bool {
	switch e.Direction {
	case Elevator.ED_Down:
		return o[e.Floor][Elevio.BT_HallDown] || o[e.Floor][Elevio.BT_Cab] || !below(e, o)
	case Elevator.ED_Up:
		return o[e.Floor][Elevio.BT_HallUp] || o[e.Floor][Elevio.BT_Cab] || !above(numFloors, e, o)
	case Elevator.ED_Stop:
		fallthrough
	default:
		return true
	}
}

func ClearAtCurrentFloor(e Elevator.ElevatorStates, o Elevator.ElevatorOrders) Elevator.ElevatorOrders {
	for btn := 0; btn < Elevator.N_BUTTONS; btn++ {
		o[e.Floor][btn] = false
	}

	return o
}

func OrdersAtCurrentFloor(e Elevator.ElevatorStates, o Elevator.ElevatorOrders) bool {
	for _, order := range o[e.Floor] {
		if order {
			return true
		}
	}
	return false
}


func above(numFloors int, e Elevator.ElevatorStates, o Elevator.ElevatorOrders) bool {
	for f := e.Floor + 1; f < numFloors; f++ {
		for btn := 0; btn < Elevator.N_BUTTONS; btn++ {
			if o[f][btn] {
				return true
			}
		}
	}
	return false
}

func below(e Elevator.ElevatorStates, o Elevator.ElevatorOrders) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < Elevator.N_BUTTONS; btn++ {
			if o[f][btn] {
				return true
			}
		}
	}
	return false
}