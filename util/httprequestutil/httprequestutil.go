package httprequestutil

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "net/url"
)

func init() {
  log.InitLog()
}

func ProcessHttpRequestParam(appContext *appcontext.AppContext) {
    var getParameter url.Values = appContext.Request.HttpRequest.URL.Query()
    var formParameter url.Values
    contentType := appContext.Request.HttpRequest.Header.Get("Content-Type")
    switch contentType {
      case "application/x-www-form-urlencoded":
        if err := appContext.Request.HttpRequest.ParseForm(); err != nil {
          log.Error("process application/x-www-form-urlencoded error")
        } else {
          formParameter = appContext.Request.HttpRequest.Form
        }
      case "multipart/form-data":
        if err := appContext.Request.HttpRequest.ParseForm(); err != nil {
         log.Error("process multipart/form-data error")
        } else {
          formParameter = appContext.Request.HttpRequest.Form
        }
    }
    log.Debug(getParameter)
    log.Debug(formParameter)
}