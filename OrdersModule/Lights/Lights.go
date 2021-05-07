package Lights

import (
	"github.com/TTK4145-Students-2021/project-gruppe48/Elevio"
	"github.com/TTK4145-Students-2021/project-gruppe48/Orders"
	"github.com/TTK4145-Students-2021/project-gruppe48/OrdersAndStates"
)

var hallLightsCurrentlySet Orders.HallOrders
var cabLightsCurrentlySet Orders.CabOrders

func Init(setNumberOfFloors int) {
	hallLightsCurrentlySet = Orders.InitializeHallOrders(setNumberOfFloors)
	cabLightsCurrentlySet = Orders.InitializeCabOrders(setNumberOfFloors)

	for floor, directions := range hallLightsCurrentlySet {
		for direction, order := range directions {
			Elevio.SetButtonLamp(Elevio.ButtonType(direction), floor, order)

		}
	}
	for floor, order := range cabLightsCurrentlySet {
		Elevio.SetButtonLamp(Elevio.BT_Cab, floor, order)
	}
}

func Update(ID int, allOrdersAndStates OrdersAndStates.AllOrdersAndStates) {
	allHallOrders := Orders.InitializeAllHallOrders(len(allOrdersAndStates), len(allOrdersAndStates[ID].HallOrders))
	for elevator := range allOrdersAndStates {
		allHallOrders[elevator] = OrdersAndStates.HallOrdersTimeToBool(elevator, allOrdersAndStates)
	}
	hallLightsToBeSet := Orders.HallOrdersLogicalOR(ID, allHallOrders)
	cabLightsToBeSet := Orders.CabOrdersTimeToBool(allOrdersAndStates[ID].CabOrders)
	hallLightsCurrentlySet = updateHallLights(hallLightsToBeSet, hallLightsCurrentlySet)
	cabLightsCurrentlySet = updateCabLights(cabLightsToBeSet, cabLightsCurrentlySet)
}

func updateHallLights(hallLightsToBeSet Orders.HallOrders, hallLightsCurrentlySet Orders.HallOrders) Orders.HallOrders {
	for floor, directions := range hallLightsToBeSet {
		for direction, order := range directions {
			if order != hallLightsCurrentlySet[floor][direction] {
				Elevio.SetButtonLamp(Elevio.ButtonType(direction), floor, order)
			}
		}
	}
	return hallLightsToBeSet
}

func updateCabLights(cabLightsToBeSet Orders.CabOrders, cabLightsCurrentlySet Orders.CabOrders) Orders.CabOrders {
	for floor, order := range cabLightsToBeSet {
		if order != cabLightsCurrentlySet[floor] {
			Elevio.SetButtonLamp(Elevio.BT_Cab, floor, order)
		}
	}
	return cabLightsToBeSet
}
