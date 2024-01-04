package domain

type User struct {
	Id       int64
	Email    string
	Password string
	Profile
}

type Profile struct {
	Nickname     string `json:"nickname"`
	Birth        string `json:"birth"`
	Introduction string `json:"introduction"`
}
