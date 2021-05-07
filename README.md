Elevator Project
================


Summary
-------
Software is created to control  `n` elevators working in parallel across `m`  floors.

The elevator project description is available [here](https://github.com/TTK4145/Project)

## Running the elevator
Build the program with the following call in a linux terminal:
> go build main.go

The file can then be run:
> ./main 

The call to ./main takes the following flags (default values in parentheses):
- -id= (0)
- -port= (15660)
- -numElevators= (3)
- -numFloors= (4)
- -orderTimeOut= (30)

**Remarks:** <br />
- The ElevatorID must be unique over a 'session'.
- The elevator only runs on a linux system due to a call to a linux terminal from the code.
- The program does not start unless there is a simulator or a physical elevator on the port provided to the port flag. 

## Our system
Our system is written in [GOLANG](https://golang.org). The system is based on a peer-to-peer design. UDP broadcast is used to send and receive a predefined network message struct between the elevators on a common broadcast-port. In order for the elevator to know if there are more elevators online all elevators broadcasts it's own unique ID on a common peers-port every 15 milliseconds. Every elevator listens to this port and generate a list of connectedNodes based on this. 

Every elevator receives and saves a local version of all button presses from the other elevators. It also receives and stores information about every online elevators states. If one elevator loses connection, gets stuck or die the remaining elevators redistribute all hallorders in the system between them. Each elevator stores its own caborders locally so that in the event that the elevator dies it can initialize itself with these orders.

If an elevator detects that a hall order has not been completed within the orderTimeOut-time, the elevator itself executes this order. 

An elevator works as a seperate elevator when it loses its connection with the others. When connection is restored the orders in the system is redistributed amoung the connected elevators. 

If an elevator has been in a state other than idle for more than 10 seconds, it defines itself as stuck, shares that information with the other elevators and redistributes the orders in the system. The elevator will thus redistribute its hall orders to the other elevators (that is connected and not stuck) in situations of engine failure and in situations where the door is held for more than 10 seconds. When the motor power is restored and it reaches a floor or when the door is closed it returns to normal operation and is ready to receive new hallorders. 


## Modules 

Our system is divided into 7 main modules. 


### Modules written by us
1.  **OrdersModule** <br />
    This module is the heart of our system. It receives from the network, sends to the network, receives local button presses, and sends orders to the local elevator. The module has one simple submodule called Lights.
    1.  **Lights** <br />
    Contains state on currently set lights and a function for updating lights called from the orders module. This reduces the stress on the elevator io module and reduces the runtime of Update() with a factor ~ (number of lights)

2.  **SingleElevatorModule** <br />
    This module consists of four sub-modules and contains all the logic needed for the elevator to receive and execute orders.
    1.  **Fsm** <br />
    Contains functions for making all state transitions.

    2.   **Request** <br />
    Contains functions that put the elevator logic into practice.

    3.   **Timer** <br />
    Contains two timers: one for the door and one to keep track of whether the lift is stuck.

    4.   **PrintElevator** <br />
    A support module used to print the elevator with its states and orders.

3.  **NetworkModule** <br />
    This module contains everything that has to do with network communication. It initializes the used functions from the provided network module and it checks which lifts are alive

4.  **DistributorModule** <br />
    This module contains a pure function shell called from the order module that uses the provided cost algorithm to make an optimal distribution of the hall orders in the system.

5.  **Types And initalize-functions** <br />
    The module contains three packages with types and initialization functions used in the system.
    1.   **Elevator** <br />
    
    2.   **Orders** <br />

    3.   **OrdersAndStates** <br />

### Modules not written by us
6.  **Elevio**
    1. Module for setting the actuators of the system. Setting light, polling buttons and floor sensors. Minor changes by us to "Elevio.PollObstructionSwitch". Written by [@klasbo](https://github.com/klasbo)
7.  **Network**
    1. Module used for UDP communication. Written by [@klasbo](https://github.com/klasbo)

### Functions not written by us
- **hall_request_assigner**
    - Precompiled D-code that distributes the orders in an optimal way in the system. Written by [@klasbo](https://github.com/klasbo)

## Testing of the system reqirements
In addition to FAT performed on 21.04.2021, we have performed a number of tests that we believe reflect the requirements of the system. You can find these tests [here](https://github.com/TTK4145-Students-2021/project-gruppe48/blob/master/Files_for_readme/testing.md).


## All goroutines and channels in the system
![GO-routines and channels](https://github.com/TTK4145-Students-2021/project-gruppe48/blob/master/Files_for_readme/goroutines.png)
All go-routines in blue. The for-loop in main is in orange.

The system has been programmed with a "circle of channels". This is generally not a good thing. A tempting solution to this is just to increase the buffer size of the channels, but since this can never guarantee that partial deadlocks will not occur and it is a bad programming practice, a different approach is chosen. We can get away with buffersize = 1 if the following is true:

1. The receiving case must always be prioritized in the for-select loop.
2. The receiving case at both ends must not send data back.

This allows for receiving data to always be prioritized and deadlocks will not occur.
    
## ClassDiagram
![Class diagram](https://github.com/TTK4145-Students-2021/project-gruppe48/blob/master/Files_for_readme/Classdiagram.png)


## SingleElevator class diagram
![SingleElevator Class diagram](https://github.com/TTK4145-Students-2021/project-gruppe48/blob/master/Files_for_readme/SingleElevator.png)

## Addtitional Comments: 
- In retrospect, we have realized that the use of channels is 'not a free lunch' and we had to reprogram parts of the system to reduce circular channels between different modules. If we were to redesign the system, we would design it so that we totally avoid circular depending channels.