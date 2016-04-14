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
  haswherekeyword  bool
  paramlist  []interface{}
  ts      *Trans
}
   
type  DbQuery interface {
   Query(tableObj interface{},fetchType FetchType,structName []string) *DbSqlBuilder
   SimpleQuery(tableObj interface{}) *DbSqlBuilder
   Insert(structVal interface{} ) *DbSqlBuilder
   Update(structVal interface{} )
   UpdateWithField(structVal interface{},fieldName[]string)
   UpdateWithWhere(structVal interface{},fieldName[]string)  *DbSqlBuilder
   
   InsertWithSql(sql string,parameter []interface{})
   UpdateWithSql(sql string,parameter []interface{})
   ExecuteWithSql(sql string,parameter []interface{})
   DeleteWithSql(sql string,parameter []interface{})
   DeleteWhere(structVal interface{}) *DbSqlBuilder
   Delete(structVal interface{} ) 
}


type DbQueryInfo interface {
   Limit(limit int) *DbSqlBuilder
   FetchOne() interface{}
   FetchAll() [] interface{}
   FetchLazyField(fieldName []string)
   WhereAnd(fieldName []string) *DbSqlBuilder
   WhereOr(fieldName []string) *DbSqlBuilder
   WhereWithSql(sql string, param []interface{}) *DbSqlBuilder
   And() *DbSqlBuilder
   Or() *DbSqlBuilder
   Start() *DbSqlBuilder
   End() *DbSqlBuilder
   SaveExecute()
   DeleteExecute()
   InsertExecute()
   GetCurrentSerialValue(structVal interface{}) int64 
}



func (this *Moudle) query(sqlStr string,ts *Trans) (*sql.Rows, error){
   return this.DbProiver.Query(sqlStr,ts)
}

func (this *Moudle) insert(sqlStr string,ts *Trans) (sql.Result ,error){
   return this.DbProiver.Insert(sqlStr,ts)
}

func (this *Moudle) update(sqlStr string,ts *Trans) (sql.Result ,error){
    return this.DbProiver.Update(sqlStr,ts)
}

func (this *Moudle)  prepareExecuteSQL(sqlStr string ,parameter []interface{},ts *Trans) {
    this.DbProiver.PrepareExecuteSQL(sqlStr,parameter,ts)
}

func (this *Moudle) executeSQL(sqlStr string,ts *Trans) (sql.Result ,error){
    return this.DbProiver.ExecuteSQL(sqlStr,ts)
}


func (this *DbSqlBuilder) deleteExecute(ts *Trans) {
  if len(this.sqlQuery) > 0 {
    if len(this.paramlist) == 0 { 
      if _,err := this.moudle.DbProiver.ExecuteSQL(this.sqlQuery,ts); err != nil {
         panic(err)
      }
    }else {
      this.moudle.DbProiver.PrepareExecuteSQL(this.sqlQuery,this.paramlist,ts)
    }
  }
}

