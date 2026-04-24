package models

const (
	TypeSimpleUtterance = "SimpleUtterance"
)

type Request struct {
	Timezone string          `json:"timezone"`
	Request  SimpleUtterance `json:"request"`
	Session  Session         `json:"session"`
	Version  string          `json:"version"`
}

type Session struct {
	New  bool        `json:"new"`
	User RequestUser `json:"user"`
}

type RequestUser struct {
	UserID string `json:"user_id"`
}

type SimpleUtterance struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

type Response struct {
	Response ResponsePayload `json:"response"`
	Version  string          `json:"version"`
}

type ResponsePayload struct {
	Text string `json:"text"`
}
