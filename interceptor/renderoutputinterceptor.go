package interceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/render"
)

func init() {
  log.InitLog()
}



func RenderOutPutinterceptor(app *appcontext.AppContext) bool {
   log.Debug("RenderOutPutinterceptor Start")
   var renderInfo *render.RenderInfo = app.Renderinfo.(*render.RenderInfo)
   renderInfo.UseLayout =  app.UseLayout
   renderInfo.LayoutName =  app.LayoutName
   renderInfo.RenderProcess()
   return true
}