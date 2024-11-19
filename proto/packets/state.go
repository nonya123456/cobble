package packets

type State int32

const (
	StateHandshaking State = iota
	StateStatus
	StateLogin
	StateTransfer
)

var allStates = map[State]struct{}{StateHandshaking: {}, StateStatus: {}, StateLogin: {}, StateTransfer: {}}

func (s State) IsValid() bool {
	_, ok := allStates[s]
	return ok
}
