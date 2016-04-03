package interceptor

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/constvalue"
  "reflect"
)

func init() {
  log.InitLog()
}

type InterceptorFunc func(*appcontext.AppContext)(bool)


type  interceptor  struct {
   Afterinterceptor  [] InterceptorFunc
   Beforeinterceptor [] InterceptorFunc
}

var interceptorManager *interceptor 

var redirectinterceptor InterceptorFunc

func Initinterceptor() {
  if interceptorManager == nil {
      interceptorManager = &interceptor{Beforeinterceptor:make([] InterceptorFunc,0),Afterinterceptor:make([] InterceptorFunc,0)}
  }
}

func ProcessBeforeinterceptor(context *appcontext.AppContext) bool {
   for _,interceptorF := range interceptorManager.Beforeinterceptor {
      res := interceptorF(context)
      if !res {
         return false
      }
   }
   return true
}

func ProcessRedirectBeforeinterceptor(context *appcontext.AppContext) bool {
   for _,interceptorF := range interceptorManager.Beforeinterceptor {
      switch reflect.TypeOf(interceptorF).Name() {
        case constvalue.ParameterParseinterceptor :
           continue
        default:
           
      }
      res := interceptorF(context)
      if !res {
         return false
      }
   }
   return true
}

func ProcessAfterinterceptor(context *appcontext.AppContext) bool {
   for _,interceptorF := range interceptorManager.Afterinterceptor {
      res := interceptorF(context)
      if !res {
         return false
      }
   }
   return true
}

func AddInitinterceptor(context *appcontext.AppContext, funcMap map[string]InterceptorFunc){
   beforeinterceptorList :=  context.GetConfigValue(constvalue.BEFORE_interceptor_KEY,constvalue.DEFULT_BEFORE_interceptor).([]string)
   afterinterceptorList :=  context.GetConfigValue(constvalue.AFTER_interceptor_KEY,constvalue.DEFULT_AFTER_interceptor).([]string)
   var beforeFuncList []InterceptorFunc = make([]InterceptorFunc,4)
   var index int = 0
   for _,key := range beforeinterceptorList{
      val := funcMap[key]
      if val != nil {
        beforeFuncList = append(beforeFuncList[:index],val)
        index+=1
      }
   }
   AddBeforeinterceptorList(beforeFuncList...)
   
   var afterFuncList []InterceptorFunc = make([]InterceptorFunc,2)
   index = 0
   for _,key := range afterinterceptorList{
      val := funcMap[key]
      if val != nil {
        afterFuncList = append(afterFuncList[:index],val)
        index+=1
      }
   }
   AddAfterinterceptorList(afterFuncList...)
}


func AddDefaultinterceptor(){
  AddBeforeinterceptorList(ParameterParseinterceptor,ParameterBinderinterceptor)
  AddAfterinterceptorList(RenderBindinterceptor,RenderOutPutinterceptor) 
}

func AddBeforeinterceptorList (InterceptorFunc... InterceptorFunc){
   for _, interceptorF := range InterceptorFunc {
       interceptorManager.AddBeforeinterceptor(interceptorF)
   }
}


func AddAfterinterceptorList (InterceptorFunc... InterceptorFunc){
   for _, interceptorF := range InterceptorFunc {
       interceptorManager.AddAfterinterceptor(interceptorF)
   }
}

func AddBeforeinterceptor (InterceptorFunc InterceptorFunc){
   interceptorManager.AddBeforeinterceptor(InterceptorFunc)
}

func AddAfterinterceptor (InterceptorFunc InterceptorFunc){
   interceptorManager.AddAfterinterceptor(InterceptorFunc)
}


func (interceptor *interceptor) AddBeforeinterceptor (InterceptorFunc InterceptorFunc){
   interceptor.Beforeinterceptor = append(interceptor.Beforeinterceptor,InterceptorFunc)
}

func (interceptor *interceptor) AddAfterinterceptor (InterceptorFunc InterceptorFunc){
   interceptor.Afterinterceptor = append(interceptor.Afterinterceptor,InterceptorFunc)
}


