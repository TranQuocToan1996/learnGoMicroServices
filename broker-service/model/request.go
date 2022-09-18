package model

type Request struct {
	Action Action      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Action string

const (
	auth Action = "auth"
)
