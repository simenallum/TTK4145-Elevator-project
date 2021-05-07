package OrdersModule

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/TTK4145-Students-2021/project-gruppe48/DistributorModule"
	"github.com/TTK4145-Students-2021/project-gruppe48/Elevator"
	"github.com/TTK4145-Students-2021/project-gruppe48/Elevio"
	"github.com/TTK4145-Students-2021/project-gruppe48/Orders"
	"github.com/TTK4145-Students-2021/project-gruppe48/OrdersAndStates"
	"github.com/TTK4145-Students-2021/project-gruppe48/OrdersModule/Lights"
)

func OrdersModule(
	ID int,
	setNumberOfElevators int,
	setNumberOfFloors int,
	orderTimeOut int,
	//Ingoing channels:
	buttonPressCh <-chan Elevio.ButtonEvent,
	elevatorStatesCh <-chan Elevator.ElevatorStates,
	connectedNodesCh <-chan []bool,
	fromNetworkCh <-chan OrdersAndStates.NetworkMessage,
	//Outgoing channels:
	elevatorOrdersCh chan<- Elevator.ElevatorOrders,
	toNetworkCh chan<- OrdersAndStates.NetworkMessage) {

	var connectedNodes []bool
	var allOrdersAndStates OrdersAndStates.AllOrdersAndStates
	var allAssignedOrders Orders.AllHallOrders

	connectedNodes = make([]bool, setNumberOfElevators)
	allOrdersAndStates = OrdersAndStates.InitializeAllOrdersAndStates(setNumberOfElevators, setNumberOfFloors)
	allAssignedOrders = Orders.InitializeAllHallOrders(setNumberOfElevators, setNumberOfFloors)

	allOrdersAndStates = readBackupLocalOrders(ID, allOrdersAndStates)
	Lights.Init(setNumberOfFloors)
	Lights.Update(ID, allOrdersAndStates)
	elevatorOrdersCh <- toElevatorFormat(ID, allOrdersAndStates, allAssignedOrders)

	checkForOldOrdersTicker := time.NewTicker(time.Second)

	for {
		select {
		//Prioritized Channel
		case elevatorStates := <-elevatorStatesCh:
			newAssignedOrders := false
			allOrdersAndStates = deleteServedHallOrders(ID, allOrdersAndStates)
			allOrdersAndStates = deleteServedCabOrders(ID, allOrdersAndStates)
			allAssignedOrders = deleteServedAssignedOrders(ID, allAssignedOrders, allOrdersAndStates)
			Lights.Update(ID, allOrdersAndStates)
			saveLocalOrders(ID, allOrdersAndStates)
			if stuckStatusChanged(ID, elevatorStates, allOrdersAndStates) {
				allOrdersAndStates[ID].States = elevatorStates
				allAssignedOrders = DistributorModule.Redistribute(ID, allOrdersAndStates, connectedNodes)
				newAssignedOrders = true
				fmt.Println(allAssignedOrders)
				elevatorOrdersCh <- toElevatorFormat(ID, allOrdersAndStates, allAssignedOrders)
			}
			allOrdersAndStates[ID].States = elevatorStates
			toNetworkCh <- toNetworkMessageFormat(ID, allOrdersAndStates, newAssignedOrders, allAssignedOrders)
		//Less Prioritized channels
		default:
			select {
			case buttonPress := <-buttonPressCh:
				var newAssignedOrders bool
				if buttonPress.Button == Elevio.BT_Cab {
					allOrdersAndStates = addCabButtonPress(ID, buttonPress.Floor, allOrdersAndStates)
					saveLocalOrders(ID, allOrdersAndStates)
					newAssignedOrders = false
				} else {
					allOrdersAndStates = addHallButtonPress(ID, buttonPress.Floor, int(buttonPress.Button), allOrdersAndStates)
					allAssignedOrders = DistributorModule.Redistribute(ID, allOrdersAndStates, connectedNodes)
					newAssignedOrders = true
				}
				Lights.Update(ID, allOrdersAndStates)
				elevatorOrdersCh <- toElevatorFormat(ID, allOrdersAndStates, allAssignedOrders)
				toNetworkCh <- toNetworkMessageFormat(ID, allOrdersAndStates, newAssignedOrders, allAssignedOrders)
			case networkMessage := <-fromNetworkCh:
				if networkMessage.ID != ID {
					allOrdersAndStates = addOrdersAndStates(networkMessage, allOrdersAndStates)
					allOrdersAndStates = deleteServedHallOrders(networkMessage.ID, allOrdersAndStates)
					Lights.Update(ID, allOrdersAndStates)
					if networkMessage.NewAssigned {
						allAssignedOrders = networkMessage.AllAssigned
						elevatorOrdersCh <- toElevatorFormat(ID, allOrdersAndStates, allAssignedOrders)
					}
				}
			case connectedNodes = <-connectedNodesCh:
				allAssignedOrders = DistributorModule.Redistribute(ID, allOrdersAndStates, connectedNodes)
				newAssignedOrders := true
				elevatorOrdersCh <- toElevatorFormat(ID, allOrdersAndStates, allAssignedOrders)
				toNetworkCh <- toNetworkMessageFormat(ID, allOrdersAndStates, newAssignedOrders, allAssignedOrders)
			case <-checkForOldOrdersTicker.C:
				var foundOldOrder bool
				foundOldOrder, allOrdersAndStates, allAssignedOrders = addOldOrders(ID, orderTimeOut, allOrdersAndStates, allAssignedOrders)
				if foundOldOrder {
					newAssignedOrders := false
					toNetworkCh <- toNetworkMessageFormat(ID, allOrdersAndStates, newAssignedOrders, allAssignedOrders)
					elevatorOrdersCh <- toElevatorFormat(ID, allOrdersAndStates, allAssignedOrders)
				}
			}
		}
	}
}

