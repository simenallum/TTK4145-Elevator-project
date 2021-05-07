package DistributorModule

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/TTK4145-Students-2021/project-gruppe48/Elevator"
	"github.com/TTK4145-Students-2021/project-gruppe48/Orders"
	"github.com/TTK4145-Students-2021/project-gruppe48/OrdersAndStates"
)

type ElevatorStatus struct {
	Behaviour   Elevator.ElevatorBehaviour `json:"behaviour"`
	Floor       int                        `json:"floor"`
	Direction   Elevator.ElevatorDirection `json:"direction"`
	CabRequests Orders.CabOrders           `json:"cabRequests"`
}

func InitializeElevatorStatus(setNumberOfFloors int) ElevatorStatus {
	return ElevatorStatus{
		Behaviour:   Elevator.ElevatorBehaviour("idle"),
		Floor:       0,
		Direction:   Elevator.ElevatorDirection("stop"),
		CabRequests: Orders.InitializeCabOrders(setNumberOfFloors),
	}

}

type AllElevatorStatus []ElevatorStatus

type AlgorithmInput struct {
	HallRequests Orders.HallOrders      `json:"hallRequests"`
	States       map[int]ElevatorStatus `json:"states"`
}

type AlgorithmOutput map[int]Orders.HallOrders

func Redistribute(ID int, allOrdersAndStates OrdersAndStates.AllOrdersAndStates, connectedNodes []bool) Orders.AllHallOrders {
	allAssignedOrders := Orders.InitializeAllHallOrders(len(connectedNodes), len(allOrdersAndStates[ID].CabOrders))
	availableNodes := findAvailableNodes(connectedNodes, allOrdersAndStates)
	algorithmInput := createAlgorithmInput(ID, allOrdersAndStates, availableNodes)
	algorithmOutput := costAlgorithm(algorithmInput)
	for elevator, hallOrders := range algorithmOutput {
		allAssignedOrders[elevator] = hallOrders
	}
	return allAssignedOrders
}

func costAlgorithm(algorithmInput AlgorithmInput) AlgorithmOutput {
	var algorithmOutput AlgorithmOutput
	inputJSON, err := json.Marshal(algorithmInput)
	PrintERR("MARSHALLING", err)
	dir, err := filepath.Abs(".")
	PrintERR("Filepath", err)
	sh := exec.Command("sh", "-c", dir+"/DistributorModule/hall_request_assigner --clearRequestType all --input '"+string(inputJSON)+"'") //Runs hall_request_assigner in a terminal "behind the scenes" with the JSON as input argument
	outputJSON, err := sh.Output()
	PrintERR("Terminal Output", err)
	err = json.Unmarshal(outputJSON, &algorithmOutput)
	PrintERR("Unmarshalling", err)
	return algorithmOutput
}

func findAvailableNodes(connectedNodes []bool, allOrdersAndStates OrdersAndStates.AllOrdersAndStates) []bool {
	availableNodes := make([]bool, len(connectedNodes))
	for elevatorID, elevator := range allOrdersAndStates {
		availableNodes[elevatorID] = connectedNodes[elevatorID] && !elevator.States.Stuck
	}
	return availableNodes
}

func createAlgorithmInput(ID int, allOrdersAndStates OrdersAndStates.AllOrdersAndStates, availableNodes []bool) AlgorithmInput {
	statusMap := make(map[int]ElevatorStatus)
	for node, available := range availableNodes {
		if available {
			statusMap[node] = ElevatorStatus{
				Behaviour:   allOrdersAndStates[node].States.Behaviour,
				Floor:       allOrdersAndStates[node].States.Floor,
				Direction:   allOrdersAndStates[node].States.Direction,
				CabRequests: Orders.CabOrdersTimeToBool(allOrdersAndStates[node].CabOrders),
			}
		}
	}
	return AlgorithmInput{
		HallRequests: OrdersAndStates.HallOrdersTimeToBool(ID, allOrdersAndStates),
		States:       statusMap,
	}
}

func PrintERR(str string, err error) {
	if err != nil {
		fmt.Println(str, err)

	}
}
