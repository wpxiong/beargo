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
   AfterFilter  map[string] FilterFunc
   BeforeFilter map[string] FilterFunc
}

var filterManager *Filter 


func InitFilter() {
  if filterManager == nil {
      filterManager = &Filter{}
  }
}

func ProcessBeforeFilter(){

}

func ProcessAfterFilter(){

}


func AddBeforeFilter (filterFunc FilterFunc){
   filterManager.AddBeforeFilter(filterFunc)
}

func AddAfterFilter (filterFunc FilterFunc){
   filterManager.AddAfterFilter(filterFunc)
}


func (filter *Filter) AddBeforeFilter (filterFunc FilterFunc){


}

func (filter *Filter) AddAfterFilter (filterFunc FilterFunc){


}