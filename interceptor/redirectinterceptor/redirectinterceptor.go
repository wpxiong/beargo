package redirectinterceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/util"
  "github.com/wpxiong/beargo/util/httprequestutil"
  "github.com/wpxiong/beargo/constvalue"
  "github.com/wpxiong/beargo/route"
  "strings"
  "net/url"
)

func init() {
  log.InitLog()
}

func Redirectinterceptor(app *appcontext.AppContext) bool {
   log.Debug("Redirectinterceptor Start")
   if app.IsRedirect() {
     u, err := url.Parse(app.RedirectPath)
     if err != nil {
        log.Error("Redirect URL Error")
        util.Redirect(constvalue.REDIRECT_ERROR)
        return false
     }
     q := u.Query()
     httprequestutil.ParseGetParameter(app,q)
     app.ClearRedirect()
     urlArray := strings.Split(app.RedirectPath,"?")
     rinfo := &route.RouteInfo{UrlParamInfo: []route.ParaInfo{}}
     res := route.GetRouteProcess().UrlRoute(urlArray[0],rinfo)
     if res {
        app.FormType = rinfo.GetFormType()
        app.ControllerMethodInfo = rinfo.GetMethodInfo()
        app.UrlPath = rinfo.UrlPath
        rinfo.CallRedirectMethod(app)
     }else {
        log.Error("Redirect URL Error: " + urlArray[0])
     }
     return false
   }
   return true
}