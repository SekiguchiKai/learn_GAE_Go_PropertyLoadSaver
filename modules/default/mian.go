package main

import (
	"github.com/gin-gonic/gin"
	"github.com/SekiguchiKai/learn_GAE_Go_PropertyLoadSaver/api"
	"net/http"
)

const _APIPath = "/api"

func init() {
	g := gin.New()
	initAPI(g)
	// gin.New()の戻り値のEngineは、ServeHTTP(ResponseWriter, *Request)メソッドを持っているので、
	// type Handler interfaceを満た
	http.Handle("/", g)
}



func initAPI(g *gin.Engine){
	apiGin := g.Group(_APIPath)
	api.InitUserAPI(apiGin)
}