func (this *DbSqlBuilder) fetchOne(ts *Trans) (interface{},bool) {
   var res []interface{}
   var rows *sql.Rows
   var result interface{}
   var err error
   rows,err = this.moudle.DbProiver.Query(this.sqlQuery,ts)
   vallist := make([]reflect.Value,0)
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


func (this *Moudle) checkInList(obj reflect.Value,list [] reflect.Value,size int ,tableinfo * DBTableInfo) (bool, *reflect.Value) {
   var k int = 0
   for _,val := range list{
      if k >= size && size >= 0{
         break
      }
      var theSame bool = true
      for _,keyIndex := range tableinfo.KeyFieldIndex {
        fieldVal1 := reflect.ValueOf(obj.Elem().Interface()).Field(keyIndex).Interface()
        fieldVal2 := reflect.ValueOf(val.Elem().Interface()).Field(keyIndex).Interface()
        if fieldVal1 != fieldVal2 {
           theSame = false
           break
        }
     }
     k += 1
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
     var k int = 0
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
          if ok, objVal := this.moudle.checkInList(structObj,*vallist,k,this.tableInfo); ok {
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
             k += 1
          }
      }
}

func (this *Moudle) fetchLazyField(structObj interface{},fieldName []string,ts *Trans) {
   structName := reflect.TypeOf(structObj).Elem().Name()
   var tableInfo *DBTableInfo
   var ok bool
   if tableInfo,ok = this.DbTableInfoByStructName[structName]; !ok {
      panic("not find the table relation with the interface :" + structName)
   }
   for _,name := range fieldName {
      if struct_name,ok := this.searchField(tableInfo,name);ok {
         table_info := this.DbTableInfoByStructName[struct_name]
         tableObj := table_info.DbStuct
         var newSqlBuilder *DbSqlBuilder =  &DbSqlBuilder{sqlQuery:"",structObj:tableObj,tableInfo:table_info,fetchjointable:make([]string,0,0),fetchtype:LAZY,moudle:this,structNameList:make([]string,0,0)}
         var refrencedcolumnName string = ""
         var columnName string = ""
         for _,val := range this.RelationInfoList {
            if val.DbTableName == tableInfo.TableName && val.StructName == table_info.StructName {
              refrencedcolumnName = val.ReferencedColumnName
              columnName = val.ColumnName
            } 
         }
         columnlist := this.listTableColumn(table_info,0)
         var sqlQuery string = "select "  + strings.Join(columnlist,",") +   " from " + table_info.TableName
         var columnInfo ColumnInfo
         if columnName != "" {
            columnInfo = tableInfo.FiledNameMap[strings.ToLower(columnName)]
            fieldVal := reflect.ValueOf(structObj).Elem().FieldByName(columnInfo.FieldName)
            sqlQuery += " where " + refrencedcolumnName + " = "  + this.getFieldValueString(&columnInfo,&fieldVal)
         }else {
             panic("not find the table relation with the interface :" + structName)
         }
         log.Debug(sqlQuery)
         newSqlBuilder.sqlQuery = sqlQuery
         vallist := make([]reflect.Value,0)
         var rows *sql.Rows
         var err error
         rows,err = this.DbProiver.Query(newSqlBuilder.sqlQuery,ts)
         defer rows.Close()
         if err != nil {
            log.Error(err)
             panic("SQL :[" + newSqlBuilder.sqlQuery + "] query error")
         }else {
            newSqlBuilder.processSelect(rows,&vallist)
            fieldType ,_ := reflect.TypeOf(structObj).Elem().FieldByName(name)
            switch  fieldType.Type.Kind() {
               case reflect.Slice:
                 nestedObj :=  reflect.ValueOf(structObj).Elem().FieldByName(name)
                 for _,elem := range vallist {
                    nestedObj.Set(reflect.Append(nestedObj,elem.Elem()))
                 }
               default:
                 nestedObj :=  reflect.ValueOf(structObj).Elem().FieldByName(name)
                 for _,elem := range vallist {
                    nestedObj.Set(elem.Elem())
                    break
                }
            }
         }
      }else {
         panic("not found fieldName db table name relation with field " + name)
      }
   }
}

func (this *DbSqlBuilder) fetchAll(ts *Trans) []interface{} {
   vallist := make([]reflect.Value,0)
   var rows *sql.Rows
   var err error
   rows,err = this.moudle.DbProiver.Query(this.sqlQuery,ts)
   defer rows.Close()
   if err != nil {
      log.Error(err)
      panic("SQL :[" + this.sqlQuery + "] query error")
   }else {
      this.processSelect(rows,&vallist)
      return this.processValueToInterface(&vallist)
   }
}

func (this *Moudle) createFieldList(tableInfo *DBTableInfo,structVal interface{},fieldNameList []string) []string {
   keyList := make([]string,0,0)
   columnList := make([]string,0,0)
   keyvalueList := make([]string,0,0)
   for _,filedName := range fieldNameList {
      if columnInfo,ok := tableInfo.FiledNameMap[strings.ToLower(filedName)];ok {
         keyList = append(keyList,strings.ToLower(columnInfo.FieldName))
         columnList = append(columnList,columnInfo.ColumnName)
      }
   }
   var valueList []string = this.listValue(tableInfo,structVal,keyList,false)
   for i:=0; i<len(columnList) ;i++ {
      keyvalueList = append(keyvalueList, columnList[i] + " = " + valueList[i])
   }
   return keyvalueList
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
         if index > 0 {
           columnList = append(columnList,"T" + strconv.Itoa(index) + "." + val.ColumnName )
         }else {
           columnList = append(columnList,val.ColumnName )
         }
      }
   }
   return columnList
}

