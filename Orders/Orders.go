package Orders

import "time"

const NUMBUTTONS = 3

type HallOrders [][]bool

func InitializeHallOrders(setNumberOfFloors int) HallOrders {
	ho := make(HallOrders, setNumberOfFloors)
	for floor := range ho {
		ho[floor] = make([]bool, NUMBUTTONS-1)
	}
	return ho
}
func HallOrdersLogicalOR(ID int, allHallOrders AllHallOrders) HallOrders {
	ho := InitializeHallOrders(len(allHallOrders[ID]))
	for _, floors := range allHallOrders {
		for floor, directions := range floors {
			for direction, order := range directions {
				ho[floor][direction] = ho[floor][direction] || order
			}
		}
	}
	return ho
}

type AllHallOrders []HallOrders

func InitializeAllHallOrders(setNumberOfElevators int, setNumberOfFloors int) AllHallOrders {
	aho := make(AllHallOrders, setNumberOfElevators)
	for elevator := range aho {
		aho[elevator] = InitializeHallOrders(setNumberOfFloors)
	}
	return aho
}

type HallOrdersTime [][]time.Time

func InitializeHallOrdersTime(setNumberOfFloors int) HallOrdersTime {
	ahot := make(HallOrdersTime, setNumberOfFloors)
	for floor := range ahot {
		ahot[floor] = make([]time.Time, NUMBUTTONS-1)
	}
	return ahot
}

type CabOrders []bool

func InitializeCabOrders(setNumberOfFloors int) CabOrders {
	co := make(CabOrders, setNumberOfFloors)
	return co
}

type AllCabOrders []CabOrders

type CabOrdersTime []time.Time

func InitializeCabOrdersTime(setNumberOfFloors int) CabOrdersTime {
	cot := make(CabOrdersTime, setNumberOfFloors)
	return cot
}

func CabOrdersTimeToBool(cabOrdersTime CabOrdersTime) CabOrders {
	co := InitializeCabOrders(len(cabOrdersTime))
	for floor, order := range cabOrdersTime {
		co[floor] = !order.IsZero()
	}
	return co
}
