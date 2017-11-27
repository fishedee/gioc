package api

type User struct {
	Id   int
	Name string
}

type UserDb struct {
	Get func(id int) User
	Add func(data User) int
}
