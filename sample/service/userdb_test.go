package service

import (
	"github.com/fishedee/gioc"
	"github.com/fishedee/gioc/sample/api"
	"testing"
)

func TestUserDb(t *testing.T) {
	userDb := gioc.New((*api.UserDb)(nil), nil, nil).(api.UserDb)
	id1 := userDb.Add(api.User{Name: "Fish"})
	id2 := userDb.Add(api.User{Name: "Fish2"})
	if id1 != 10001 ||
		id2 != 10002 {
		t.Errorf("Error!")
	}
	user1 := userDb.Get(10001)
	user2 := userDb.Get(10002)
	if user1.Id != 10001 ||
		user1.Name != "Fish" ||
		user2.Id != 10002 ||
		user2.Name != "Fish2" {
		t.Errorf("Error2!")
	}
}
