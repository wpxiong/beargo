package appcontext

import (
  "os"
  "bufio"
  "strings"
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/webhttp"
  "github.com/wpxiong/beargo/moudle"
  "github.com/wpxiong/beargo/constvalue"
  "reflect"
  "regexp"
  "mime/multipart"
)


func init() {
  log.InitLog()
}

type ConvertFunc func(string)(interface{})

type DBConnectionInfo struct {
   Dailect_Type  string
   DB_Name  string
   DB_Url   string
   DB_User  string
   DB_Pass  string
   DB_Session_Name string
}

type AppConfigContext struct {
   ConfigPath  string
   Port  int
   ConvertList  map[string] ConvertFunc
   ConfigParam  map[string]interface{}
   dbconfiglist []DBConnectionInfo
}
  
type AppContext struct {
  ConfigContext  *AppConfigContext
  Data map[interface{}]interface{}
  Parameter     map[string]interface{}
  Request        *webhttp.HttpRequest
  Writer         *webhttp.HttpResponse
  ControllerMethodInfo  *reflect.Method
  FormType          reflect.Type
  Form              interface{}
  UrlPath        string
  Renderinfo     interface{}
  redirect       bool
  RedirectPath   string
  FileList       map[string][]multipart.File
  Trans          map[string]* moudle.Trans
  DBSession      map[string]* moudle.Moudle
  ErrorInfo      map[string] []string
  RenderData     map[string] interface{}
  UseLayout    bool
  LayoutName   string
}

func (ctx *AppContext)  IsRedirect() bool{
   return ctx.redirect 
}

func (ctx *AppContext)  SetRedirect() {
   ctx.redirect = true
}

func (ctx *AppContext)  ClearRedirect() {
   ctx.redirect = false
}

func (ctx *AppContext) GetConfigValue(key string,defaultValue interface{}) interface{} {
   val := ctx.ConfigContext.ConfigParam[key]
   if val == nil {
      return defaultValue
   }
   switch reflect.TypeOf(val).Kind() {
      case reflect.String:
         return val
      case reflect.Slice:
        return val.([]string)
      default:
         return val
   }
}


