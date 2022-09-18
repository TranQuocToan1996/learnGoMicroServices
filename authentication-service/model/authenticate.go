package model

type ReqAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