func addCabButtonPress(ID int, floor int, allOrdersAndStates OrdersAndStates.AllOrdersAndStates) OrdersAndStates.AllOrdersAndStates {
	allOrdersAndStates[ID].CabOrders[floor] = time.Now()
	return allOrdersAndStates
}

func saveLocalOrders(ID int, allOrdersAndStates OrdersAndStates.AllOrdersAndStates) {
	localCabOrders := Orders.CabOrdersTimeToBool(allOrdersAndStates[ID].CabOrders)
	fileName := "./OrdersModule/CabOrdersBackup/localOrdersBackup" + fmt.Sprint(ID) + ".txt"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error in creating file. Error: ", err)
	} else {
		defer file.Close()
	}

	for _, order := range localCabOrders {
		if order {
			fmt.Fprintf(file, "1")
		} else {
			fmt.Fprintf(file, "0")
		}
	}

}

func readBackupLocalOrders(ID int, allOrdersAndStates OrdersAndStates.AllOrdersAndStates) OrdersAndStates.AllOrdersAndStates {
	fileName := "./OrdersModule/CabOrdersBackup/localOrdersBackup" + fmt.Sprint(ID) + ".txt"
	fromFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		saveLocalOrders(ID, allOrdersAndStates)
	} else {
		for floor, order := range fromFile {
			if order == 49 { // 49 is the ASCII for 1 as a char
				allOrdersAndStates[ID].CabOrders[floor] = time.Now()
			}
		}

	}
	return allOrdersAndStates
}

func toElevatorFormat(ID int, allOrdersAndStates OrdersAndStates.AllOrdersAndStates, allAssignedOrders Orders.AllHallOrders) Elevator.ElevatorOrders {
	localCabOrders := Orders.CabOrdersTimeToBool(allOrdersAndStates[ID].CabOrders)
	N_BUTTONS := 3
	eo := Elevator.InitializeElevatorOrders(len(localCabOrders), N_BUTTONS)
	for floor, directions := range allAssignedOrders[ID] {
		for direction, order := range directions {
			eo[floor][direction] = order
		}
	}
	cab := 2
	for floor, order := range localCabOrders {
		eo[floor][cab] = order
	}
	return eo
}

