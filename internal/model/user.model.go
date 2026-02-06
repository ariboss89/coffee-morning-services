package model

type Login struct {
	Id       int    `db:"password"`
	Email    string `db:"id"`
	Password string `db:"password"`
	Role     string `db:"role"`
}

type User struct {
	Fullname string `db:"fullname"`
	Avatar   string `db:"avatar"`
	Bio      string `db:"bio"`
}
