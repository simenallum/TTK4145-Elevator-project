package Timer

import (
	"time"
)

const (
	DOOR_OPEN_TIME = 3
	STUCK_TIME     = 10
)

var doorTimerActive bool
var stuckTimerActive bool
var doorTimerLastReset time.Time
var stuckTimerLastReset time.Time

func TimeOut(stuckTimeOutCh chan<- bool, doorTimeOutCh chan<- bool) {
	for {
		if time.Since(doorTimerLastReset) > time.Duration(DOOR_OPEN_TIME)*time.Second && doorTimerActive {
			stopDoorTimer()
			doorTimeOutCh <- true
		} else if time.Since(stuckTimerLastReset) > time.Duration(STUCK_TIME)*time.Second && stuckTimerActive {
			StopStuckTimer()
			stuckTimeOutCh <- true
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func ResetDoorTimer() {
	doorTimerLastReset = time.Now()
	doorTimerActive = true
}

func ResetStuckTimer() {
	stuckTimerLastReset = time.Now()
	stuckTimerActive = true
}

func stopDoorTimer() {
	doorTimerActive = false
}

func StopStuckTimer() {
	stuckTimerActive = false
}
