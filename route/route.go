package route

import (
  "strings"
  "regexp"
  "reflect"
  "github.com/wpxiong/beargo/controller"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/webhttp"
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/interceptor"
  "github.com/wpxiong/beargo/constvalue"
  "github.com/wpxiong/beargo/render"
)

func init() {
  log.InitLog()
}

type ParaInfo struct {
    ParaName string
    ParaType string
    ParaValue string
}

const (
  DEFAULT_URL = "/index"
)

type TreeNode struct {
    parent     *TreeNode
    left       *TreeNode
    right      *TreeNode
    left_children   map[string]*TreeNode
    right_children  map[string]*TreeNode
    isLeafNode bool
    isReg      bool
    nodeValue  string
    controller controller.ControllerMethod
    funcmap    reflect.Value
    methodInfo reflect.Method
    paraInfo   []ParaInfo
    regstr     *regexp.Regexp
    formType       reflect.Type
    UrlPath        string
}


type RouteProcess struct {
   ctx       *appcontext.AppContext
   treeNode  *TreeNode
}


type RouteInfo struct {
   requestUrl     string
   UrlParamInfo   []ParaInfo
   result         bool
   controller     controller.ControllerMethod
   Request        *webhttp.HttpRequest
   Writer         *webhttp.HttpResponse
   funcmap        *reflect.Value
   methodInfo     *reflect.Method
   methodName     string
   ResultChan     chan int
   formType       reflect.Type
   UrlPath        string
}

func initRootTreeNode() *TreeNode {
   node := &TreeNode{}
   node.left_children = make(map[string]*TreeNode)
   node.right_children = make(map[string]*TreeNode)
   return node
}



var routeProcess *RouteProcess
func NewRouteProcess(appContext *appcontext.AppContext) *RouteProcess {
   if routeProcess == nil {
     routeProcess = &RouteProcess{ctx : appContext}
     treeNode := initRootTreeNode()
     treeNode.nodeValue = "/"
     routeProcess.treeNode = treeNode
   }
   return routeProcess
}

func (rtp *RouteInfo ) getFuncmap() reflect.Value {
   return (*rtp.funcmap)
}

func (rtp *RouteInfo ) GetMethodInfo() *reflect.Method {
   return rtp.methodInfo
}

func (rtp *RouteInfo ) GetFormType() reflect.Type {
   return  rtp.formType
}

func (rtp *RouteInfo ) GetResult() bool {
   return rtp.result
}

func (rtp *RouteInfo ) InitAppContext(app *appcontext.AppContext) {
   app.Request = rtp.Request
   app.Writer = rtp.Writer
   app.UrlPath = rtp.UrlPath
   app.Parameter = make(map[string]interface{})
   for _,param := range rtp.UrlParamInfo {
       app.Parameter[param.ParaName] = app.Convert(param.ParaValue,param.ParaType)
   }
}


func (rtp *RouteInfo ) ResourceClean(appContext *appcontext.AppContext) {
  log.Debug("ResourceClean")
}


func (rtp *RouteInfo ) CallRedirectMethod(appContext *appcontext.AppContext){
    var funcmap reflect.Value = rtp.getFuncmap()
    res := interceptor.ProcessRedirectBeforeinterceptor(appContext)
    if !res {
       return 
    }
    v := make([]reflect.Value,2)
    v[0] = reflect.ValueOf(appContext)
    v[1] = reflect.ValueOf(appContext.Form)
    defer func() {
        if err := recover(); err != nil {
            log.Error("Call Controller Method Error")
            res = interceptor.ProcessAfterinterceptor(appContext)
        }
    }()
    
    beforeFunc := reflect.ValueOf(rtp.controller).MethodByName(constvalue.BEFORE_FUNC)
    afterFunc := reflect.ValueOf(rtp.controller).MethodByName(constvalue.AFTER_FUNC)
    var result []reflect.Value = beforeFunc.Call(v)

    if (result[0].Interface()).(bool) == false {
      return 
    }
    
    funcmap.Call(v)
    result = afterFunc.Call(v)
    if (result[0].Interface()).(bool) == false {
      return 
    }
    res = interceptor.ProcessAfterinterceptor(appContext)
    if !res {
       return 
    }
}


