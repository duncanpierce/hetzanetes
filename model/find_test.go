package model

import (
	"testing"
	"time"
)

func TestPhaseUpTo(t *testing.T) {
	status := NodeStatus{Phases: PhaseChanges{
		PhaseChange{
			Phase:  Active,
			Time:   time.Now(),
			Reason: "bootstrapped",
		},
	}}
	if PhaseUpTo(Creating)(status) {
		t.Fail()
	}
	if PhaseUpTo(Joining)(status) {
		t.Fail()
	}
	if !PhaseUpTo(Active)(status) {
		t.Fail()
	}
	if !PhaseUpTo(Delete)(status) {
		t.Fail()
	}
}
