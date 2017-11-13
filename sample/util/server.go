package util

import (
	"net/http"
)

type internal struct {
	internalServer *Server
}

func (this *internal) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	this.internalServer.handle(writer, request)
}

type ServerHandler func(http.ResponseWriter, *http.Request)
type MiddleHandler func([]interface{}) interface{}

type Server struct {
	router        map[string]ServerHandler
	middleHandler []MiddleHandler
	draftRouter   map[string]interface{}
}

func (this *Server) Route(path string, handler interface{}) {
	this.draftRouter[path] = handler
}

func (this *Server) Use(middleHandler MiddleHandler) {
	this.middleHandler = append(this.middleHandler, middleHandler)
}

func (this *Server) handle(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	pathHandler, isExist := this.router[path]
	if isExist == false {
		writer.WriteHeader(404)
		writer.Write([]byte("File Not Found!by fishedee"))
		return
	}
	pathHandler(writer, request)
}

func (this *Server) getSingleHandler(handler interface{}) ServerHandler {
	allHandler := []interface{}{handler}
	for i := len(this.middleHandler) - 1; i >= 0; i-- {
		curHandler := this.middleHandler[i](allHandler)
		allHandler = append(allHandler, curHandler)
	}
	resultHandler := allHandler[len(allHandler)-1]
	return ServerHandler(resultHandler.(func(http.ResponseWriter, *http.Request)))
}

func (this *Server) Run(addr string) {
	for key, singleHandler := range this.draftRouter {
		this.router[key] = this.getSingleHandler(singleHandler)
	}

	a := &internal{internalServer: this}
	err := http.ListenAndServe(addr, a)
	if err != nil {
		panic(err)
	}
}

func NewServer() *Server {
	result := &Server{}
	result.router = map[string]ServerHandler{}
	result.middleHandler = []MiddleHandler{}
	result.draftRouter = map[string]interface{}{}
	return result
}