func (rtp *RouteInfo ) CallMethod() {
    var funcmap reflect.Value = rtp.getFuncmap()
    var appContext *appcontext.AppContext = &appcontext.AppContext{ControllerMethodInfo :rtp.methodInfo,FormType:rtp.formType}
    appContext.CopyAppContext(routeProcess.ctx)
    rtp.InitAppContext(appContext)
    res := interceptor.ProcessBeforeinterceptor(appContext)
    defer rtp.ResourceClean(appContext)
    if !res {
       return 
    }
    v := make([]reflect.Value,2)
    v[0] = reflect.ValueOf(appContext)
    v[1] = reflect.ValueOf(appContext.Form)
    defer func() {
        if err := recover(); err != nil {
            log.Error("Call Controller Method Error")
            //500 Error
            RedirectTo500(appContext)
        }
        appContext.DestoryAppContext()
    }()
    
    beforeFunc := reflect.ValueOf(rtp.controller).MethodByName(constvalue.BEFORE_FUNC)
    afterFunc := reflect.ValueOf(rtp.controller).MethodByName(constvalue.AFTER_FUNC)
    var result []reflect.Value = beforeFunc.Call(v)

    if (result[0].Interface()).(bool) == false {
      return 
    }
    
    funcmap.Call(v)
    result = afterFunc.Call(v)

    if (result[0].Interface()).(bool) == false {
      return 
    }
    
    res = interceptor.ProcessAfterinterceptor(appContext)
    if !res {
       return 
    }
}


func (rtp *RouteProcess ) match(urlcom []string,index int , treeNodemap map[string]*TreeNode, paraList *[][]ParaInfo) (bool, controller.ControllerMethod,*reflect.Value,*reflect.Method,reflect.Type,string) {
   *paraList = (append(*paraList,[]ParaInfo{}))
   if index >= len(urlcom) {
      return false,nil,nil,nil,nil,""
   }
   var res bool = false
   var funcm *reflect.Value
   var methodInfo *reflect.Method
   var formType  reflect.Type
   var urlpath string
   for _,treeNode := range treeNodemap {
      if !treeNode.isReg {
          if urlcom[index] == treeNode.nodeValue {
              if index == len(urlcom) -1 {
                 if !treeNode.isLeafNode {
                    return false,nil,nil,nil,nil,""
                 }else {
                    return true,treeNode.controller,&treeNode.funcmap,&treeNode.methodInfo,treeNode.formType,treeNode.UrlPath
                 }
              } else {
                 var control  controller.ControllerMethod
                 res,control,funcm,methodInfo,formType,urlpath = rtp.match(urlcom,index + 1,treeNode.left_children,paraList)
                 if !res {
                   *paraList = (*paraList)[:len(*paraList)-1]
                   res,control,funcm,methodInfo,formType,urlpath = rtp.match(urlcom,index + 1,treeNode.right_children,paraList)
                   return res,control,funcm,methodInfo,formType,urlpath
                 }else {
                   return res,control,funcm,methodInfo,formType,urlpath
                 }
              }
          }
      }else {
          res = treeNode.regstr.MatchString(urlcom[index])
          if res {
              reglist := treeNode.regstr.FindStringSubmatch(urlcom[index])
              pa := []ParaInfo{}
              for i,reval := range reglist {
                 if i >0 {
                    pa = append(pa,ParaInfo{ParaName : treeNode.paraInfo[i-1].ParaName, ParaType : treeNode.paraInfo[i-1].ParaType, ParaValue: reval})
                 }
              }
              (*paraList)[len(*paraList)-1] = pa
              if index == len(urlcom) -1 {
                 if !treeNode.isLeafNode {
                    return false,nil,nil,nil,nil,""
                 }else {
                    return true,treeNode.controller,&treeNode.funcmap,&treeNode.methodInfo,treeNode.formType,treeNode.UrlPath
                 }
              } else {
                 var control  controller.ControllerMethod
                 res,control,funcm,methodInfo,formType,urlpath = rtp.match(urlcom,index + 1,treeNode.left_children,paraList)
                 if !res {
                   *paraList = (*paraList)[:len(*paraList)-1]
                   res,control,funcm,methodInfo,formType,urlpath = rtp.match(urlcom,index + 1,treeNode.right_children,paraList)
                   return res,control,funcm,methodInfo,formType,urlpath
                 }else {
                   return res,control,funcm,methodInfo,formType,urlpath
                 }
              }
          }
      }
     
   }
   return res,nil,nil,nil,nil,""
}


