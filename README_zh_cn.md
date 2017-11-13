# gioc

gioc 是一个纯粹的golang类spring实现，实现spring中的IOC(控制依赖)和AOP(面向切面编程)两方面。而且，不需要改动编译器，接口一共就只有3个，简单易用。

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

定义两个bean的接口

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

使用gioc注册UserDb接口的实现，构造器的参数表明实现的依赖，返回值表明实现的输出

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

使用gioc注册UserAo接口的视线，构造器显示需要依赖UserDb接口才能构造输出

```go
userAo := gioc.New((*api.UserAo)(nil), nil, nil).(api.UserAo)
```

使用gioc的New接口就能获得UserAo的实例，gioc会自动处理bean容器中的依赖问题

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

使用gioc的New接口时，第二个参数代表覆盖容器中的bean实现，用来实现单元测试中的mock和stub。

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

要实现gioc的AOP编程，需要用工具自动生成各个接口的Hook实现，以上是UserAo接口的通用Hook实现，然后用RegisterHook注入

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

然后用gioc的hook接口来生成bean，自动回调hookHandler，实现不改动UserAo实现的情况下注入hook

# TODO

* 多线程支持
* Register时支持单例模式
* Hook时用aspectJ的正则表达式，以及用完整的包名
* 注入STUB的参数不太漂亮
* 自动化的mock和hook代码生成

