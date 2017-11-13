package gioc

import (
	"reflect"
)

type typeInfo struct {
	depType []reflect.Type
	builder reflect.Value
}

var typeBuilder map[reflect.Type]typeInfo = map[reflect.Type]typeInfo{}
var hookTypeBuilder map[reflect.Type]reflect.Value = map[reflect.Type]reflect.Value{}

func getRealType(t reflect.Type, hook map[string]interface{}, visit map[reflect.Type]bool, cache map[reflect.Type]reflect.Value, myTypeBuilder map[reflect.Type]typeInfo) reflect.Value {
	result, isExist := cache[t]
	if isExist {
		return result
	}
	_, isVisit := visit[t]
	if isVisit {
		panic("loop dependence")
	}
	visit[t] = true

	info, isExist := myTypeBuilder[t]
	if isExist == false {
		panic("unknown type")
	}
	args := []reflect.Value{}
	for _, singleDepType := range info.depType {
		args = append(args, dfs(singleDepType, hook, visit, cache, myTypeBuilder))
	}
	lastResult := info.builder.Call(args)
	cache[t] = lastResult[0]
	return lastResult[0]
}
func dfs(t reflect.Type, hook map[string]interface{}, visit map[reflect.Type]bool, cache map[reflect.Type]reflect.Value, myTypeBuilder map[reflect.Type]typeInfo) reflect.Value {
	realResult := getRealType(t, hook, visit, cache, myTypeBuilder)
	typeName := realResult.Type()
	hookFun := map[string]interface{}{}
	for i := 0; i != typeName.NumMethod(); i++ {
		methodName := typeName.Method(i).Name
		allName := typeName.Name() + "." + methodName
		hookHandler, isExist := hook[allName]
		if isExist {
			hookFun[methodName] = hookHandler
		}
	}
	if len(hookFun) == 0 {
		return realResult
	}
	myHookTypeBuilder, isExist := hookTypeBuilder[t]
	if isExist == false {
		panic("need hook but no hook builder")
	}
	hookResultList := myHookTypeBuilder.Call([]reflect.Value{})
	hookResult := hookResultList[0]
	for i := 0; i != hookResult.Type().NumMethod(); i++ {
		methodName := hookResult.Type().Method(i).Name
		fieldName := methodName + "Handler"
		hookHandler, isExist := hookFun[methodName]
		if isExist == false {
			hookResult.Elem().Elem().FieldByName(fieldName).Set(realResult.MethodByName(methodName))
		} else {
			origin := realResult.MethodByName(methodName).Interface()
			hookHandlerH := hookHandler.(func(interface{}) interface{})
			target := reflect.ValueOf(hookHandlerH(origin))
			hookResult.Elem().Elem().FieldByName(fieldName).Set(target)
		}
	}
	return hookResult
}

func New(a interface{}, moc []interface{}, hook map[string]interface{}) interface{} {
	myTypeBuilder := map[reflect.Type]typeInfo{}
	for key, value := range typeBuilder {
		myTypeBuilder[key] = value
	}

	if moc != nil && len(moc) != 0 {
		for _, singleMoc := range moc {
			a, b := getRegisterInfo(singleMoc)
			myTypeBuilder[a] = b
		}

	}

	if hook == nil {
		hook = map[string]interface{}{}
	}

	visit := map[reflect.Type]bool{}
	cache := map[reflect.Type]reflect.Value{}
	targetType := reflect.ValueOf(a).Type()
	if targetType.Elem().Kind() == reflect.Interface {
		targetType = targetType.Elem()
	}
	return dfs(targetType, hook, visit, cache, myTypeBuilder).Interface()
}

func getRegisterInfo(createFun interface{}) (reflect.Type, typeInfo) {
	typeValue := reflect.ValueOf(createFun)
	typeType := typeValue.Type()
	if typeType.Kind() != reflect.Func {
		panic("invalid type")
	}
	numIn := []reflect.Type{}
	for i := 0; i != typeType.NumIn(); i++ {
		numIn = append(numIn, typeType.In(i))
	}
	if typeType.NumOut() != 1 {
		panic("invalid num out")
	}
	numOut := typeType.Out(0)
	return numOut, typeInfo{
		depType: numIn,
		builder: typeValue,
	}
}
func Register(createFun interface{}) {
	a, b := getRegisterInfo(createFun)
	typeBuilder[a] = b
}

func RegisterHook(createFun interface{}) {
	a, b := getRegisterInfo(createFun)
	hookTypeBuilder[a] = b.builder
}
