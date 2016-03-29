package moudle

import (
  "github.com/wpxiong/beargo/log"
  "reflect"
)

func init() {
  log.InitLog()
}



type FetchType int

const (
    EAGER   FetchType = iota
    LAZY
)


type QueryInfo struct {
  sqlQuery  string
  structObj  interface{}
  tableInfo  *DBTableInfo
  fetchtype  FetchType
  fetchjointable []string
  moudle   *Moudle
}

type  DbQuery interface {
   Query(tableObj interface{},fetchType FetchType,structName []string) QueryInfo
   SimpleQuery(tableObj interface{}) QueryInfo
   Limit(limit int) QueryInfo
   GetOneResult() interface{}
   GetResultList() [] interface{}
}

func (this *QueryInfo) GetResultList() []interface{} {
   res := make([]interface{},0,0)
   return res
}

func (this *Moudle) searchField(tableInfo  *DBTableInfo, fieldName string) (string,bool) {
    for _,val := range tableInfo.FiledNameMap {
       if val.FieldName == fieldName {
          return val.RelationStructName,true
       }
    }
    return "",false
}

func (this *Moudle) createJoinSqlString(tableName1 string,tableName2 string) (string,bool){
   var keymap map[string][]ForeignKeyInfo
   for key,val := range this.RelationMap {
     if key == tableName1 {
        keymap = val.(map[string][]ForeignKeyInfo)
        for key2,val2 := range keymap {
           for _,info := range val2 {
             if info.TableName == tableName2 {
                return tableName1 + "." + key2 + " = " + tableName2 + "."  + info.KeyColumnName ,true
             }
           }
        }
     }else if key == tableName2 {
        keymap = val.(map[string][]ForeignKeyInfo)
        for key2,val2 := range keymap {
            for _,info := range val2 {
             if info.TableName == tableName1 {
                return tableName2 + "." + key2 + " = " + tableName1 + "."  + info.KeyColumnName ,true
             }
           }
        }
     }
   }
   return "",false
}

func (this *QueryInfo) createJoinSql(structName []string) string {
   var sql string = ""
   for _,name := range structName {
      if struct_name,ok := this.moudle.searchField(this.tableInfo,name);ok {
         tableName := this.moudle.DbTableInfoByStructName[struct_name].TableName
         if joinstr,ok := this.moudle.createJoinSqlString(this.tableInfo.TableName,tableName) ; ok {
             tableName := this.moudle.DbTableInfoByStructName[struct_name].TableName
             sql = sql + " left join " + tableName  + " on " + joinstr
         }
      }
   }
   return sql
}

func (this *Moudle) Query(tableObj interface{},fetchType FetchType,structName []string) QueryInfo {
   var table_info *DBTableInfo = nil
   info := QueryInfo{sqlQuery:"",structObj:tableObj,tableInfo:table_info,fetchjointable:make([]string,0,0),fetchtype:fetchType,moudle:this}
   struct_Name := reflect.TypeOf(tableObj).Name()
   if val,ok := this.DbTableInfoByStructName[struct_Name]; ok {
      info.tableInfo = val
      info.sqlQuery = "select * from " + info.tableInfo.TableName + info.createJoinSql(structName)
   }
   log.Debug(info.sqlQuery)
   return info
}


func (this *Moudle) SimpleQuery(tableObj interface{}) QueryInfo {
   var table_info *DBTableInfo = nil
   info := QueryInfo{sqlQuery:"",structObj:tableObj,tableInfo:table_info,fetchjointable:make([]string,0,0),fetchtype:LAZY,moudle:this}
   structName := reflect.TypeOf(tableObj).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      info.tableInfo = val
      info.sqlQuery = "select * from " + info.tableInfo.TableName;
   }
   return info
}