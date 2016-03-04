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
	Before(ctx   *appcontext.AppContext)
	After(ctx   *appcontext.AppContext)
	Get(ctx   *appcontext.AppContext)
	Post(ctx   *appcontext.AppContext)
	Delete(ctx   *appcontext.AppContext)
	Put(ctx   *appcontext.AppContext)
	Head(ctx   *appcontext.AppContext)
}


func (con *Controller) Before(ctx  *appcontext.AppContext) {

}

func (con *Controller) After(ctx  *appcontext.AppContext) {

}

func (con *Controller) Get(ctx  *appcontext.AppContext) {

}

func (con *Controller) Post(ctx  *appcontext.AppContext) {

}

func (con *Controller) Delete(ctx  *appcontext.AppContext) {

}

func (con *Controller) Put(ctx   *appcontext.AppContext) {

}

func (con *Controller) Head(ctx   *appcontext.AppContext) {

}





