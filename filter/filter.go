package filter

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
  "github.com/wpxiong/beargo/constvalue"
  "reflect"
)

func init() {
  log.InitLog()
}

type FilterFunc func(*appcontext.AppContext)(bool)


type  Filter  struct {
   AfterFilter  [] FilterFunc
   BeforeFilter [] FilterFunc
}

var filterManager *Filter 

var redirectfilter FilterFunc

func InitFilter() {
  if filterManager == nil {
      filterManager = &Filter{BeforeFilter:make([] FilterFunc,0),AfterFilter:make([] FilterFunc,0)}
  }
}

func ProcessBeforeFilter(context *appcontext.AppContext) bool {
   for _,filterF := range filterManager.BeforeFilter {
      res := filterF(context)
      if !res {
         return false
      }
   }
   return true
}

func ProcessRedirectBeforeFilter(context *appcontext.AppContext) bool {
   for _,filterF := range filterManager.BeforeFilter {
      switch reflect.TypeOf(filterF).Name() {
        case constvalue.ParameterParseFilter :
           continue
        default:
           
      }
      res := filterF(context)
      if !res {
         return false
      }
   }
   return true
}

func ProcessAfterFilter(context *appcontext.AppContext) bool {
   for _,filterF := range filterManager.AfterFilter {
      res := filterF(context)
      if !res {
         return false
      }
   }
   return true
}

func AddInitFilter(context *appcontext.AppContext, funcMap map[string]FilterFunc){
   beforeFilterList :=  context.GetConfigValue(constvalue.BEFORE_FILTER_KEY,constvalue.DEFULT_BEFORE_FILTER).([]string)
   afterFilterList :=  context.GetConfigValue(constvalue.AFTER_FILTER_KEY,constvalue.DEFULT_AFTER_FILTER).([]string)
   var beforeFuncList []FilterFunc = make([]FilterFunc ,0,0)
   for _,key := range beforeFilterList{
      val := funcMap[key]
      if val != nil {
        beforeFuncList = append(beforeFuncList,val)
      }
   }
   
   var afterFuncList []FilterFunc = make([]FilterFunc ,0,0)
   for _,key := range afterFilterList{
      val := funcMap[key]
      if val != nil {
        afterFuncList = append(afterFuncList,val)
      }
   }
   AddBeforeFilterList(beforeFuncList...)
   AddAfterFilterList(afterFuncList...)
}


func AddDefaultFilter(){
  AddBeforeFilterList(ParameterParseFilter,ParameterBinderFilter)
  AddAfterFilterList(RenderBindFilter,RenderOutPutFilter) 
}

func AddBeforeFilterList (filterFunc... FilterFunc){
   for _, filterF := range filterFunc {
       filterManager.AddBeforeFilter(filterF)
   }
}


func AddAfterFilterList (filterFunc... FilterFunc){
   for _, filterF := range filterFunc {
       filterManager.AddAfterFilter(filterF)
   }
}

func AddBeforeFilter (filterFunc FilterFunc){
   filterManager.AddBeforeFilter(filterFunc)
}

func AddAfterFilter (filterFunc FilterFunc){
   filterManager.AddAfterFilter(filterFunc)
}


func (filter *Filter) AddBeforeFilter (filterFunc FilterFunc){
   filter.BeforeFilter = append(filter.BeforeFilter,filterFunc)
}

func (filter *Filter) AddAfterFilter (filterFunc FilterFunc){
   filter.AfterFilter = append(filter.AfterFilter,filterFunc)
}


