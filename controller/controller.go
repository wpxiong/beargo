package controller

import (
  "../log"
  "../appcontext"
)


func init() {
  log.InitLog()
}


type Controller struct {
	data map[interface{}]interface{}
	controllerName string
	actionName     string
	methodMapping  map[string]func()
	ctx   *appcontext.AppContext
}

type ControllerMethod interface {
	Before()
	After()
	Get()
	Post()
	Delete()
	Put()
	Head()
}


func (con *Controller) getData() map[interface{}]interface{} {
   return con.data
}

func (con *Controller) getControllerName() string {
   return con.controllerName
}

func (con *Controller) getActionName() string {
   return con.actionName
}

func (con *Controller) getAppContext() *appcontext.AppContext {
   return con.ctx
}


