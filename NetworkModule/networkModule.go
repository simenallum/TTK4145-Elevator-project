package NetworkModule

import (
	"fmt"
	"strconv"

	"github.com/TTK4145-Students-2021/project-gruppe48/NetworkModule/network/bcast"
	"github.com/TTK4145-Students-2021/project-gruppe48/NetworkModule/network/peers"
	"github.com/TTK4145-Students-2021/project-gruppe48/OrdersAndStates"
)

const (
	PEERPORT  = 20258
	BCASTPORT = 20259
)

func Init(elevatorID int,
	numElevators int,
	fromNeworkCh chan<- OrdersAndStates.NetworkMessage,
	connectedNodesCh chan<- []bool,

	toNeworkCh <-chan OrdersAndStates.NetworkMessage) {

	bufsize := 1
	peerUpdateCh := make(chan peers.PeerUpdate, bufsize)
	peerTxEnableCh := make(chan bool, bufsize)

	go peers.Transmitter(PEERPORT, fmt.Sprint(elevatorID), peerTxEnableCh)
	go peers.Receiver(PEERPORT, peerUpdateCh)
	go bcast.Transmitter(BCASTPORT, toNeworkCh)
	go bcast.Receiver(BCASTPORT, fromNeworkCh)

	go elevatorsAlive(elevatorID, numElevators, peerUpdateCh, connectedNodesCh)
}

func elevatorsAlive(elevatorID int,
	numElevators int,

	peerUpdateCh <-chan peers.PeerUpdate,

	connectedNodesCh chan<- []bool) {

	alive := make([]bool, numElevators)
	alive[elevatorID] = true

	for {
		select {
		case peerList := <-peerUpdateCh:
			if len(peerList.New) != 0 {
				ID, _ := strconv.Atoi(peerList.New)
				alive[ID] = true
			}
			if len(peerList.Lost) != 0 {
				for _, lostNode := range peerList.Lost {
					ID, _ := strconv.Atoi(lostNode)
					if ID != elevatorID {
						alive[ID] = false
					}

				}
			}

			connectedNodesCh <- alive
		}
	}
}
