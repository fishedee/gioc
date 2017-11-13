# gioc

awesome golang ioc and aop!Only Hava 3 function,Very Very Easy To Use!

# Language

* [English](https://github.com/fishedee/gioc/blob/master/README.md)
* [简体中文](https://github.com/fishedee/gioc/blob/master/README_zh_cn.md)

# IOC

```go
type User struct {
	Id   int
	Name string
}
type UserAo interface {
	Get(id int) User
	Add(data User) int
}
type UserDb interface {
	Get(id int) User
	Add(data User) int
}
```

define two bean interface

```go
type userDbImpl struct {
}

func (this *userDbImpl) Get(id int) api.User {
	return api.User{Name:"Fish"}
}

func (this *userDbImpl) Add(data api.User) int {
	return 10001
}

func newUserDbImpl() api.UserDb {
	return &userDbImpl{}
}
func init() {
	gioc.Register(newUserDbImpl)
}
```

define UserDb implment and register it by gioc.Register.You should notice that constructor function argument is rely and return value is target

```go
type userAoImpl struct {
	userDb api.UserDb
}

func (this *userAoImpl) Get(id int) api.User {
	return this.userDb.Get(id)
}

func (this *userAoImpl) Add(data api.User) int {
	return this.userDb.Add(data)
}

func newUserAoImpl(userDb api.UserDb) api.UserAo {
	userAo := &userAoImpl{}
	userAo.userDb = userDb
	return userAo
}

func init() {
	gioc.Register(newUserAoImpl)
}
```

define UserAo implment and register it by gioc.Register.You should notice that UserAo need UserDb to construct

```go
userAo := gioc.New((*api.UserAo)(nil), nil, nil).(api.UserAo)
```

use gioc.New to get userAo.gioc will handle all rely,awesome!second argument and third argument is empty。

# STUB

```go
type UserDbStub struct {
}

func (this *UserDbStub) Get(id int) api.User {
	return api.User{Id: 10001, Name: "Fish"}
}

func (this *UserDbStub) Add(data api.User) int {
	return 10002
}

userAo := gioc.New((*api.UserAo)(nil), []interface{}{
	func() api.UserDb {
		return &UserDbStub{}
	},
}, nil).(api.UserAo)
```

use gioc.New to get userAo,but second argument represent which bean is replace by the new bean,this trick is used to unit test fake mock and stub.

# AOP

```go
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
```

if you want to use gioc aop,you have to code the general aop implement for interface.but no surprise, this suck code can auto generated by tool.notice that the aop implement is register by gioc.RegisterHook not gioc.Register

```go
func hookHandler(data interface{}) interface{} {
	dataValue := reflect.ValueOf(data)
	newValue := reflect.MakeFunc(dataValue.Type(), func(args []reflect.Value) []reflect.Value {
		fmt.Println("Hook Begin!")
		result := dataValue.Call(args)
		fmt.Println("Hook End!")
		return result
	})
	return newValue.Interface()
}
hook := map[string]interface{}{
	"UserAo.Get": hookHandler,
	"UserAo.Add": hookHandler,
}
userAo := gioc.New((*api.UserAo)(nil), nil, hook).(api.UserAo)
```

at last,use gioc.New to create bean,the third argument is hook setting.and now,we can never modify userAoImpl code to achieve aop!!

# TODO

* mutli thread support
* Register for singleton
* aop by regex expression like aspectJ and use complete package name
* more elegant register interface
* auto generate hook code

# License

gioc source code is licensed under the Apache Licence, Version 2.0 (http://www.apache.org/licenses/LICENSE-2.0.html)

