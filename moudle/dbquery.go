package moudle

import (
  "github.com/wpxiong/beargo/log"
  "reflect"
)

func init() {
  log.InitLog()
}

type QueryInfo struct {
  sqlQuery  string
  structObj  interface{}
  tableInfo  *DBTableInfo
}

type  DbQuery interface {
   Query(tableObj interface{}) QueryInfo
}


func (this *Moudle) Query(tableObj interface{}) QueryInfo {
   var table_info *DBTableInfo = nil
   info := QueryInfo{sqlQuery:"",structObj:tableObj,tableInfo:table_info}
   structName := reflect.TypeOf(tableObj).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      info.tableInfo = val
   }
   return info
}