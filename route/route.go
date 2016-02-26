package route

import (
  "../controller"
  "../appcontext"
  "../http"
  "../log"
)

func init() {
  log.InitLog()
}

type RouteProcess struct {
   ctx   *appcontext.AppContext  
}

var routeProcess *RouteProcess
func NewRouteProcess(appContext *appcontext.AppContext) *RouteProcess {
   if routeProcess == nil {
     routeProcess = &RouteProcess{ctx : appContext}
   }
   return routeProcess
}


func (rtp *RouteProcess ) processRequest(request * http.HttpRequest) *RouteInfo {
   rinfo := &RouteInfo{}
   return rinfo
}

func (rtp *RouteProcess ) Add() {
   
}


type RouteInfo struct {
   requestUrl     string
   paramInfo      map[string]string
   result         bool
   controller     controller.Controller
   request        * http.HttpRequest
}

func (this *RouteInfo) getRequestInfo() string {
   return this.requestUrl
}

func (this *RouteInfo) getParamInfo() map[string]string {
   return this.paramInfo
}

func (this *RouteInfo) getMatchResult() bool {
   return this.result
}

func (this *RouteInfo) Init(requestUrl string, paramInfo map[string]string,result bool,controller controller.Controller,request * http.HttpRequest) {
   this.requestUrl = requestUrl
   this.paramInfo = paramInfo
   this.result = result
   this.controller = controller
   this.request = request
}