func addHallButtonPress(ID int, floor int, direction int, allOrdersAndStates OrdersAndStates.AllOrdersAndStates) OrdersAndStates.AllOrdersAndStates {
	allOrdersAndStates[ID].HallOrders[floor][direction] = time.Now()
	return allOrdersAndStates
}

func toNetworkMessageFormat(ID int, allOrdersAndStates OrdersAndStates.AllOrdersAndStates, newAssignedOrders bool, allAssignedOrders Orders.AllHallOrders) OrdersAndStates.NetworkMessage {
	nm := OrdersAndStates.NetworkMessage{
		ID:           ID,
		OrdersStates: allOrdersAndStates[ID],
		NewAssigned:  newAssignedOrders,
		AllAssigned:  allAssignedOrders,
	}
	return nm
}

func deleteServedHallOrders(ID int, allOrdersAndStates OrdersAndStates.AllOrdersAndStates) OrdersAndStates.AllOrdersAndStates {
	for elevator, ordersAndStates := range allOrdersAndStates {
		for floor, directions := range ordersAndStates.HallOrders {
			servedTime := allOrdersAndStates[ID].States.LastVisit[floor]
			for direction, orderTime := range directions {
				if servedTime.After(orderTime) && !orderTime.IsZero() {
					allOrdersAndStates[elevator].HallOrders[floor][direction] = time.Time{}
				}
			}
		}

	}
	return allOrdersAndStates
}

func deleteServedCabOrders(ID int, allOrdersAndStates OrdersAndStates.AllOrdersAndStates) OrdersAndStates.AllOrdersAndStates {
	for floor, orderTime := range allOrdersAndStates[ID].CabOrders {
		servedTime := allOrdersAndStates[ID].States.LastVisit[floor]
		if servedTime.After(orderTime) && !orderTime.IsZero() {
			allOrdersAndStates[ID].CabOrders[floor] = time.Time{}
		}
	}
	return allOrdersAndStates
}

func deleteServedAssignedOrders(ID int, allAssignedOrders Orders.AllHallOrders, allOrdersAndStates OrdersAndStates.AllOrdersAndStates) Orders.AllHallOrders {
	if allOrdersAndStates[ID].States.Behaviour == Elevator.EB_DoorOpen {
		for direction := range allAssignedOrders[ID][allOrdersAndStates[ID].States.Floor] {
			floor := allOrdersAndStates[ID].States.Floor
			allAssignedOrders[ID][floor][direction] = false
		}
	}
	return allAssignedOrders
}

func stuckStatusChanged(ID int, elevatorStates Elevator.ElevatorStates, allOrdersAndStates OrdersAndStates.AllOrdersAndStates) bool {
	return elevatorStates.Stuck != allOrdersAndStates[ID].States.Stuck
}
func addOrdersAndStates(networkMessage OrdersAndStates.NetworkMessage, allOrdersAndStates OrdersAndStates.AllOrdersAndStates) OrdersAndStates.AllOrdersAndStates {
	allOrdersAndStates[networkMessage.ID] = networkMessage.OrdersStates
	return allOrdersAndStates
}

func addOldOrders(ID int, orderTimeOut int, allOrdersAndStates OrdersAndStates.AllOrdersAndStates, allAssignedOrders Orders.AllHallOrders) (bool, OrdersAndStates.AllOrdersAndStates, Orders.AllHallOrders) {
	foundOldOrder := false
	for elevators, ordersAndStates := range allOrdersAndStates {
		for floor, directions := range ordersAndStates.HallOrders {
			for direction, order := range directions {
				if time.Since(order) > time.Duration(orderTimeOut)*time.Second && !order.IsZero() {
					foundOldOrder = true
					allOrdersAndStates[elevators].HallOrders[floor][direction] = time.Now()
					allAssignedOrders[ID][floor][direction] = true //Add the order to this elevators assigned orders to make sure it is served
				}
			}
		}
	}
	return foundOldOrder, allOrdersAndStates, allAssignedOrders
}
