package OrdersAndStates

import (
	"github.com/TTK4145-Students-2021/project-gruppe48/Elevator"
	"github.com/TTK4145-Students-2021/project-gruppe48/Orders"
)

type OrdersAndStates struct {
	HallOrders Orders.HallOrdersTime
	CabOrders  Orders.CabOrdersTime
	States     Elevator.ElevatorStates
}

func InitializeOrdersAndStates(setNumberOfFloors int) OrdersAndStates {
	oas := OrdersAndStates{
		HallOrders: Orders.InitializeHallOrdersTime(setNumberOfFloors),
		CabOrders:  Orders.InitializeCabOrdersTime(setNumberOfFloors),
		States:     Elevator.InitializeElevatorStates(setNumberOfFloors),
	}
	return oas
}

type NetworkMessage struct {
	ID           int
	OrdersStates OrdersAndStates
	NewAssigned  bool
	AllAssigned  Orders.AllHallOrders
}

type AllOrdersAndStates []OrdersAndStates

func InitializeAllOrdersAndStates(setNumberOfElevators int, setNumberOfFloors int) AllOrdersAndStates {
	aoas := make(AllOrdersAndStates, setNumberOfElevators)
	for elevator := range aoas {
		aoas[elevator] = InitializeOrdersAndStates(setNumberOfFloors)
	}
	return aoas
}

func HallOrdersTimeToBool(ID int, allOrdersAndStates AllOrdersAndStates) Orders.HallOrders {
	ho := Orders.InitializeHallOrders(len(allOrdersAndStates[ID].HallOrders))
	for _, elevator := range allOrdersAndStates {
		for floor, directions := range elevator.HallOrders {
			for direction, orderTime := range directions {
				ho[floor][direction] = ho[floor][direction] || !orderTime.IsZero()
			}
		}
	}
	return ho
}
