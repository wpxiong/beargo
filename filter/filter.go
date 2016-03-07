package filter

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/appcontext"
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

func ProcessAfterFilter(context *appcontext.AppContext) bool {
   for _,filterF := range filterManager.AfterFilter {
      res := filterF(context)
      if !res {
         return false
      }
   }
   return true
}

func AddDefaultFilter(){
  AddBeforeFilterList(ParameterParseFilter,ParameterBinderFilter)
  AddAfterFilterList(RenderFilter) 
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