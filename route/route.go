package route

import (
  "strings"
  "regexp"
  "../controller"
  "../appcontext"
  "../http"
  "../log"
)

func init() {
  log.InitLog()
}

type ParaInfo struct {
    paraName string
    paraType string
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
    controller *controller.Controller
    funcmap    *map[string] func()
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
   result         bool
   controller     *controller.Controller
   request        *http.HttpRequest
   methodName     string
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

func (rtp *RouteProcess ) match(urlcom []string,index int , treeNodemap map[string]*TreeNode) (bool, *controller.Controller) {
   var res bool = false
   return res,nil
}


func (rtp *RouteProcess ) urlRoute (request string, rtinfo *RouteInfo) bool{
   request = strings.Trim(request," ")
   if len(request) == 0 {
      return false
   }
   componentArray := strings.Split(request,"/")
   var res bool
   res,_  = rtp.match(componentArray,0,rtp.treeNode.left_children)
   if !res {
      res,_  = rtp.match(componentArray,0,rtp.treeNode.right_children)
   }
   return res
} 

func (rtp *RouteProcess ) ProcessRequest(request * http.HttpRequest) *RouteInfo {
   urlArray := strings.Split(request.Urlpath,"?")
   length := len(urlArray)
   if length > 0 {
      rinfo := &RouteInfo{paramInfo : make(map[string]string)}
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
      rtp.urlRoute(urlpath,rinfo)
      return rinfo
   }
   return nil
}


func (node *TreeNode ) DebugInfo() {
    log.Info(node.nodeValue)
    for _,nodep := range node.left_children {
       if nodep != nil {
         nodep.DebugInfo()
       }
    }
    for _,nodep := range node.right_children {
       if nodep != nil {
         nodep.DebugInfo()
       }
    }
}

func (rtp *RouteProcess ) DebugInfo() {
  rtp.treeNode.DebugInfo()
}


func (rtp *RouteProcess ) paramTypeCheck(typestr string) (string,string) {
    types := strings.ToLower(typestr)
    switch types {
      case "int" :
        return "int",`([-+][0-9]+)`
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


/**
 * <name:int>
 *
**/
func (rtp *RouteProcess ) Add(pathPattern string) {
   //controller * controller.Controller, method string
   pathPattern = strings.Trim(pathPattern," ")
   if len(pathPattern) == 0 {
      return
   }
   componentArray := strings.Split(pathPattern,"/")
   var parentNode *TreeNode = rtp.treeNode
   var treeNode *TreeNode = nil
   for _, com := range componentArray {
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
        var parseName = true
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
              par.paraName = paramname + reg
              nodename = nodename + par.paraName
              paraArray = append(paraArray,*par)
           }else if ch == ':' {
              parseName = false
           }else if ch == ' ' {
              continue
           }else {
              if parseName {
                paramname += string(ch)
              }else {
                paramtype += string(ch)
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
        }else{
          nodename = paramname
        }
        treeNode.nodeValue = nodename
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
   if treeNode != nil {
      treeNode.isLeafNode = true
   }
}

/**
 * <name:int>
 *
**/
func (rtp *RouteProcess ) AddAuto(controller * controller.Controller) {
   
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

func (this *RouteInfo) Init(requestUrl string, paramInfo map[string]string,result bool,controller controller.Controller,request * http.HttpRequest) {
   this.requestUrl = requestUrl
   this.paramInfo = paramInfo
   this.result = result
   this.controller = &controller
   this.request = request
}