func (rtp *RouteProcess ) UrlRoute (request string, rtinfo *RouteInfo) bool{
   request = strings.Trim(request," ")
   if len(request) == 0 {
      return false
   }
   componentArray := strings.Split(request,"/")
   var startIndex int
   for index,urlCom := range componentArray {
       if len(urlCom) != 0 {
          startIndex = index 
          break
       }
   }
   componentArray = componentArray[startIndex:]
   var res bool
   var paraList [][]ParaInfo
   var controller controller.ControllerMethod
   var funcmap *reflect.Value
   var methodInfo *reflect.Method
   var formType reflect.Type
   var urlpath string
   res,controller,funcmap,methodInfo,formType,urlpath  = rtp.match(componentArray,0,rtp.treeNode.left_children,&paraList)
   if !res {
      res,controller,funcmap,methodInfo,formType,urlpath  = rtp.match(componentArray,0,rtp.treeNode.right_children,&paraList)
   }
   if res {
     rtinfo.controller = controller
     for _,pamlist := range paraList{
        for _,pam := range pamlist{
           rtinfo.UrlParamInfo = append(rtinfo.UrlParamInfo,pam)
        }
     }
   }
   rtinfo.funcmap = funcmap
   rtinfo.methodInfo = methodInfo
   rtinfo.formType = formType
   rtinfo.UrlPath = urlpath
   return res
} 

func (rtp *RouteProcess ) ProcessRequest(request * webhttp.HttpRequest) *RouteInfo {
   urlArray := strings.Split(request.Urlpath,"?")
   rinfo := &RouteInfo{UrlParamInfo: []ParaInfo{}}
   rinfo.result = false
   length := len(urlArray)
   if length > 0 {
      urlpath := urlArray[0]
      rinfo.requestUrl = urlpath
      res := rtp.UrlRoute(urlpath,rinfo)
      if res {
          rinfo.result = true
          rinfo.ResultChan = make(chan int)
      }
   }
   return rinfo
}

func (node *TreeNode ) debugInfo() {
    log.DebugNoReturn("{nodeValue:" + node.nodeValue + ",left_children:[")
    var k int = 0
    for _,nodep := range node.left_children {
       if nodep != nil {
         nodep.debugInfo()
         if k != len(node.left_children) -1 {
           log.DebugNoReturn(",")
         }
       }
       k++
    }
    k=0
    log.DebugNoReturn("],right_children:[")
    for _,nodep := range node.right_children {
       if nodep != nil {
         nodep.debugInfo()
         if k != len(node.right_children) -1 {
           log.DebugNoReturn(",")
         }
       }
       k++
    }
    log.DebugNoReturn("]")
    log.DebugNoReturn("}\n")
}

func (routeInfo *RouteInfo ) DebugInfo() {
   log.DebugNoReturn("{result:")
   log.DebugNoReturn(routeInfo.result)
   log.DebugNoReturn(",UrlParamInfo:[")
   for k,pam := range routeInfo.UrlParamInfo {
      log.DebugNoReturn(pam)
      if k != len(routeInfo.UrlParamInfo) -1 {
        log.DebugNoReturn(",")
      }
   }
   log.DebugNoReturn("],methodName:" + routeInfo.methodName)
   log.DebugNoReturn("}\n")
}

func (rtp *RouteProcess ) DebugInfo() {
  rtp.treeNode.debugInfo()
}


func (rtp *RouteProcess ) paramTypeCheck(typestr string) (string,string) {
    types := strings.ToLower(typestr)
    switch types {
      case "int" :
        return "int",`([-+]?[0-9]+)`
      case "double" :
        return "double",`([-+]?[0-9]*\\.?[0-9]+)`
      case "float" :
        return "float",`([-+]?[0-9]*\\.?[0-9]+)`
      case "string" :
        return "string",`(.*)`
      default:
        str := `(` + typestr + `)`
        return "reg",str
    }
}


func (rtp *RouteProcess ) checkMethod (controller controller.ControllerMethod,method string) (*reflect.Value,*reflect.Method) {
  controllerType := reflect.TypeOf(controller)
  for i := 0; i < controllerType.NumMethod(); i++ {
     methodstr := controllerType.Method(i)
     methodName := methodstr.Name
     firstLetter := string(methodName[0])
     if  methodstr.Type.NumIn() == 3 && firstLetter == strings.ToUpper(firstLetter) && strings.ToLower(methodName) ==  strings.ToLower(method) {
        var methodfu reflect.Value = reflect.ValueOf(controller).MethodByName(methodName)
        return &methodfu,&methodstr
     }
  }
  return nil,nil
}


func (rtp *RouteProcess ) getMethodInfo (controller controller.ControllerMethod,method string) (*reflect.Value,*reflect.Method) {
   controllerType := reflect.TypeOf(controller)
   var methodfu reflect.Value = reflect.ValueOf(controller).MethodByName(method)
   methodstr,res := controllerType.MethodByName(method)
   if  res  && methodstr.Type.NumIn() == 3  {
     return &methodfu,&methodstr
   }
   return nil,nil
}


