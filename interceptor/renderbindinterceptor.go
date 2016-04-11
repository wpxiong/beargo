package interceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/render"
  "github.com/wpxiong/beargo/session"
  "github.com/wpxiong/beargo/constvalue"
  "reflect"
)

func init() {
  log.InitLog()
}

func makeOutputData(app *appcontext.AppContext) map[string] interface{} {
   output := make(map[string] interface{})
   fieldNum := app.FormType.NumField()
   for i:= 0 ;i<fieldNum ;i++ {
      name :=  app.FormType.Field(i).Name
      field := reflect.ValueOf(app.Form).Elem().Field(i)
      output[name] = field.Interface()
   }
   for key,val := range app.ErrorInfo {
      output[key] = val
   }
   for key,val := range app.RenderData {
      output[key] = val
   }
   
   request := app.Request.HttpRequest
   response := app.Writer.HttpResponseWriter
   var sess session.Session = session.NewSession(request , *response)
   output[constvalue.APP_SESSION_DATA] = sess.SessionValue
   log.Debug(output)
   return output
}

func RenderBindinterceptor(app *appcontext.AppContext) bool {
   log.Debug("RenderBindinterceptor Start")
   var renderInfo *render.RenderInfo = render.CreateRenderInfo(app)
   renderInfo.InitRenderInfo(app)
   renderInfo.ErrorInfo = app.ErrorInfo
   renderInfo.OutPutData = makeOutputData(app)
   app.Renderinfo = renderInfo
   return true
}