func (this *Moudle) listColumn(tableinfo *DBTableInfo,fieldList []string) []string {
   columnList := make([]string,0,0)
   for _,fieldName := range tableinfo.FieldList {
      val := tableinfo.FiledNameMap[fieldName]
      if val.RelationStructName == "" && val.AutoIncrement == false {
         if len(fieldList) == 0 {
            columnList = append(columnList,val.ColumnName)
         }else if InSliceStringList(fieldList,val.FieldName) {
            columnList = append(columnList,val.ColumnName)
         }
      }
   }
   return columnList
}

func (this *Moudle) listValue(tableinfo *DBTableInfo,structVal interface{},fieldNameList []string,insert bool) []string {
   valueList := make([]string,0,0)
   for _,fieldName := range tableinfo.FieldList {
      val := tableinfo.FiledNameMap[fieldName]
      if len(fieldNameList) !=0 && !InSliceStringList(fieldNameList,fieldName) {
          continue
      }
      log.Debug(fieldName)
      if val.RelationStructName == "" && (insert && val.AutoIncrement == false  || !insert) {
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
              if elemtype == "uint8" {
                valstr = this.DbProiver.GetInsertByteDataSql(value.Interface().([]byte))
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

func (this *Moudle)  getFieldValueString(val *ColumnInfo,value *reflect.Value) string {
  var valstr string = ""
  switch val.FieldType.Kind(){
    case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
       valstr = strconv.FormatInt(value.Int(),10)
    case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
       valstr = strconv.FormatUint(value.Uint(), 10)
    case reflect.Uintptr:
    case reflect.Float32,reflect.Float64:
       valstr = strconv.FormatFloat(value.Float(),'f', -1, 64)
    case reflect.Complex64:
       valstr = this.DbProiver.GetInsertDBComplex64Sql(value.Complex())
    case reflect.Complex128:
       valstr = this.DbProiver.GetInsertDBComplex128Sql(value.Complex())
    case reflect.Array:
    case reflect.Bool:
       if value.Bool() {
          valstr = "'1'"
       }else {
          valstr = "'0'"
       }
    case reflect.Ptr:
    case reflect.Struct:
       if val.FieldType.Name() == "Time" {
         ti := value.Interface().(time.Time)
         valstr = this.DbProiver.GetInsertDBTimeSql(ti)
       }
    case reflect.Slice:
       elemtype := val.FieldType.Elem().Name()
       if elemtype == "uint8" {
          valstr = this.DbProiver.GetInsertByteDataSql(value.Interface().([]byte))
       }
    case reflect.Map:
    case reflect.String:
       valstr = "'" + value.String() + "'"
    }
    return valstr
}

func (this *Moudle) createUpdateSql(tableinfo *DBTableInfo,structVal interface{},fieldNameList []string) string {
   valueList := make([]string,0,0)
   conditionList := make([]string,0,0)
   for _,fieldName := range tableinfo.FieldList {
      val := tableinfo.FiledNameMap[fieldName]
      if val.RelationStructName == "" && val.AutoIncrement == false {
         value := reflect.ValueOf(structVal).FieldByName(val.FieldName)
         var valstr string
         valstr = this.getFieldValueString(&val,&value)
         if valstr == "" {
            continue
         }
         if val.IsId {
            conditionList = append(conditionList,val.ColumnName + "=" + valstr)
         }else if len(fieldNameList) == 0 || InSliceStringList(fieldNameList,val.FieldName) {
            valueList = append(valueList,val.ColumnName + "=" + valstr)
         }
      }
   }
   return strings.Join(valueList,",") + " where " + strings.Join(conditionList," and ") 
}


func (this *DbSqlBuilder) insertExecute(ts *Trans) {
  if len(this.sqlQuery) > 0 {
    if _,err := this.moudle.DbProiver.Insert(this.sqlQuery,ts); err != nil {
       panic(err)
    }
  }
}


func (this *DbSqlBuilder) saveExecute(ts *Trans) {
  if len(this.sqlQuery) > 0 {
     if len(this.paramlist) == 0 { 
       if _,err := this.moudle.DbProiver.Update(this.sqlQuery,ts); err != nil {
         panic(err)
       }
     }else {
       this.moudle.DbProiver.PrepareExecuteSQL(this.sqlQuery,this.paramlist,ts)
     }
  }
}

func (this *Moudle) updateDB(structVal interface{} ,ts *Trans){
   structName := reflect.TypeOf(structVal).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      updatestr := this.createUpdateSql(val,structVal,[]string{})
      if len(updatestr) > 0 {
         sqlstr := "update " + val.TableName + " set " + updatestr
         if _,err := this.DbProiver.Update(sqlstr,ts); err != nil {
            panic(err)
         }
      }
   }else {
      panic("No table relation with struct: " + structName)
   }
}

func (this *Moudle) updateWithField(structVal interface{},fieldName[]string,ts* Trans) {
   structName := reflect.TypeOf(structVal).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      updatestr := this.createUpdateSql(val,structVal,fieldName)
      if len(updatestr) > 0 {
        sqlstr := "update " + val.TableName + " set " +  updatestr
        if _,err := this.DbProiver.Update(sqlstr,ts); err != nil {
           panic(err)
        }
      }
   }else {
      panic("No table relation with struct: " + structName)
   }
}



func (this *Moudle) createKeyWhere(tableInfo *DBTableInfo,structVal interface{}) string {
   keyList := make([]string,0,0)
   columnList := make([]string,0,0)
   wherList := make([]string,0,0)
   for _,columnInfo := range tableInfo.FiledNameMap {
      if columnInfo.IsId {
         keyList = append(keyList,strings.ToLower(columnInfo.FieldName))
         columnList = append(columnList,columnInfo.ColumnName)
      }
   }
   var valueList []string = this.listValue(tableInfo,structVal,keyList,false)
   for i:=0; i<len(columnList) ;i++ {
      wherList = append(wherList, columnList[i] + " = " + valueList[i])
   }
   return " where " + strings.Join(wherList, " and ")
}

func (this *Moudle)  delete(structVal interface{} ,ts *Trans) {
   structName := reflect.TypeOf(structVal).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      sqlstr := "delete from " + val.TableName + this.createKeyWhere(val,structVal)
      if _,err := this.DbProiver.ExecuteSQL(sqlstr,ts); err != nil {
         panic(err)
      }
   }else {
      panic("No table relation with struct: " + structName)
   }
}

func (this *Moudle) getCurrentSerialValue(structVal interface{},ts *Trans) int64 {
   structName := reflect.TypeOf(structVal).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      for _,columnInfo := range val.FiledNameMap {
        if columnInfo.IsId && columnInfo.AutoIncrement {
           return this.DbProiver.GetCurrentSerialValue(val.TableName,columnInfo.ColumnName,ts)
        }
      }
      panic("The Table [" + val.TableName + "] has no Serial column")
   }else {
      panic("No table relation with struct: " + structName)
   }
   
}

