package common

type AppState int

const (
	Main AppState = iota
	PeerChat
)

type StateManager struct {
	state AppState
}

func InitStateManager() StateManager {
	return StateManager{state: Main}
}

func (s StateManager) GetState() AppState {
	return s.state
}

func (s *StateManager) SetState(state AppState) {
	s.state = state
}
