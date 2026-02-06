package dto

type User struct {
	Fullname string `db:"fullname"`
	Avatar   string `db:"avatar"`
	Bio      string `db:"bio"`
}
