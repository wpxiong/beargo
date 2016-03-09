package filter

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/render"
)

func init() {
  log.InitLog()
}



func RenderOutPutFilter(app *appcontext.AppContext) bool {
   log.Debug("RenderOutPutFilter Start")
   var renderInfo *render.RenderInfo = app.Renderinfo.(*render.RenderInfo)
   renderInfo.RenderProcess(renderInfo.OutPutData)
   return true
}