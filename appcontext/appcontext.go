package appcontext

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/webhttp"
)

func init() {
  log.InitLog()
}

type ConvertFunc func(string)(interface{})

type AppConfigContext struct {
   ConfigPath  string
   Port  int
   ConvertList map[string] ConvertFunc
}
  
type AppContext struct {
  ConfigContext  *AppConfigContext
  Data map[interface{}]interface{}
  Parameter     map[string]interface{}
  Request        *webhttp.HttpRequest
  Writer         *webhttp.HttpResponse
}

func (ctx *AppContext) InitAppContext(ConfigPath string , Port int) {
   if ctx.ConfigContext == nil {
      ctx.ConfigContext = &AppConfigContext{ConfigPath :ConfigPath, Port:Port}
   }
   ctx.ConfigContext.ConvertList = make(map[string]ConvertFunc)
}

func (ctx *AppContext) AddConvertFunctiont(paramType string,function ConvertFunc) {
    ctx.ConfigContext.ConvertList[paramType] = function
}

func (ctx *AppContext) CopyAppContext(frmctx *AppContext) {
   ctx.Data = make(map[interface{}]interface{})
   ctx.ConfigContext = frmctx.ConfigContext
}

func (ctx *AppContext) Convert(valStr string,valType string) interface{} {
    function := ctx.ConfigContext.ConvertList[valType]
    if function != nil {
       return function(valStr)
    }else {
       return valStr
    }
}

