package moudle

import (
  "github.com/wpxiong/beargo/log"
  "reflect"
  "strings"
  "time"
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


type DbSqlBuilder struct {
  sqlQuery  string
  structObj  interface{}
  tableInfo  *DBTableInfo
  fetchtype  FetchType
  fetchjointable []string
  structNameList []string
  moudle   *Moudle
}

type  DbQuery interface {
   Query(tableObj interface{},fetchType FetchType,structName []string) *DbSqlBuilder
   SimpleQuery(tableObj interface{}) *DbSqlBuilder
   Insert(structVal interface{} ) *DbSqlBuilder
   Update(structVal interface{} ) *DbSqlBuilder
   UpdateWithField(structVal interface{},fieldName[]string ) *DbSqlBuilder
   Delete(structVal interface{} ) *DbSqlBuilder
   
   InsertWithSql(sql string,parameter []interface{})
   UpdateWithSql(sql string,parameter []interface{})
   ExecuteWithSql(sql string,parameter []interface{})
   DeleteWithSql(sql string) 
}


type DbQueryInfo interface {
   Limit(limit int) *DbSqlBuilder
   FetchOne() interface{}
   FetchAll() [] interface{}
   WhereAnd(fieldName []string) *DbSqlBuilder
   WhereOr(fieldName []string) *DbSqlBuilder
   WhereWithSql(sql string, param []interface{}) *DbSqlBuilder
   And() *DbSqlBuilder
   Or() *DbSqlBuilder
   SaveExecute()
   DeleteExecute()
   InsertExecute()
}


func (this *DbSqlBuilder) FetchOne() (interface{},bool) {
   var res []interface{}
   var rows *sql.Rows
   var result interface{}
   vallist := make([]reflect.Value,0,0)
   var err error
   rows,err = this.moudle.DbProiver.Query(this.sqlQuery)
   defer rows.Close()
   if err != nil {
      log.Error(err)
      panic("SQL :[" + this.sqlQuery + "] query error")
   }else {
      this.processSelect(rows,&vallist)
      res = this.processValueToInterface(&vallist)
   }
   if len(res) > 0 {
      return res[0],true
   }else {
      return result,false
   }
}

func (this *Moudle) listField(list *[]interface{},tableInfo *DBTableInfo,obj *reflect.Value,typeVal reflect.Type) {
   for _,fieldName := range tableInfo.FieldList {
      val := tableInfo.FiledNameMap[fieldName]
      if val.RelationStructName == "" {
         fieldVal := (*obj).Elem().FieldByName(val.FieldName).Addr().Interface()
         switch val.FieldType.Kind() {
            case reflect.Complex64,reflect.Complex128:
               this.DbProiver.AppendScanComplexField(list)
            case reflect.Bool:
               var boolVield string
               (*list) = append(*list,&boolVield)
            default:
               (*list) = append(*list,fieldVal)
         }
      }
   }
}


func (this *Moudle) checkInList(obj reflect.Value,list [] reflect.Value,tableinfo * DBTableInfo) (bool, *reflect.Value) {
   for _,val := range list{
      var theSame bool = true
      for _,keyIndex := range tableinfo.KeyFieldIndex {
        fieldVal1 := reflect.ValueOf(obj.Elem().Interface()).Field(keyIndex).Interface()
        fieldVal2 := reflect.ValueOf(val.Elem().Interface()).Field(keyIndex).Interface()
        if fieldVal1 != fieldVal2 {
           theSame = false
           break
        }
     }
     if theSame {
        return true,&val
     }
   }
   return false,nil
}


func (this *DbSqlBuilder) processValueToInterface(valueList *[]reflect.Value) []interface{} {
   res := make([]interface{},0,0)
   for _,val := range *valueList{
      res = append(res,val.Elem().Interface())
   }
   return res
}

func (this *DbSqlBuilder) processSelect(rows *sql.Rows,vallist *[]reflect.Value) {
     for rows.Next() {
          objtype := reflect.TypeOf(this.structObj)
          structObj := reflect.New(objtype)
          var joinStructObj []reflect.Value = make([]reflect.Value,0,0)
          var columnObj []interface{} = make([]interface{},0,0)
          this.moudle.listField(&columnObj,this.tableInfo,&structObj,objtype)
          for index,table_Name := range this.fetchjointable {
             info := this.moudle.DbTableInfoByTableName[table_Name]
             joinObjtype := reflect.TypeOf(info.DbStuct)
             joinObj := reflect.New(joinObjtype)
             joinStructObj = append(joinStructObj,joinObj)
             this.moudle.listField(&columnObj,info,&joinObj,joinObjtype)
             columnInfo := this.tableInfo.FiledNameMap[strings.ToLower(this.structNameList[index])]
             if columnInfo.IsArray {
               nestedObj := structObj.Elem().FieldByName(this.structNameList[index])
               nestedObj.SetLen(0)
             }
          }
          if err := rows.Scan(columnObj...) ;err != nil {
            panic(err)
          }
          if ok, objVal := this.moudle.checkInList(structObj,*vallist,this.tableInfo); ok {
              structObj = *objVal
              for k,joinObj :=  range joinStructObj {
                columnInfo := this.tableInfo.FiledNameMap[strings.ToLower(this.structNameList[k])]
                nestedObj := structObj.Elem().FieldByName(this.structNameList[k])
                if columnInfo.IsArray {
                   nestedObj.Set(reflect.Append(nestedObj,joinObj.Elem()))
                }else {
                   nestedObj.Set(joinObj.Elem())
                }
             }
          }else {
             for k,joinObj :=  range joinStructObj {
                columnInfo := this.tableInfo.FiledNameMap[strings.ToLower(this.structNameList[k])]
                nestedObj := structObj.Elem().FieldByName(this.structNameList[k])
                if columnInfo.IsArray {
                   nestedObj.Set(reflect.Append(nestedObj,joinObj.Elem()))
                }else {
                   nestedObj.Set(joinObj.Elem())
                }
             }
             *vallist = append(*vallist,structObj)
          }
      }
}

func (this *DbSqlBuilder) FetchAll() []interface{} {
   vallist := make([]reflect.Value,0,0)
   var rows *sql.Rows
   var err error
   rows,err = this.moudle.DbProiver.Query(this.sqlQuery)
   defer rows.Close()
   if err != nil {
      log.Error(err)
      panic("SQL :[" + this.sqlQuery + "] query error")
   }else {
      this.processSelect(rows,&vallist)
      return this.processValueToInterface(&vallist)
   }
}

func (this *DbSqlBuilder) Limit(limit int) *DbSqlBuilder {
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

func (this *DbSqlBuilder) createJoinSql(structName []string) string {
   var sql string = ""
   var index int = 2
   for _,name := range structName {
      if struct_name,ok := this.moudle.searchField(this.tableInfo,name);ok {
         tableName := this.moudle.DbTableInfoByStructName[struct_name].TableName
         if joinstr,ok := this.moudle.createJoinSqlString(this.tableInfo.TableName,tableName,index) ; ok {
             tableName := this.moudle.DbTableInfoByStructName[struct_name].TableName
             this.fetchjointable = append(this.fetchjointable,tableName)
             this.structNameList = append(this.structNameList,name)
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
   for _,fieldName := range tableinfo.FieldList {
      val := tableinfo.FiledNameMap[fieldName]
      if val.RelationStructName == "" {
         columnList = append(columnList,"T" + strconv.Itoa(index) + "." + val.ColumnName )
      }
   }
   return columnList
}

func (this *Moudle) listColumn(tableinfo *DBTableInfo) []string {
   columnList := make([]string,0,0)
   for _,fieldName := range tableinfo.FieldList {
      val := tableinfo.FiledNameMap[fieldName]
      if val.RelationStructName == "" && val.AutoIncrement == false {
         columnList = append(columnList,val.ColumnName)
      }
   }
   return columnList
}

func (this *Moudle) listValue(tableinfo *DBTableInfo,structVal interface{}) []string {
   valueList := make([]string,0,0)
   for _,fieldName := range tableinfo.FieldList {
      val := tableinfo.FiledNameMap[fieldName]
      if val.RelationStructName == "" && val.AutoIncrement == false {
         value := reflect.ValueOf(structVal).FieldByName(val.FieldName)
         var valstr string
         switch val.FieldType.Kind(){
            case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
              valstr = strconv.FormatInt(value.Int(),10)
            case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
              valstr = strconv.FormatUint(value.Uint(), 10)
            case reflect.Uintptr:
              continue
            case reflect.Float32,reflect.Float64:
              valstr = strconv.FormatFloat(value.Float(),'f', -1, 64)
            case reflect.Complex64:
              valstr = this.DbProiver.GetInsertDBComplex64Sql(value.Complex())
            case reflect.Complex128:
              valstr = this.DbProiver.GetInsertDBComplex128Sql(value.Complex())
            case reflect.Array:
              continue
            case reflect.Bool:
               if value.Bool() {
                 valstr = "'1'"
               }else {
                 valstr = "'0'"
               }
            case reflect.Ptr:
              continue
            case reflect.Struct:
              if val.FieldType.Name() == "Time" {
                ti := value.Interface().(time.Time)
                valstr = this.DbProiver.GetInsertDBTimeSql(ti)
              }else {
                continue
              }
            case reflect.Slice:
              elemtype := val.FieldType.Elem().Name()
              if elemtype == "byte" {
                
              }else {
                continue
              }
            case reflect.Map:
              continue
            case reflect.String:
              valstr = "'" + value.String() + "'"
         }
         valueList = append(valueList,valstr )
      }
   }
   return valueList
}

func (this *Moudle) Query(tableObj interface{},fetchType FetchType,structName []string) *DbSqlBuilder {
   var table_info *DBTableInfo = nil
   info := DbSqlBuilder{sqlQuery:"",structObj:tableObj,tableInfo:table_info,fetchjointable:make([]string,0,0),fetchtype:fetchType,moudle:this,structNameList:make([]string,0,0)}
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
        columnlist := this.listTableColumn(info.tableInfo,1)
        info.sqlQuery = "select "  + strings.Join(columnlist,",") +   " from " + info.tableInfo.TableName  + " T1 "
      }
   }else {
     log.Error("not found the table relation with struct :" + struct_Name)
   }
   return &info
}


func (this *Moudle) SimpleQuery(tableObj interface{}) *DbSqlBuilder {
   var table_info *DBTableInfo = nil
   info := DbSqlBuilder{sqlQuery:"",structObj:tableObj,tableInfo:table_info,fetchjointable:make([]string,0,0),fetchtype:LAZY,moudle:this,structNameList:make([]string,0,0)}
   structName := reflect.TypeOf(tableObj).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      info.tableInfo = val
      info.sqlQuery = "select * from " + info.tableInfo.TableName;
   }else {
     log.Error("not found the table relation with struct :" + structName)
   }
   return &info
}


func (this *Moudle) createUpdateSql(tableinfo *DBTableInfo,structVal interface{},fieldNameList []string) string {
   valueList := make([]string,0,0)
   conditionList := make([]string,0,0)
   for _,fieldName := range tableinfo.FieldList {
      val := tableinfo.FiledNameMap[fieldName]
      if val.RelationStructName == "" && val.AutoIncrement == false {
         value := reflect.ValueOf(structVal).FieldByName(val.FieldName)
         var valstr string
         switch val.FieldType.Kind(){
            case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
              valstr = strconv.FormatInt(value.Int(),10)
            case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
              valstr = strconv.FormatUint(value.Uint(), 10)
            case reflect.Uintptr:
              continue
            case reflect.Float32,reflect.Float64:
              valstr = strconv.FormatFloat(value.Float(),'f', -1, 64)
            case reflect.Complex64:
              valstr = this.DbProiver.GetInsertDBComplex64Sql(value.Complex())
            case reflect.Complex128:
              valstr = this.DbProiver.GetInsertDBComplex128Sql(value.Complex())
            case reflect.Array:
              continue
            case reflect.Bool:
               if value.Bool() {
                 valstr = "'1'"
               }else {
                 valstr = "'0'"
               }
            case reflect.Ptr:
              continue
            case reflect.Struct:
              if val.FieldType.Name() == "Time" {
                ti := value.Interface().(time.Time)
                valstr = this.DbProiver.GetInsertDBTimeSql(ti)
              }else {
                continue
              }
            case reflect.Slice:
              elemtype := val.FieldType.Elem().Name()
              if elemtype == "byte" {
                
              }else {
                continue
              }
            case reflect.Map:
              continue
            case reflect.String:
              valstr = "'" + value.String() + "'"
         }
         if val.IsId {
            conditionList = append(conditionList,val.ColumnName + "=" + valstr)
         }else if len(fieldNameList) == 0 || InSliceTableList(fieldNameList,val.FieldName) {
            valueList = append(valueList,val.ColumnName + "=" + valstr)
         }
      }
   }
   return strings.Join(valueList,",") + " where " + strings.Join(conditionList," and ") 
}

func (this *DbSqlBuilder) InsertExecute() {
  if len(this.sqlQuery) > 0 {
    if _,err := this.moudle.DbProiver.Insert(this.sqlQuery); err != nil {
       panic(err)
    }
  }
}


func (this *Moudle) Insert(structVal interface{} ) *DbSqlBuilder {
   var table_info *DBTableInfo = nil
   info := DbSqlBuilder{sqlQuery:"",structObj:structVal,tableInfo:table_info,moudle:this}
   structName := reflect.TypeOf(structVal).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      info.tableInfo = val
      var valuelist []string = this.listValue(val,structVal)
      var columnlist []string = this.listColumn(val) 
      sqlstr := "insert into " + val.TableName + "(" + strings.Join(columnlist,",") + ") values (" +  strings.Join(valuelist,",")  +")"
      info.sqlQuery = sqlstr
   }else {
      panic("No table relation with struct: " + structName)
   }
   return &info
}

func (this *DbSqlBuilder) SaveExecute() {
  if len(this.sqlQuery) > 0 {
     if _,err := this.moudle.DbProiver.Update(this.sqlQuery); err != nil {
         panic(err)
     }
  }
}
      
func (this *Moudle) Update(structVal interface{} ) *DbSqlBuilder{
   var table_info *DBTableInfo = nil
   info := DbSqlBuilder{sqlQuery:"",structObj:structVal,tableInfo:table_info,moudle:this}
   structName := reflect.TypeOf(structVal).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      info.tableInfo = val
      updatestr := this.createUpdateSql(val,structVal,[]string{})
      if len(updatestr) > 0 {
         sqlstr := "update " + val.TableName + " set " + updatestr
         info.sqlQuery = sqlstr
      }
   }else {
      panic("No table relation with struct: " + structName)
   }
   return &info
}


func (this *Moudle) UpdateWithField(structVal interface{},fieldName[]string ) *DbSqlBuilder{
   var table_info *DBTableInfo = nil
   info := DbSqlBuilder{sqlQuery:"",structObj:structVal,tableInfo:table_info,moudle:this}
   structName := reflect.TypeOf(structVal).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      info.tableInfo = val
      updatestr := this.createUpdateSql(val,structVal,fieldName)
      if len(updatestr) > 0 {
        sqlstr := "update " + val.TableName + " set " +  updatestr
        info.sqlQuery = sqlstr
      }
   }else {
      panic("No table relation with struct: " + structName)
   }
   return &info
}

func (this *Moudle) InsertWithSql(sql string,parameter []interface{})  {

}

func (this *Moudle) UpdateWithSql(sql string,parameter []interface{}) {

}

func (this *Moudle) ExecuteWithSql(sql string,parameter []interface{}) {

}
