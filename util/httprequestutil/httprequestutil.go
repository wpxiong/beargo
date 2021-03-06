package httprequestutil

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "net/url"
  "mime/multipart"
  "strings"
  "encoding/json"
  "io/ioutil"
  "reflect"
)

func init() {
  log.InitLog()
}

func MapMerge(dst, src url.Values) *url.Values{
    var valuesMap url.Values = make(map[string][]string)
    for key, value := range dst {
        if valuesMap[key] == nil {
          valuesMap[key] = make([]string,0)
        }
        valuesMap[key] = append(valuesMap[key],value...)
    }
    for key, value := range src {
        if valuesMap[key] == nil {
          valuesMap[key] = make([]string,0)
        }
        valuesMap[key] = append(valuesMap[key],value...)
    }
    return &valuesMap
}

func processFileParameter(filesParameter map[string][]*multipart.FileHeader,appContext *appcontext.AppContext) {
   appContext.FileList = make(map[string][]multipart.File)
   for key,fileheaderlist := range filesParameter {
      fileList := make([]multipart.File,0)
      for _,fileheader := range fileheaderlist {
        if file,ok := fileheader.Open();ok == nil {
           fileList = append(fileList,file)
        }
      }
      appContext.FileList[key] = fileList
   }
}


func ProcessHttpRequestParam(appContext *appcontext.AppContext) {
    var getParameter url.Values = appContext.Request.HttpRequest.URL.Query()
    var formParameter url.Values
    var filesParameter map[string][]*multipart.FileHeader  
    contentType := strings.ToLower(appContext.Request.HttpRequest.Header.Get("Content-Type"))
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
      case "application/json" :
         body, err := ioutil.ReadAll(appContext.Request.HttpRequest.Body)
         if err != nil {
            log.Error(err)
         }
         var objmap map[string] interface{}
         err = json.Unmarshal(body, &objmap)
         if err != nil {
            log.Error("process application/json error")
         }else {
            for key,val := range objmap {
               appContext.Parameter[key] = val
            }
         }
    }
    paramMap := MapMerge(getParameter,formParameter)
    for key,value := range *paramMap {
        if len(value) == 1{
           appContext.Parameter[key] = value[0]
        }else {
           appContext.Parameter[key] = value
        }
    }
    pam := reflect.ValueOf(appContext.Parameter).Interface()
    ConvertMapKeyToLower(&pam)
    processFileParameter(filesParameter,appContext)
}


func ConvertMapKeyToLower(mapInfo *interface{}){
  switch reflect.TypeOf(*mapInfo).Kind(){
     case reflect.Map:
       newParamap := (*mapInfo).(map[string]interface{})
       for key,mapval := range newParamap {
           if mapval != nil {
              ConvertMapKeyToLower(&mapval)
              newParamap[strings.ToLower(key)] = mapval
              if key != strings.ToLower(key) {
                 delete(newParamap,key)
              }
           }
       }
    case reflect.Array:
       if mapInfo != nil {
          var arraylist []interface{} = (*mapInfo).([]interface{})
          for i,val2 := range arraylist {
             ConvertMapKeyToLower(&val2)
             arraylist[i] = val2
          }
       }
    default:
   }
}

func ParseGetParameter(app *appcontext.AppContext, urlParam url.Values){
    for key,value := range urlParam {
       if len(value) == 1{
           app.Parameter[strings.ToLower(key)] = value[0]
       }else {
           app.Parameter[strings.ToLower(key)] = value
       }
   }
}