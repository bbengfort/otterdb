package replica

import (
	"fmt"
)

// Replica states for distributed consensus.
const (
	Stopped State = iota // stopped should be the zero value and default
	Initialized
	Running
	Follower
	Candidate
	Leader
)

// Names of the states for serialization
var stateStrings = [...]string{
	"stopped", "initialized", "running", "follower", "candidate", "leader",
}

//===========================================================================
// State Enumeration
//===========================================================================

// State is an enumeration of the possible status of a replica.
type State uint8

// String returns a human readable representation of the state.
func (s State) String() string {
	return stateStrings[s]
}

//===========================================================================
// State Transitions
//===========================================================================

// SetState updates the state of the local replica, performing any actions
// related to multiple states, modifying internal private variables as
// needed and calling the correct internal state setting function.
//
// NOTE: These methods are not thread-safe.
func (r *Replica) setState(state State) (err error) {
	switch state {
	case Stopped:
		err = r.setStoppedState()
	case Initialized:
		err = r.setInitializedState()
	case Running:
		err = r.setRunningState()
	case Follower:
		err = r.setFollowerState()
	case Candidate:
		err = r.setCandidateState()
	case Leader:
		err = r.setLeaderState()
	default:
		err = fmt.Errorf("unknown state %q", state)
	}

	if err == nil {
		r.state = state
	}

	return err
}

// Stops all timers that might be running.
func (r *Replica) setStoppedState() error {
	return nil
}

// Resets any volatile variables on the local replica and is called when the
// replica becomes a follower or a candidate.
func (r *Replica) setInitializedState() error {
	return nil
}

// Should only be called once after initialization to bootstrap the quorum by
// starting the leader's heartbeat or starting the election timeout for all
// other replicas.
func (r *Replica) setRunningState() error {
	return nil
}

func (r *Replica) setFollowerState() error {
	return nil
}

func (r *Replica) setCandidateState() error {
	return nil
}

func (r *Replica) setLeaderState() error {
	return nil
}
