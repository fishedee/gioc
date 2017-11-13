package service

import (
	"github.com/fishedee/gioc"
	"github.com/fishedee/gioc/sample/api"
	"github.com/fishedee/gioc/sample/util"
)

type userDbImpl struct {
	db util.Db
}

func (this *userDbImpl) Get(id int) api.User {
	result := this.db.Select(id)
	if result == nil {
		panic("404 not found id")
	}
	return result.(api.User)
}

func (this *userDbImpl) Add(data api.User) int {
	return this.db.Insert(data)
}

func newUserDbImpl(db util.Db) api.UserDb {
	userDbImpl := &userDbImpl{}
	userDbImpl.db = db
	return userDbImpl
}

func init() {
	gioc.Register(newUserDbImpl)
}
