package webapp

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/process"
  "github.com/wpxiong/beargo/route"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/webhttp"
  "github.com/wpxiong/beargo/controller"
  "github.com/wpxiong/beargo/util"
  "strconv"
  "net/http"
)

func init() {
  log.InitLog()
}


type  WebApplication struct {
   WorkProcess *process.WorkProcess
   RouteProcess  *route.RouteProcess
   AppContext *appcontext.AppContext
   IsStart bool
   control  chan int
}

var webApp *WebApplication

func InitDefaultConvertFunction(appContext *appcontext.AppContext){
  appContext.InitAppContext(appContext.ConfigContext.ConfigPath ,appContext.ConfigContext.Port)
  appContext.AddConvertFunctiont("Int",util.StringToInt)
  appContext.AddConvertFunctiont("Float",util.StringToFloat)
  appContext.AddConvertFunctiont("Double",util.StringToDouble)
  appContext.AddConvertFunctiont("Bool",util.StringToBool)
}

func New(appContext *appcontext.AppContext) *WebApplication {
   if webApp == nil {
      webApp = &WebApplication{WorkProcess : process.New(),RouteProcess : route.NewRouteProcess(appContext) , AppContext : appContext , control:make(chan int ) }
      InitDefaultConvertFunction(appContext)
   }
   return webApp;
}

func processWebRequest(param interface{}) interface{} {
   var rti *route.RouteInfo = param.(*route.RouteInfo)
   rti.CallMethod()
   rti.ResultChan <- 1
   return true
}

func processRequest(w http.ResponseWriter, r *http.Request){
    path := r.URL.Path
    log.Info(path)
    request := webhttp.HttpRequest{Urlpath : path, HttpRequest : r }
    response := webhttp.HttpResponse{HttpResponseWriter:&w}
    var rti *route.RouteInfo
    rti = webApp.RouteProcess.ProcessRequest(&request)
    if rti.GetResult() {
       rti.Request = &request
       rti.Writer =&response
       workjob := &process.WorkJob{Parameter : rti }
       workjob.WorkProcess = processWebRequest
       process.AddJob(workjob)
       _ = <- rti.ResultChan
        
    }else {
       log.Error("Error: not found page")
    }
}

func startProcess(web *WebApplication){
    err := http.ListenAndServe(":" + strconv.Itoa(web.AppContext.ConfigContext.Port) ,nil)
    if err != nil {
        log.DebugNoReturn("ListenAndServe: ")
        log.ErrorArray("Error",err)
        web.control <- 1
    }
}

func (web *WebApplication) AddRoute(urlPattern string,controller controller.ControllerMethod,method string) {
   web.RouteProcess.Add(urlPattern,controller,method)
}

func (web *WebApplication) Start() {
    http.HandleFunc("/", processRequest) 
    go startProcess(web)
    web.WorkProcess.Init_Default()
    res := <- web.control
    if res == 1 {
       process.StopWork()
       log.Debug("Stop WebApplication")
    }
}