func readLines(path string) (lines []string, err error) {
    var linesarray = make([]string,100) 
    filereader,error := os.Open(path)
    if error != nil {
       return make([]string,0),error
    }
    scanner := bufio.NewScanner(filereader)
	for scanner.Scan() {
	    linesarray = append(linesarray,scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.ErrorArray("config file Read Error",err)
		return make([]string,0),err
	}
	return linesarray,nil
}

func (ctx *AppConfigContext) LoadConfig() {
   pwd, _ := os.Getwd()
   ctx.dbconfiglist = make([]DBConnectionInfo,0)
   lines,error := readLines(pwd + ctx.ConfigPath)
   if error != nil {
      log.Error("read config file error!")
   }else{
     var currentDbConfig *DBConnectionInfo = nil
     for _,line := range lines {
       line = strings.Trim(line," ")
       if len(line) == 0 || strings.HasPrefix(line,"#"){
          continue
       }else {
          if line == "[" + constvalue.DB_CONFIG_SECTION + "]" {
            if currentDbConfig != nil {
               ctx.dbconfiglist = append(ctx.dbconfiglist,*currentDbConfig)
            }
            currentDbConfig = &DBConnectionInfo{}
          }
          words := strings.Split(line,"=")
          if len(words) >= 2 { 
            pamname,pamval := words[0],words[1]
            pamname = strings.Trim(pamname," ")
            pamval = strings.Trim(pamval," ")
            switch pamname {
             case constvalue.DB_DIALECT_TYPE:
                 currentDbConfig.Dailect_Type = pamval
                 continue
             case constvalue.DB_NAME:
                 currentDbConfig.DB_Name = pamval
                 continue
             case constvalue.DB_URL:
                 currentDbConfig.DB_Url = pamval
                 continue
             case constvalue.DB_USER:
                 currentDbConfig.DB_User = pamval
                 continue
             case constvalue.DB_PASSWORD:
                 currentDbConfig.DB_Pass = pamval
                 continue
             case constvalue.DB_SESSION_NAME:
                 currentDbConfig.DB_Session_Name = pamval
                 continue
             default:
            }
            var pamvalList []string 
            var isArray bool = true
            result,err := regexp.MatchString(`^\[.*\]$`,pamval)
            if !result || err != nil {
              isArray = false
            }else {
              pamvalList = strings.Split(pamval[1:len(pamval)-1],",")
            }
            if (pamname != ""){
              if ctx.ConfigParam[pamname] == nil {
                if !isArray{
                   ctx.ConfigParam[pamname] = pamval
                }else {
                   ctx.ConfigParam[pamname] = pamvalList
                }
              }else {
                switch  ctx.ConfigParam[pamname].(type) {
                  case string:
                     preval := ctx.ConfigParam[pamname].(string)
                     var list []string = make([]string,2)
                     list[0] = preval
                     if !isArray{
                       list[1] = pamval
                     }else {
                       list = append(list[:1],pamvalList...)
                     }
                     ctx.ConfigParam[pamname] = list
                  case []string:
                     var list []string = ctx.ConfigParam[pamname].([]string)
                     if !isArray{
                       list = append(list,pamval)
                     }else {
                       list = append(list,pamvalList...)
                     }
                     ctx.ConfigParam[pamname] = list
               }
             }
           }
         }
       }
     }
     if currentDbConfig != nil {
        ctx.dbconfiglist = append(ctx.dbconfiglist,*currentDbConfig)
     }
   }
}

func (ctx *AppContext) InitAppContext(ConfigPath string , Port int) {
   if ctx.ConfigContext == nil {
      ctx.ConfigContext = &AppConfigContext{ConfigPath :ConfigPath, Port:Port}
   }
   ctx.ConfigContext.ConvertList = make(map[string]ConvertFunc)
   ctx.ConfigContext.ConfigParam = make(map[string](interface{}))
   if ctx.ConfigContext.ConfigPath != ""{
      ctx.ConfigContext.LoadConfig()
      log.Info(ctx.ConfigContext.ConfigParam)
   }
   
}

func (ctx *AppContext) AddConvertFunctiont(paramType string,function ConvertFunc) {
    ctx.ConfigContext.ConvertList[paramType] = function
}

func (ctx *AppContext) CopyAppContext(frmctx *AppContext) {
   ctx.Data = make(map[interface{}]interface{})
   ctx.ConfigContext = frmctx.ConfigContext
   ctx.ErrorInfo = make(map[string][]string)
   ctx.RenderData = make(map[string]interface{})
   ctx.DBSession = frmctx.DBSession
}

func (ctx *AppContext) Convert(valStr string,valType string) interface{} {
    function := ctx.ConfigContext.ConvertList[valType]
    if function != nil {
       log.Debug(function)
       return function(valStr)
    }else {
       return valStr
    }
}

func (ctx *AppContext)  DestoryAppContext() {
   for _,fileheaderlist := range ctx.FileList {
      for _,file := range fileheaderlist {
        if err := file.Close(); err != nil {
           log.Error(err)
        }
      }
   }
}

func (ctx *AppContext) GetPostFileByParameterName(parameterName string) [] multipart.File{
   if val,ok := ctx.FileList[parameterName] ;ok {
      return val
   }else {
      return make([]multipart.File,0) 
   }
}

func (ctx *AppConfigContext) GetDBConfigParameter() []DBConnectionInfo {
   return ctx.dbconfiglist
}


func (ctx *AppContext) GetDefaultDB() *moudle.Moudle {
   for _,val :=  range ctx.DBSession {
      return val
   }
   return nil
}


func (ctx *AppContext) GetDBByName(dbname string) *moudle.Moudle {
   if val,ok := ctx.DBSession[dbname]; ok {
      return val
   }else{
      return nil
   }
}


func (ctx *AppContext) GetDefaultDBTransaction() *moudle.Trans {
   for _,val :=  range ctx.Trans {
      return val
   }
   return nil
}


func (ctx *AppContext) GetDBTranscationByName(dbname string) *moudle.Trans {
   if val,ok := ctx.Trans[dbname]; ok {
      return val
   }else{
      return nil
   }
}

func (ctx *AppContext) SetError(errorKey string,errorMessage string)  {
   if errorList,ok := ctx.ErrorInfo[errorKey];!ok {
      errorList = make([]string,1)
      errorList[0] = errorMessage
      ctx.ErrorInfo[errorKey] = errorList
   }else {
      errorList = append(errorList,errorMessage)
      ctx.ErrorInfo[errorKey] = errorList
   }
}

func (ctx *AppContext) SetRenderData(key string, renderdata interface{})  {
    ctx.RenderData[key] = renderdata
}

func (ctx *AppContext) ClearError(errorKey string)  {
   if _,ok := ctx.ErrorInfo[errorKey];ok {
      ctx.ErrorInfo[errorKey] = make([]string,0)
   }
}


func (ctx *AppContext) ClearAllError()  {
   ctx.ErrorInfo = make(map[string][]string)
}


func (ctx *AppContext) SetLayoutBaseName(baseName string )  {
   ctx.UseLayout = true
   ctx.LayoutName = baseName
}
