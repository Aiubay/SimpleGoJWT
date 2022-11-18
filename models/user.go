package models

type User struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"Email" gorm:"unique"`
	Passwords []byte `json:"-"`
}
