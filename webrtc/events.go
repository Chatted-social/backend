package webrtc

const (
	EventTypeJoinRoom = "join_room"
	EventTypeIceCandidate = "ice-candidate"
	EventTypeAnswer = "answer"
	EventTypeOffer = "offer"
	EventTypeUserJoined = "new_user"
)

type Event struct {
	ID string `json:"id,omitempty"`
}

// WebRTC events
type (
	EventJoinRoom struct {
		Event
		RoomID string `json:"room_id" mapstructure:"room"`

		// Sends only webrtc
		OtherUsers []*Client `json:"other_users,omitempty"`
		UserID string `json:"user_id"`
	}

	EventUserJoined struct {
		Event

		UserID string `json:"user_id"`
	}

	EventHandshake struct {
		Event

		Target string `json:"target"`
		Caller string `json:"caller"`
		SDP interface{} `json:"sdp"`
	}

	EventIceCandidate struct {
		Event

		Target string `json:"target"`
		From string `json:"from,omitempty"`
		Candidate interface{} `json:"candidate"`
	}
)