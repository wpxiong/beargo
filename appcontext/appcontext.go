package appcontext

import (
  "os"
  "bufio"
  "strings"
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/webhttp"
  "reflect"
  "regexp"
)


func init() {
  log.InitLog()
}

type ConvertFunc func(string)(interface{})

type AppConfigContext struct {
   ConfigPath  string
   Port  int
   ConvertList  map[string] ConvertFunc
   ConfigParam  map[string]interface{}
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
   lines,error := readLines(pwd + ctx.ConfigPath)
   if error != nil {
      log.Error("read config file error!")
   }else{
      for _,line := range lines {
        line = strings.Trim(line," ")
        if len(line) == 0 || strings.HasPrefix(line,"#"){
           continue
        }
        words := strings.Split(line,"=")
        if len(words) >= 2 { 
           pamname,pamval := words[0],words[1]
           pamname = strings.Trim(pamname," ")
           var pamvalList []string 
           var isArray bool = true
           pamval = strings.Trim(pamval," ")
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
                     list = append(list,preval)
                     if !isArray{
                       list = append(list,pamval)
                     }else {
                       list = append(list,pamvalList...)
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