/**
 * <name:int>
 *
**/
func (rtp *RouteProcess ) AddAuto(pathPattern string,controller controller.ControllerMethod,formType reflect.Type) {
   rtp.addRoute(pathPattern,controller,constvalue.DEFAULT_FUNC_NAME,formType,nil,true)
   var pathernStr = pathPattern
   if pathPattern[len(pathPattern)-1] == '/' {
      pathernStr = pathPattern[:len(pathPattern)-1]
   }
   controllerType := reflect.TypeOf(controller)
   for i := 0; i < controllerType.NumMethod(); i++ {
      methodstr := controllerType.Method(i)
      methodName := methodstr.Name
      if methodName != constvalue.BEFORE_FUNC &&  methodName != constvalue.AFTER_FUNC {
         rtp.addRoute(pathernStr + "/"+ strings.ToLower(methodName),controller,methodName,formType,nil,true)
      }
   }
}

/**
 * <name:int>
 *
**/
func (rtp *RouteProcess ) AddAutoWithViewPath(pathPattern string,controller controller.ControllerMethod,formType reflect.Type,viewPath string) {
   rtp.addRoute(pathPattern,controller,constvalue.DEFAULT_FUNC_NAME,formType,&viewPath,true)
   var pathernStr = pathPattern
   if pathPattern[len(pathPattern)-1] == '/' {
      pathernStr = pathPattern[:len(pathPattern)-1]
   }
   controllerType := reflect.TypeOf(controller)
   for i := 0; i < controllerType.NumMethod(); i++ {
      methodstr := controllerType.Method(i)
      methodName := methodstr.Name
      if methodName != constvalue.BEFORE_FUNC &&  methodName != constvalue.AFTER_FUNC {
         rtp.addRoute(pathernStr + "/"+ strings.ToLower(methodName),controller,methodName,formType,&viewPath,true)
      }
   }
}

/**
 * <name:int>
 *
**/
func (rtp *RouteProcess ) Add(pathPattern string,controller controller.ControllerMethod,method string,formType reflect.Type) {
   rtp.addRoute(pathPattern,controller,method,formType,nil,false)
}

func (rtp *RouteProcess ) AddWithViewPath(pathPattern string,controller controller.ControllerMethod,method string,formType reflect.Type,viewPath string) {
   rtp.addRoute(pathPattern,controller,method,formType,&viewPath,false)
}


func (rtp *RouteProcess ) addRoute(pathPattern string,controller controller.ControllerMethod,method string,formType reflect.Type,viewPath *string,auto bool) {
   pathPattern = strings.Trim(pathPattern," ")
   if len(pathPattern) == 0 {
      return
   }
   componentArray := strings.Split(pathPattern,"/")
   var parentNode *TreeNode = rtp.treeNode
   var urlPathIndex int = len(componentArray)
   var firstParameter bool = false
   for k, com := range componentArray {
     if  (k < urlPathIndex -1) && (len(com) == 0) {
        continue
     }else{
        hasReg := false
        var paraArray []ParaInfo
        treeNode := &TreeNode{}
        treeNode.left_children = make(map[string]*TreeNode)
        treeNode.right_children = make(map[string]*TreeNode)
        var paramname string = ""
        var paramtype string = ""
        var par *ParaInfo = nil
        var parseName ,parseType = false,false
        var nodename  string = ""
        var reg string
        var regstr *regexp.Regexp
        for _,ch := range com {
           if ch == '<' {
              if !firstParameter {
                urlPathIndex = k
              }
              firstParameter = true
              par = &ParaInfo{}
              parseName = true
              paramname = ""
              paramtype = ""
           }else if ch == '>' {
              par.ParaType, reg = rtp.paramTypeCheck(paramtype)
              par.ParaName = paramname
              nodename = nodename + reg
              paraArray = append(paraArray,*par)
              parseType = false
              parseName = false
           }else if ch == ':' {
             if parseName {
               parseName = false
               parseType = true
             }
           }else if ch == ' ' {
              continue
           }else {
              if parseName {
                paramname += string(ch)
              }else if parseType {
                paramtype += string(ch)
              }else {
                nodename += string(ch)
              }
           }
        }
        if len(paraArray) > 0 {
          treeNode.paraInfo = paraArray
          hasReg = true
          nodename = `^` + nodename + `$`
          var error_Info error
          regstr,error_Info =  regexp.Compile(nodename)
          if error_Info != nil {
             log.Error("Error : " + pathPattern)
          }
        }
        treeNode.nodeValue = nodename
        if k == len(componentArray) -1 {
           treeNode.isLeafNode = true
           treeNode.controller = controller
           treeNode.formType = formType
           if viewPath == nil {
              treeNode.UrlPath = strings.Join(componentArray[:urlPathIndex],"/")
           }else {
              treeNode.UrlPath = *viewPath
           }
           if strings.Trim(treeNode.UrlPath," ") == ""{
              treeNode.UrlPath = DEFAULT_URL
           }
           var methodInfo *reflect.Method
           var meth *reflect.Value
           if !auto {
             meth,methodInfo = rtp.checkMethod(controller,method)
           }else {
             meth,methodInfo = rtp.getMethodInfo(controller,method)
           }
           if meth != nil {
               treeNode.funcmap  = *meth
               treeNode.methodInfo  = *methodInfo
           }
        }
        if(!hasReg){
           treeNode.parent = parentNode
           treeNode.isReg = hasReg
           if parentNode.left_children[nodename] == nil {
              parentNode.left_children[nodename] = treeNode
              parentNode = treeNode
           }else {
              parentNode = parentNode.left_children[nodename]
           }
        }else {
           treeNode.parent = parentNode
           treeNode.isReg = hasReg
           treeNode.regstr = regstr
           if parentNode.right_children[nodename] == nil {
              parentNode.right_children[nodename] = treeNode
              parentNode = treeNode
           }else {
              parentNode = parentNode.right_children[nodename]
           }
        }
      }
   }
}


