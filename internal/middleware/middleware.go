package middleware

type Role struct {
	Id    int    `db:"id"`
	Title string `db:"title"`
}
