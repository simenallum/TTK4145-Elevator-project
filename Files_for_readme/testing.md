Tests for Elevator Project
================

Disclaimer: This are tests made by us to make sure that the system fulfills the system requirements. This is not the official TTK4145 FAT. 

### No orders are lost

| testNumber | Requirements | Test | Passed / not-passed |
|-|-|-|-|
| 1 | Once the light on a hall call button is turned on, an elevator should arrive at that floor. | Press all Hall order buttons on one of the connected elevators. All Orders should be cleared. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 2 | Once the light on a cab call button is turned on the elevator at that specific workspace should take the order. | Give the elevators cab orders and check that all elevators complete its own cab orders.  | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 3 | The elevator should handle network packet loss | Set package loss on udp port 20258 and 20259 to 80%. Press single orders until one of the clearedFloor messages on the network gets lost. Two or three elevators should serve the 'lost floor'.  | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 4 | The elevator should handle losing network connection entirely. | Give a sensible amount of orders in the system. Block udp port 20258 and 20259. The three elevators should work as seperate elevators and every elevator should complete all orders in the system. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 5 | The elevator should handle software that crashes. | Press the stop button on one of the elevators (This stops the program). The other elevators in the system should take over the disconnected elevators orders. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 6 | The elevator should handle losing motor power. | Disconnect the elevators motor. The other elevators should redistribute the stuck elevators orders and complete them. | <ul><li>- [ ] Simulator</li><li>- [ ] At Lab</li></ul> |
| 7 | The system should handle losing power to one of the machines that controls the elevator. | Turn off one of the computers using 'svenske knappen'. The other elevators should redistribute the disconnected elevators orders and complete them. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 8 | In the case of power/software crash cab orders should be executed once service is restored. | Give a single elevator some cab orders. Press stop. Start the elevator again. The elevator should continue by clearing its cab orders.  | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 9 | The time used to detect these failures should be reasonable, ie. on the order of magnitude of seconds (not minutes) | Time case 3-7 above and check 'time to fix' | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 10 | If the elevator is disconnected from the network, it should still serve all the currently active orders (ie. whatever lights are showing). | Tested by testNumber 4. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 11 | The elevator should keep taking new cab calls whens its disconnected, so that people can exit the elevator even if it is disconnected from the network. | Block udp port 20258 and 20259. Give cab orders to the elevators. All cab orders should be served. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 12 | The elevator software should not require reinitialization (manual restart) after intermittent network loss. | Block udp port 20258 and 20259. Give some orders to all elevators. Unblock network. All orders should be redistributed and served. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 13 | The elevator software should not require reinitialization (manual restart) after motor power loss. | After testNumber 6 connect the motor power. The elevator should drive to the next floor. Give an order on that floor on another elevator and check if the elevator that was stuck serves it. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |

### Multiple elevators should be more efficient than one
| testNumber | Requirements | Test | Passed / not-passed |
|-|-|-|-|
| 14 | The orders should be distributed across the elevators in a reasonable way | If all three elevators are idle. Two of them are at the bottom floor. The last one is at floor 1. A new order at the top floor should be handled by the closest elevator | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 15 | Multiple elevator should be more efficient than one | Test with single elevator: Starts at floor 0. Orders given: U0, D3, U2, D1. Test with two elevators. Both starts at floor 0. Give given order sequence at one of the elevators. Check time. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
 
### An individual elevator should behave sensibly and efficiently
 | testNumber | Requirements | Test | Passed / not-passed |
|-|-|-|-|
| 16 | No stopping at every floor "just to be safe" | Start at floor 0. Give cab order at floor 3. Should drive directly to floor 3. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 17 | The hall "call upward" and "call downward" buttons should behave differently | Start at floor 0. Give order: U1, U2. Check sequence. Start at floor 0. Give order: D1, D2. Check sequence. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
 
### The lights and buttons should function as expected
 | testNumber | Requirements | Test | Passed / not-passed |
|-|-|-|-|
| 18 | The hall call buttons on all workspaces should let you summon an elevator | Tested by testNumber 1 | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 19 | Under normal circumstances, the lights on the hall buttons should show the same thing on all workspaces | Give all orders to one elevator. Check that all lights is set on the two other elevators. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 20 | Under circumstances with high packet loss, at least one light must work as expected | Packetloss on port 20259 to 80%. Press some hall orders on one of the elavators. Check that lights is set on at least on elevator. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 21 | The cab button lights should not be shared between elevators | Give a single elevator some cab orders. Check that cab orders is set on the correct elevator. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 22 | The cab and hall button lights should turn on as soon as is reasonable after the button has been pressed | Tested by testNumber 19 and 21. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 23 | The cab and hall button lights should turn off when the corresponding order has been serviced | All order lights should be cleared in every elevator after test 19 and 21. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |
| 24 | The "door open" lamp should be used as a substitute for an actual door, and as such should not be switched on while the elevator is moving |  | <ul><li>- [ ] Simulator</li><li>- [ ] At Lab</li></ul> |
| 25 | The duration for keeping the door open should be in the 1-5 second range | Time the door-lamp. | <ul><li>- [ ] Simulator</li><li>- [X] At Lab</li></ul> |
| 26 | The door should not close while it is obstructed by the obstruction switch. | Flip the obstruction switch before a elevator reaches a floor. The door should remain open as long as the obstruction is active after it reached the floor. | <ul><li>- [X] Simulator</li><li>- [ ] At Lab</li></ul> |


Unspecified behaviour
---------------------
Which orders are cleared when stopping at a floor
 - The system assumes that everyone enters/exits the elevator when the door opens.
 
How the elevator behaves when it cannot connect to the network (router) during initialization
 - The system enters 'single-elevator-mode' and is available for network connection.
 
How the hall (call up, call down) buttons work when the elevator is disconnected from the network
 - The system takes all new orders.
 - The orders is redestributed when the connection is back.
 
What the stop button does
 - The stop button terminates the corrensponding elevator program.

How the door behaves
 - The door is open for 3 seconds.
 - The door timer is reset when an order is pressed on that floor.

   
Permitted assumptions
---------------------

The following assumptions will always be true during testing:
 1. At least one elevator is always working normally
 2. No multiple simultaneous errors: Only one error happens at a time, but the system must still return to a fully operational state after this error
    - Recall that network packet loss is *not* an error in this context, and must be considered regardless of any other (single) error that can occur
 3. No network partitioning: There will never be a situation where there are multiple sets of two or more elevators with no connection between them
 4. Cab call redundancy with a single elevator is not required
    - Given assumptions **1** and **2**, a system containing only one elevator is assumed to be unable to fail