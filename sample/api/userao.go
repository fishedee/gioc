package api

import (
	"github.com/fishedee/gioc"
)

type UserAo interface {
	Get(id int) User
	Add(data User) int
}

type UserAoHook struct {
	GetHandler func(id int) User
	AddHandler func(data User) int
}

func (this *UserAoHook) Get(id int) User {
	return this.GetHandler(id)
}

func (this *UserAoHook) Add(data User) int {
	return this.AddHandler(data)
}

func NewUserAoHook() UserAo {
	return &UserAoHook{}
}

func init() {
	gioc.RegisterHook(NewUserAoHook)
}