func (this *Moudle)  deleteWithSql(sql string,parameter []interface{},ts *Trans) {
   this.DbProiver.PrepareExecuteSQL(sql,parameter,ts)
}

func (this *Moudle) insertWithSql(sql string,parameter []interface{},ts *Trans)  {
  this.DbProiver.PrepareExecuteSQL(sql,parameter,ts)
}

func (this *Moudle) updateWithSql(sql string,parameter []interface{},ts *Trans) {
   this.DbProiver.PrepareExecuteSQL(sql,parameter,ts)
}


func (this *Moudle) executeWithSql(sql string,parameter []interface{},ts *Trans) {
   this.DbProiver.PrepareExecuteSQL(sql,parameter,ts)
}


//Public Method

func (this *Moudle)  Delete(structVal interface{} ) {
   this.delete(structVal,nil)
}

func (this *DbSqlBuilder) FetchOne() (interface{},bool) {
    return this.fetchOne(nil)
}

func (this *DbSqlBuilder) Limit(limit int) *DbSqlBuilder {
   this.sqlQuery += this.moudle.DbProiver.LimitSql(limit)
   return this
}

func (this *Moudle) FetchLazyField(structObj interface{} ,fieldName []string) {
   this.fetchLazyField(structObj,fieldName,nil)
}

