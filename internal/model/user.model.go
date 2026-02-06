package model

type User struct {
	Id       int    `db:"password"`
	Email    string `db:"id"`
	Password string `db:"password"`
	Role     string `db:"role"`
}
