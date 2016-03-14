package redirectfilter

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/util"
  "github.com/wpxiong/beargo/constvalue"
  "github.com/wpxiong/beargo/util/httprequestutil"
  "github.com/wpxiong/beargo/route"
  "github.com/wpxiong/beargo/filter"
  "strings"
  "net/url"
)

func init() {
  log.InitLog()
}

func RedirectFilter(app *appcontext.AppContext) bool {
   log.Debug("RedirectFilter Start")
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
     app.UrlPath = app.RedirectPath
     urlArray := strings.Split(app.UrlPath,"?")
     urlArray := strings.Split(app.UrlPath,"?")
     rinfo := &RouteInfo{UrlParamInfo: []ParaInfo{}}
     res := route.routeProcess.urlRoute(urlArray[0],rinfo)
     if res {
        rinfo.CallRedirectMethod(app)
     }else {
        log.Error("Redirect URL Error: " + urlArray[0])
     }
     return false
   }
   return true
}