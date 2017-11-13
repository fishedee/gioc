package api

type User struct {
	Id   int
	Name string
}

type UserDb interface {
	Get(id int) User
	Add(data User) int
}
