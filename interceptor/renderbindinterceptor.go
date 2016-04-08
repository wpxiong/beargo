package interceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/render"
)

func init() {
  log.InitLog()
}


func RenderBindinterceptor(app *appcontext.AppContext) bool {
   log.Debug("RenderBindinterceptor Start")
   var renderInfo *render.RenderInfo = render.CreateRenderInfo(app)
   renderInfo.InitRenderInfo(app)
   renderInfo.ErrorInfo = app.ErrorInfo
   renderInfo.OutPutData = app.Form
   app.Renderinfo = renderInfo
   return true
}