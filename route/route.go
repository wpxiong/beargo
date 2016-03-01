package route

import (
  "strings"
  "regexp"
  "reflect"
  "../controller"
  "../appcontext"
  "../webhttp"
  "../log"
)

func init() {
  log.InitLog()
}

type ParaInfo struct {
    paraName string
    paraType string
    paraValue string
}


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
    paraInfo   []ParaInfo
    regstr     *regexp.Regexp
}


type RouteProcess struct {
   ctx       *appcontext.AppContext
   treeNode  *TreeNode
}


type RouteInfo struct {
   requestUrl     string
   paramInfo      map[string]string
   UrlParamInfo   []ParaInfo
   result         bool
   controller     controller.ControllerMethod
   Request        *webhttp.HttpRequest
   Writer         *webhttp.HttpResponse
   funcmap        reflect.Value
   methodName     string
   ResultChan     chan int
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
   return rtp.funcmap
}

func (rtp *RouteInfo ) GetResult() bool {
   return rtp.result
}


func (rtp *RouteInfo ) CallMethod() {
    var funcmap reflect.Value = rtp.getFuncmap()
    v := make([]reflect.Value, 1)
    v[0] = reflect.ValueOf(rtp)
    funcmap.Call(v)
}


func (rtp *RouteProcess ) match(urlcom []string,index int , treeNodemap map[string]*TreeNode, paraList *[][]ParaInfo) (bool, controller.ControllerMethod,reflect.Value) {
   *paraList = (append(*paraList,[]ParaInfo{}))
   if index >= len(urlcom) {
      return false,nil,reflect.Value{}
   }
   var res bool = false
   var funcm reflect.Value
   for _,treeNode := range treeNodemap {
      if !treeNode.isReg {
          if urlcom[index] == treeNode.nodeValue {
              if !treeNode.isLeafNode {
                var control  controller.ControllerMethod
                res,control,funcm = rtp.match(urlcom,index + 1,treeNode.left_children,paraList)
                if !res {
                   *paraList = (*paraList)[:len(*paraList)-1]
                   res,control,funcm = rtp.match(urlcom,index + 1,treeNode.right_children,paraList)
                   return res,control,funcm
                }else {
                   return res,control,funcm
                }
              }else {
                return true,treeNode.controller,treeNode.funcmap
              }
          }
      }else {
          res = treeNode.regstr.MatchString(urlcom[index])
          if res {
              reglist := treeNode.regstr.FindStringSubmatch(urlcom[index])
              pa := []ParaInfo{}
              for i,reval := range reglist {
                 if i >0 {
                    pa = append(pa,ParaInfo{paraName : treeNode.paraInfo[i-1].paraName, paraType : treeNode.paraInfo[i-1].paraType, paraValue: reval})
                 }
              }
              (*paraList)[len(*paraList)-1] = pa
              if !treeNode.isLeafNode {
                var control  controller.ControllerMethod
                res,control,funcm = rtp.match(urlcom,index + 1,treeNode.left_children,paraList)
                if !res {
                   *paraList = (*paraList)[:len(*paraList)-1]
                   res,control,funcm = rtp.match(urlcom,index + 1,treeNode.right_children,paraList)
                   return res,control,funcm
                }else {
                   return res,control,funcm
                }
              }else {
                return true,treeNode.controller,treeNode.funcmap
              }
          }
      }
     
   }
   return res,nil,reflect.Value{}
}


