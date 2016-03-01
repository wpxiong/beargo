package controller

import (
  "github.com/beargo/log"
  "github.com/beargo/appcontext"
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

func (con *Controller) Before() {

}

func (con *Controller) After() {

}

func (con *Controller) Get() {

}

func (con *Controller) Post() {

}

func (con *Controller) Delete() {

}

func (con *Controller) Put() {

}

func (con *Controller) Head() {

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