func (this *DbSqlBuilder) FetchAll() []interface{} {
   return this.fetchAll(nil)
}

func (this *DbSqlBuilder) WhereAnd(fieldName []string) *DbSqlBuilder {
   if this.tableInfo != nil && this.sqlQuery != "" {
     if !this.haswherekeyword {
        this.sqlQuery += " where "
        this.haswherekeyword = true
     }
     this.sqlQuery += " " + strings.Join(this.moudle.createFieldList(this.tableInfo,this.structObj,fieldName)," and ")
   }else {
      panic("create where sql error!")
   }
   return this
}



func (this *DbSqlBuilder) WhereOr(fieldName []string) *DbSqlBuilder {
   if this.tableInfo != nil && this.sqlQuery != "" {
     if !this.haswherekeyword {
        this.sqlQuery += " where "
        this.haswherekeyword = true
     }
     this.sqlQuery += " " + strings.Join(this.moudle.createFieldList(this.tableInfo,this.structObj,fieldName)," or ")
   }else {
      panic("create where sql error!")
   }
   return this
}


func (this *DbSqlBuilder) WhereWithSql(sql string, param []interface{}) *DbSqlBuilder {
   if this.tableInfo != nil && this.sqlQuery != "" {
     if !this.haswherekeyword {
        this.sqlQuery += " where "
        this.haswherekeyword = true
     }
     this.sqlQuery += " " + sql
     this.paramlist = append(this.paramlist,param...)
   }else {
      panic("create where sql error!")
   }
   return this
}

func (this *DbSqlBuilder) And() *DbSqlBuilder {
   if this.tableInfo != nil && this.sqlQuery != "" {
     this.sqlQuery += " and " 
   }else {
      panic("create where sql error!")
   }
   return this
}

func (this *DbSqlBuilder) Or(sql string, param []interface{}) *DbSqlBuilder {
   if this.tableInfo != nil && this.sqlQuery != "" {
     this.sqlQuery += " or " 
   }else {
      panic("create where sql error!")
   }
   return this
}

func (this *DbSqlBuilder) End() *DbSqlBuilder {
   if this.tableInfo != nil && this.sqlQuery != "" {
     this.sqlQuery += " ( " 
   }else {
      panic("create where sql error!")
   }
   return this
}

func (this *DbSqlBuilder) Start() *DbSqlBuilder {
   if this.tableInfo != nil && this.sqlQuery != "" {
     this.sqlQuery += " ) " 
   }else {
      panic("create where sql error!")
   }
   return this
}

