package controller

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
)


func init() {
  log.InitLog()
}


type Controller struct {
}

type ControllerMethod interface {
	Before(ctx   *appcontext.AppContext,form interface{}) bool
	After(ctx   *appcontext.AppContext,form interface{}) bool
	Get(ctx   *appcontext.AppContext,form interface{})
	Post(ctx   *appcontext.AppContext,form interface{})
	Delete(ctx   *appcontext.AppContext,form interface{})
	Put(ctx   *appcontext.AppContext,form interface{})
	Head(ctx   *appcontext.AppContext,form interface{})
}


func (con *Controller) Before(ctx  *appcontext.AppContext,form interface{}) bool {
  log.Debug("Before Function Start")
  return true
}

func (con *Controller) After(ctx  *appcontext.AppContext,form interface{}) bool {
  log.Debug("After Function Start")
  return true
}

func (con *Controller) Get(ctx  *appcontext.AppContext,form interface{}) {

}

func (con *Controller) Post(ctx  *appcontext.AppContext,form interface{}) {

}

func (con *Controller) Delete(ctx  *appcontext.AppContext,form interface{}) {

}

func (con *Controller) Put(ctx   *appcontext.AppContext,form interface{}) {

}

func (con *Controller) Head(ctx   *appcontext.AppContext,form interface{}) {

}