func (this *RouteInfo) getRequestInfo() string {
   return this.requestUrl
}

func (this *RouteInfo) getMatchResult() bool {
   return this.result
}

func (this *RouteInfo) Init(requestUrl string,result bool,controller controller.ControllerMethod,request * webhttp.HttpRequest) {
   this.requestUrl = requestUrl
   this.result = result
   this.controller = controller
   this.Request = request
}


func GetRouteProcess() *RouteProcess {
   return routeProcess
}

func redireToErrorPage(app *appcontext.AppContext, urlpath string, ErrorData interface{}) {
   app.CopyAppContext(routeProcess.ctx)
   switch urlpath {
      case constvalue.ERROR_403:
        app.UrlPath = app.GetConfigValue(constvalue.ERROR_403_PATH_KEY,constvalue.DEFAULT_ERROR_403_PATH).(string)
      case constvalue.ERROR_404:
        app.UrlPath = app.GetConfigValue(constvalue.ERROR_404_PATH_KEY,constvalue.DEFAULT_ERROR_404_PATH).(string)
      case constvalue.ERROR_405:
        app.UrlPath = app.GetConfigValue(constvalue.ERROR_405_PATH_KEY,constvalue.DEFAULT_ERROR_405_PATH).(string)
      case constvalue.ERROR_500:
        app.UrlPath = app.GetConfigValue(constvalue.ERROR_500_PATH_KEY,constvalue.DEFAULT_ERROR_500_PATH).(string)
   }
   var renderInfo *render.RenderInfo = render.CreateRenderInfo(app)
   renderInfo.InitRenderInfo(app)
   renderInfo.OutPutData = ErrorData
   app.Renderinfo = renderInfo
   res := interceptor.RenderOutPutinterceptor(app)
   if !res {
     log.Error("Can not Render Errer Page")
   }
}

func RedirectTo403(app *appcontext.AppContext) {
    redireToErrorPage(app,constvalue.ERROR_403,make(map[string]interface{}))
}

func RedirectTo404(app *appcontext.AppContext) {
    redireToErrorPage(app,constvalue.ERROR_404,make(map[string]interface{}))
}

func RedirectTo405(app *appcontext.AppContext) {
    redireToErrorPage(app,constvalue.ERROR_405,make(map[string]interface{}))
}

func RedirectTo500(app *appcontext.AppContext) {
    redireToErrorPage(app,constvalue.ERROR_405,make(map[string]interface{}))
}


func RedirectTo500AndErrorData(app *appcontext.AppContext) {
    redireToErrorPage(app,constvalue.ERROR_500,make(map[string]interface{}))
}


func RedirectTo403AndErrorData(app *appcontext.AppContext,ErrorData interface{}) {
    redireToErrorPage(app,constvalue.ERROR_403,ErrorData)
}


func RedirectTo404AndErrorData(app *appcontext.AppContext,ErrorData interface{}) {
    redireToErrorPage(app,constvalue.ERROR_404,ErrorData)
}


func RedirectTo405AndErrorData(app *appcontext.AppContext,ErrorData interface{}) {
    redireToErrorPage(app,constvalue.ERROR_405,ErrorData)
}
