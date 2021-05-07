package PrintElevator

import (
	"fmt"
	"strconv"

	"github.com/TTK4145-Students-2021/project-gruppe48/Elevator"
	"github.com/TTK4145-Students-2021/project-gruppe48/Elevio"
)

func Print(numFloors int, e Elevator.ElevatorStates, o Elevator.ElevatorOrders) {
	p := fmt.Printf
	p("  +--------------------+\n")
	p("  |floor = %-2d          |\n", e.Floor)
	p("  |dirn  = %-12.12s|\n", dirnToString(e.Direction))
	p("  |behav = %-12.12s|\n", ebToString(e.Behaviour))
	p("  |Stuck = %-12.12s|\n", strconv.FormatBool(e.Stuck))
	p("  +--------------------+\n")
	p("  |  | up  | dn  | cab |\n")
	for f := numFloors - 1; f >= 0; f-- {
		p("  | %d", f)
		for btn := 0; btn < Elevator.N_BUTTONS; btn++ {
			if f == numFloors-1 && Elevio.ButtonType(btn) == Elevio.BT_HallUp || (f == 0 && Elevio.ButtonType(btn) == Elevio.BT_HallDown) {
				p("|     ")
			} else if o[f][btn] {
				p("|  #  ")
			} else {
				p("|  -  ")
			}
		}
		p("|\n")
	}
	p("  +--------------------+\n")

}

func ebToString(eb Elevator.ElevatorBehaviour) string {
	switch eb {
	case Elevator.EB_Idle:
		return "EB_Idle"
	case Elevator.EB_DoorOpen:
		return "EB_DoorOpen"
	case Elevator.EB_Moving:
		return "EB_Moving"
	default:
		return "EB_Undefined"
	}
}

func dirnToString(direction Elevator.ElevatorDirection) string {
	switch direction {
	case Elevator.ED_Down:
		return "D_Down"
	case Elevator.ED_Up:
		return "D_Up"
	case Elevator.ED_Stop:
		return "D_Stop"
	default:
		return "D_UNDEFINED"
	}
}