func (rtp *RouteProcess ) urlRoute (request string, rtinfo *RouteInfo) bool{
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
   var funcmap reflect.Value
   res,controller,funcmap  = rtp.match(componentArray,0,rtp.treeNode.left_children,&paraList)
   if !res {
      res,controller,funcmap  = rtp.match(componentArray,0,rtp.treeNode.right_children,&paraList)
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
   return res
} 

func (rtp *RouteProcess ) ProcessRequest(request * webhttp.HttpRequest) *RouteInfo {
   urlArray := strings.Split(request.Urlpath,"?")
   rinfo := &RouteInfo{paramInfo : make(map[string]string), UrlParamInfo: []ParaInfo{}}
   rinfo.result = false
   length := len(urlArray)
   if length > 0 {
      urlpath := urlArray[0]
      rinfo.requestUrl = urlpath
      if length >1 {
         paramArray := strings.Split( urlArray[1],"&")
         for _,param := range paramArray {
            ind := strings.Index(param,"=")
            if ind >=0 && ind < len(param) {
               rinfo.paramInfo[param[0:ind]] = param[ind+1:]
            }
         }
      }
      res := rtp.urlRoute(urlpath,rinfo)
      if res {
          rinfo.result = true
          rinfo.ResultChan = make(chan int)
      }
   }
   return rinfo
}


func (node *TreeNode ) debugInfo() {
    log.InfoNoReturn("{nodeValue:" + node.nodeValue + ",left_children:[")
    var k int = 0
    for _,nodep := range node.left_children {
       if nodep != nil {
         nodep.debugInfo()
         if k != len(node.left_children) -1 {
           log.InfoNoReturn(",")
         }
       }
       k++
    }
    k=0
    log.InfoNoReturn("],right_children:[")
    for _,nodep := range node.right_children {
       if nodep != nil {
         nodep.debugInfo()
         if k != len(node.right_children) -1 {
           log.InfoNoReturn(",")
         }
       }
       k++
    }
    log.InfoNoReturn("]")
    log.Info("}")
}

func (routeInfo *RouteInfo ) DebugInfo() {
   log.InfoNoReturn("{result:")
   log.InfoNoReturn(routeInfo.result)
   log.InfoNoReturn(",UrlParamInfo:[")
   for k,pam := range routeInfo.UrlParamInfo {
      log.InfoNoReturn(pam)
      if k != len(routeInfo.UrlParamInfo) -1 {
        log.InfoNoReturn(",")
      }
   }
   log.InfoNoReturn("],methodName:" + routeInfo.methodName)
   log.Info("}")
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


func (rtp *RouteProcess ) checkMethod (controller controller.ControllerMethod,method string) *reflect.Value {
  controllerType := reflect.TypeOf(controller)
  for i := 0; i < controllerType.NumMethod(); i++ {
     methodstr := controllerType.Method(i)
     methodName := methodstr.Name
     firstLetter := string(methodName[0])
     if methodstr.Type.NumIn() == 2  && firstLetter == strings.ToUpper(firstLetter) && strings.ToLower(methodName) ==  strings.ToLower(method) {
        var methodfu reflect.Value = reflect.ValueOf(controller).MethodByName(methodName)
        return &methodfu
     }
  }
  return nil
}



/**
 * <name:int>
 *
**/
func (rtp *RouteProcess ) Add(pathPattern string,controller controller.ControllerMethod,method string) {
   pathPattern = strings.Trim(pathPattern," ")
   if len(pathPattern) == 0 {
      return
   }
   componentArray := strings.Split(pathPattern,"/")
   var parentNode *TreeNode = rtp.treeNode
   for k, com := range componentArray {
     if len(com) == 0 {
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
        if len(com) == 0 {
          continue
        }
        for _,ch := range com {
           if ch == '<' {
              par = &ParaInfo{}
              parseName = true
              paramname = ""
              paramtype = ""
           }else if ch == '>' {
              par.paraType, reg = rtp.paramTypeCheck(paramtype)
              par.paraName = paramname
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
             log.Info("Error : " + pathPattern)
          }
        }
        treeNode.nodeValue = nodename
        if k == len(componentArray) -1 {
           treeNode.isLeafNode = true
           treeNode.controller = controller
           var meth *reflect.Value= rtp.checkMethod(controller,method)
           if meth != nil {
               treeNode.funcmap  = *meth
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

/**
 * <name:int>
 *
**/
func (rtp *RouteProcess ) AddAuto(controller controller.ControllerMethod) {
   
}



func (this *RouteInfo) getRequestInfo() string {
   return this.requestUrl
}

func (this *RouteInfo) getParamInfo() map[string]string {
   return this.paramInfo
}

func (this *RouteInfo) getMatchResult() bool {
   return this.result
}

func (this *RouteInfo) Init(requestUrl string, paramInfo map[string]string,result bool,controller controller.ControllerMethod,request * webhttp.HttpRequest) {
   this.requestUrl = requestUrl
   this.paramInfo = paramInfo
   this.result = result
   this.controller = controller
   this.Request = request
}