func (this *Moudle) Insert(structVal interface{} ) *DbSqlBuilder {
   var table_info *DBTableInfo = nil
   info := DbSqlBuilder{sqlQuery:"",structObj:structVal,tableInfo:table_info,moudle:this}
   structName := reflect.TypeOf(structVal).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      info.tableInfo = val
      var valuelist []string = this.listValue(val,structVal,[]string{},true)
      var columnlist []string = this.listColumn(val,[]string{}) 
      sqlstr := "insert into " + val.TableName + "(" + strings.Join(columnlist,",") + ") values (" +  strings.Join(valuelist,",")  +")"
      info.sqlQuery = sqlstr
   }else {
      panic("No table relation with struct: " + structName)
   }
   return &info
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
        columnlist := this.listTableColumn(info.tableInfo,0)
        info.sqlQuery = "select "  + strings.Join(columnlist,",") +   " from " + info.tableInfo.TableName
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
      columnlist := this.listTableColumn(info.tableInfo,0)
      info.sqlQuery = "select "  + strings.Join(columnlist,",") +   " from " + info.tableInfo.TableName
   }else {
     log.Error("not found the table relation with struct :" + structName)
   }
   return &info
}


func (this *Moudle) Update(structVal interface{} ){
   this.updateDB(structVal,nil)
}

func (this *Moudle) UpdateWithWhere(structVal interface{},fieldName[]string)  *DbSqlBuilder {
   var table_info *DBTableInfo = nil
   info := DbSqlBuilder{sqlQuery:"",structObj:structVal,tableInfo:table_info,moudle:this}
   structName := reflect.TypeOf(structVal).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      info.tableInfo = val
      var valuelist []string
      var columnlist []string 
      var columnValueList []string
      if len(fieldName) == 0 {
         valuelist = this.listValue(val,structVal,[]string{},false)
         columnlist = this.listColumn(val,[]string{})
      }else {
         fieldList := make([]string,len(fieldName))
         for key,field := range fieldName {
             fieldList[key] = strings.ToLower(field)
         }
         valuelist = this.listValue(val,structVal,fieldList,false)
         columnlist = this.listColumn(val,fieldName)
         columnValueList = make([]string,len(valuelist))
         for key,columnName := range columnlist {
            columnValueList[key] = columnName + "=" + valuelist[key]
         }
      }
      if len(columnValueList) > 0 {
        sqlstr := "update " + val.TableName + " set " +  strings.Join(columnValueList,",")
        info.sqlQuery = sqlstr
        log.Debug(info.sqlQuery)
      }
   }else {
      panic("No table relation with struct: " + structName)
   }
   return &info
}

func (this *Moudle) UpdateWithField(structVal interface{},fieldName[]string ) {
   this.updateWithField(structVal,fieldName,nil)
}


func (this *Moudle) DeleteWhere(structVal interface{}) *DbSqlBuilder {
   var table_info *DBTableInfo = nil
   info := DbSqlBuilder{sqlQuery:"",structObj:structVal,tableInfo:table_info,moudle:this}
   structName := reflect.TypeOf(structVal).Name()
   if val,ok := this.DbTableInfoByStructName[structName]; ok {
      info.tableInfo = val
      sqlstr := "delete from " + val.TableName
      info.sqlQuery = sqlstr
   }else {
      panic("No table relation with struct: " + structName)
   }
   return &info
}


func (this *Moudle)  DeleteWithSql(sql string,parameter []interface{}) {
   this.deleteWithSql(sql,parameter,nil)
}


func (this *Moudle) InsertWithSql(sql string,parameter []interface{})  {
  this.insertWithSql(sql,parameter,nil)
}


func (this *Moudle) UpdateWithSql(sql string,parameter []interface{}) {
   this.prepareExecuteSQL(sql,parameter,nil)
}

func (this *Moudle) ExecuteWithSql(sql string,parameter []interface{}) {
   this.executeWithSql(sql,parameter,nil)
}

func (this *Moudle) GetCurrentSerialValue(structVal interface{}) int64 {
   return this.getCurrentSerialValue(structVal,nil)
}
