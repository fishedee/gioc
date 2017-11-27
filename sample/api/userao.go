package api

type UserAo struct {
	Get func (id int) User
	Add func(data User) int
}
