package filter

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/render"
  "reflect"
)

func init() {
  log.InitLog()
}


func RenderFilter(app *appcontext.AppContext) bool {
   log.Debug("RenderFilter Start")
   var renderInfo *render.RenderInfo = render.CreateRenderInfo(app)
   var outputData interface{} = reflect.ValueOf([]string{"xxxx"}).Interface()
   renderInfo.InitRenderInfo(app)
   renderInfo.RenderProcess(outputData)
   return true
}