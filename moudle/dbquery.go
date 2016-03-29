package moudle

import (
  "github.com/wpxiong/beargo/log"
  "reflect"
  "strings"
  "strconv"
  "database/sql"
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
}

type DbQueryInfo interface {
   Limit(limit int) *QueryInfo
   GetOneResult() interface{}
   GetResultList() [] interface{}
}


func (this *QueryInfo) GetOneResult() interface{} {
   var res interface{}
   var rows *sql.Rows
   var err error
   rows,err = this.moudle.DbProiver.Query(this.sqlQuery)
   if err != nil {
      log.Error(err)
      panic("SQL :[" + this.sqlQuery + "] query error")
   }else {
      log.Debug(rows)
   }
   return res
}

func (this *QueryInfo) GetResultList() []interface{} {
   res := make([]interface{},0,0)
   var rows *sql.Rows
   var err error
   rows,err = this.moudle.DbProiver.Query(this.sqlQuery)
   if err != nil {
      log.Error(err)
      panic("SQL :[" + this.sqlQuery + "] query error")
   }else {
      log.Debug(rows)
   }
   return res
}

func (this *QueryInfo) Limit(limit int) *QueryInfo {
   this.sqlQuery += this.moudle.DbProiver.LimitSql(limit)
   return this
}

func (this *Moudle) searchField(tableInfo  *DBTableInfo, fieldName string) (string,bool) {
    for _,val := range tableInfo.FiledNameMap {
       if val.FieldName == fieldName {
          return val.RelationStructName,true
       }
    }
    return "",false
}

func (this *Moudle) createJoinSqlString(tableName1 string,tableName2 string,tableIndex int) (string,bool){
   var keymap map[string][]ForeignKeyInfo
   for key,val := range this.RelationMap {
     if key == tableName1 {
        keymap = val.(map[string][]ForeignKeyInfo)
        for key2,val2 := range keymap {
           for _,info := range val2 {
             if info.TableName == tableName2 {
                return "T1" + "." + key2 + " = " + "T" + strconv.Itoa(tableIndex) + "."  + info.KeyColumnName ,true
             }
           }
        }
     }else if key == tableName2 {
        keymap = val.(map[string][]ForeignKeyInfo)
        for key2,val2 := range keymap {
            for _,info := range val2 {
             if info.TableName == tableName1 {
                return "T" + strconv.Itoa(tableIndex) + "." + key2 + " = " + "T1" + "."  + info.KeyColumnName ,true
             }
           }
        }
     }
   }
   return "",false
}

func (this *QueryInfo) createJoinSql(structName []string) string {
   var sql string = ""
   var index int = 2
   for _,name := range structName {
      if struct_name,ok := this.moudle.searchField(this.tableInfo,name);ok {
         tableName := this.moudle.DbTableInfoByStructName[struct_name].TableName
         if joinstr,ok := this.moudle.createJoinSqlString(this.tableInfo.TableName,tableName,index) ; ok {
             tableName := this.moudle.DbTableInfoByStructName[struct_name].TableName
             sql = sql + " left join " + tableName  + " T"  + strconv.Itoa(index) + " on " + joinstr
         }else{ 
           panic("create join condition error type : " + struct_name)
         }
      }else {
          panic("create join condition error fieldName : " + name)
      }
      index += 1
   }
   return sql
}

func (this *Moudle) listTableColumn(tableinfo *DBTableInfo,index int) []string {
   columnList := make([]string,0,0)
   for _,val := range tableinfo.FiledNameMap {
      if val.RelationStructName == "" {
         columnList = append(columnList,"T" + strconv.Itoa(index) + "." + val.ColumnName )
      }
   }
   return columnList
}

func (this *Moudle) Query(tableObj interface{},fetchType FetchType,structName []string) QueryInfo {
   var table_info *DBTableInfo = nil
   info := QueryInfo{sqlQuery:"",structObj:tableObj,tableInfo:table_info,fetchjointable:make([]string,0,0),fetchtype:fetchType,moudle:this}
   defer func() {
      if err := recover(); err != nil {
         log.Error(err)
         info.tableInfo = nil
      }
   }()
   struct_Name := reflect.TypeOf(tableObj).Name()
   if val,ok := this.DbTableInfoByStructName[struct_Name]; ok {
      if fetchType == EAGER {
        info.tableInfo = val
        var tableIndex int = 1
        columnlist := this.listTableColumn(info.tableInfo,tableIndex)
        for _,name := range structName {
          if struct_name,ok := this.searchField(val,name);ok {
              table_info := this.DbTableInfoByStructName[struct_name]
              tableIndex += 1
              list := this.listTableColumn(table_info,tableIndex);
              columnlist = append(columnlist,list...)
          }else {
              panic("not found fieldName db table name relation with field " + name)
          }
        }
        info.sqlQuery = "select "  + strings.Join(columnlist,",") +   " from " + info.tableInfo.TableName  + " T1 " + info.createJoinSql(structName)
      }else {
        info.tableInfo = val
        info.sqlQuery = "select * from " + info.tableInfo.TableName
      }
   }else {
     log.Error("not found the table relation with struct :" + struct_Name)
   }
   return info
}


func (this *Moudle) SimpleQuery(tableObj interface{}) QueryInfo {
   var table_info *DBTableInfo = nil
   info := QueryInfo{sqlQuery:"",structObj:tableObj,tableInfo:table_info,fetchjointable:make([]string,0,0),fetchtype:LAZY,moudle:this}
   structName := reflect.TypeOf(tableObj).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      info.tableInfo = val
      info.sqlQuery = "select * from " + info.tableInfo.TableName;
   }else {
     log.Error("not found the table relation with struct :" + structName)
   }
   return info
}