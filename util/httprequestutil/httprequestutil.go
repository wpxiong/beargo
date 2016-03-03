package httprequestutil

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "net/url"
  "mime/multipart"
)

func init() {
  log.InitLog()
}

func MapMerge(dst, src url.Values) *url.Values{
    var valuesMap url.Values = make(map[string][]string)
    for key, value := range dst {
        if valuesMap[key] == nil {
          valuesMap[key] = make([]string,len(value))
        }
        valuesMap[key] = append(valuesMap[key],value...)
    }
    for key, value := range src {
       if valuesMap[key] == nil {
         valuesMap[key] = value
       }else {
         var slice1 []string = valuesMap[key]
         valuesMap[key] = append(slice1,value...)
       }
    }
    return &valuesMap
}

func ProcessHttpRequestParam(appContext *appcontext.AppContext) {
    var getParameter url.Values = appContext.Request.HttpRequest.URL.Query()
    var formParameter url.Values
    var filesParameter map[string][]*multipart.FileHeader  
    contentType := appContext.Request.HttpRequest.Header.Get("Content-Type")
    switch contentType {
      case "application/x-www-form-urlencoded":
        if err := appContext.Request.HttpRequest.ParseForm(); err != nil {
          log.Error("process application/x-www-form-urlencoded error")
        } else {
          formParameter = appContext.Request.HttpRequest.Form
        }
      case "multipart/form-data":
        if err := appContext.Request.HttpRequest.ParseMultipartForm(32 << 20 /* 32 MB */); err != nil {
            log.Error("process multipart/form-data error")
        } else {
            formParameter = appContext.Request.HttpRequest.MultipartForm.Value
            filesParameter = appContext.Request.HttpRequest.MultipartForm.File
        }
    }
    formParameter = make(map[string][]string)
    formParameter["xiong"] = []string{"gag"}
    paramMap := MapMerge(getParameter,formParameter)
    for key,value := range *paramMap {
        if len(value) == 1{
           appContext.Parameter[key] = value[0]
        }else {
           appContext.Parameter[key] = value
        }
    }
    log.Debug(filesParameter)
    log.Debug(appContext.Parameter)
}