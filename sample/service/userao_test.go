package service

import (
	"github.com/fishedee/gioc"
	"github.com/fishedee/gioc/sample/api"
	"testing"
)

func NewUserDbStub() api.UserAo {
	return api.UserAo{
		Get: func(id int) api.User {
			return api.User{Id: 10001, Name: "Fish"}
		},
		Add: func(data api.User) int {
			return 10002
		},
	}
}

func TestUserAoGet(t *testing.T) {
	userAo := gioc.New(api.UserAo{}, map[interface{}]interface{}{
		&api.UserAo{}: NewUserDbStub,
	}, nil).(api.UserAo)
	left := userAo.Get(0)
	right := api.User{Id: 10001, Name: "Fish"}
	if left.Id != right.Id ||
		left.Name != right.Name {
		t.Errorf("Error!")
	}
}

func TestUserAoAdd(t *testing.T) {
	userAo := gioc.New(api.UserAo{}, map[interface{}]interface{}{
		api.UserAo{}: NewUserDbStub,
	}, nil).(api.UserAo)
	left := userAo.Add(api.User{})
	right := 10002
	if left != right {
		t.Errorf("Error!")
	}
}
