package util

import (
	"github.com/fishedee/gioc"
	"reflect"
)

type Db interface {
	Insert(data interface{}) int
	Select(id int) interface{}
}

type dbImpl struct {
	totalId int
	data    map[int]interface{}
}

func (this *dbImpl) Select(id int) interface{} {
	return this.data[id]
}

func (this *dbImpl) Insert(data interface{}) int {
	this.totalId++
	origin := reflect.ValueOf(data)
	temp := reflect.New(origin.Type()).Elem()
	temp.Set(origin)
	temp.FieldByName("Id").Set(reflect.ValueOf(this.totalId))
	this.data[this.totalId] = temp.Interface()
	return this.totalId
}

func newDb() Db {
	dbImpl := &dbImpl{}
	dbImpl.totalId = 10000
	dbImpl.data = map[int]interface{}{}
	return dbImpl
}

func init() {
	gioc.Register(newDb)
}
