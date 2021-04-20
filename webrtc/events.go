package webrtc

const (
	EventTypeJoinRoom = "join_room"
	EventTypeIceCandidate = "ice-candidate"
	EventTypeAnswer = "answer"
	EventTypeOffer = "offer"
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
		OtherUsers *Room `json:"other_users,omitempty"`
		UserID string `json:"user_id"`
	}

	EventUserJoined struct {
		Event

		UserID string `json:"user_id"`
	}

	EventOffer struct {
		Event

		Target string `json:"target"`
	}

	EventAnswer struct {
		Event

		Target string `json:"target"`
	}

	EventIceCandidate struct {
		Event

		Target string `json:"target"`
	}
)