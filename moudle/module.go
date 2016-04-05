package moudle

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/constvalue"
  "github.com/wpxiong/beargo/util"
  "database/sql"
  "errors"
  "reflect"
  "strings"
)

func init() {
  log.InitLog()
}


type DbDialectType int

const (
    MYSQL   DbDialectType = iota
    SQLITE
    POSTGRESQL
)

type DbRelationType int

const (
    ONE_TO_MANY   DbRelationType = iota
    MANY_TO_ONE
    ONE_TO_ONE
)


type ColumnInfo struct {
  ColumnName    string
  NotNull       bool
  SqlType       string
  IsId          bool
  Length        int
  Scale         int
  Unique        bool
  DefaultValue  interface{}
  FieldName     string
  FieldType     reflect.Type
  RelationStructName  string
  IsArray       bool
  AutoIncrement bool
}


type ForeignKeyInfo struct {
  TableName        string
  KeyColumnName          string
}

type DBTableInfo struct {
  TableName        string
  DbSchema      string
  DbStuct       interface{}
  FiledNameMap  map[string] ColumnInfo
  FieldList      []string
  DbTableExist  bool
  StructName    string
  KeyFieldIndex []int
}

type RelationInfo struct {
  DbTableName  string
  RelationType DbRelationType
  ColumnName string
  ReferencedColumnName  string
  StructName    string
}

type Moudle struct {
  DbDialect        DbDialectType
  DbName           string
  DbConnectionUrl  string
  DbUserName       string
  DbPassword       string
  DbTableInfoByTableName       map[string]*DBTableInfo
  DbTableInfoByStructName      map[string]*DBTableInfo
  DbProiver        DbProviderInterface
  connectionStatus bool
  RelationInfoList []RelationInfo
  RelationMap      map[string]interface{}
}


func CreateModuleInstance(DbDialect DbDialectType,DbName string,DbConnectionUrl string, DbUserName string,DbPassword string) *Moudle {
   module :=  &Moudle{DbDialect:DbDialect,DbName:DbName,DbConnectionUrl:DbConnectionUrl,DbUserName:DbUserName,DbPassword:DbPassword,RelationInfoList:make([]RelationInfo,0,0)}
   module.initModuleInstance()
   return module
}

func (this *Moudle) initModuleInstance(){
   this.DbTableInfoByTableName = make(map[string]*DBTableInfo)
   this.DbTableInfoByStructName = make(map[string]*DBTableInfo)
   var connectionUrl string  
   switch this.DbDialect {
      case MYSQL :
        connectionUrl = this.DbUserName  + ":" + this.DbPassword +   "@" + this.DbConnectionUrl + "/" + this.DbName + "?parseTime=true"
        this.DbProiver = &MysqlDBProvider{}
      case POSTGRESQL :
   }
   err := this.DbProiver.ConnectionDb(connectionUrl)
   if err != nil {
      log.Error("DB Connection Error!")
   }else {
      this.connectionStatus = true
      this.DbProiver.SetMinConnection(constvalue.DEFAULT_MIN_DB_CONNECTION)
      this.DbProiver.SetMaxConnection(constvalue.DEFAULT_MAX_DB_CONNECTION)
   }
}


func isEmpty(strvalue string) bool {
  if strings.Trim(strvalue," ") == "" {
     return true
  }
  return false
}

func (this *Moudle) AddTableWithTableName(dbtable interface{},tableName string){
   this.addTable(dbtable,tableName,"")
}


func (this *Moudle) AddTable(dbtable interface{}){
   this.addTable(dbtable,"","")
}

func (this *Moudle) AddTableWithSchema(dbtable interface{},tableName string ,tableSchema string){
   this.addTable(dbtable,tableName,tableSchema)
}


func (this *Moudle) droptable(tablename string) {
   _,err := this.DbProiver.DropTable(tablename)
   if err != nil {
      log.Error(err)
   }
}

func (this *Moudle) createtable(tableName,sqlstr string,primaryKey []string) {
   _,err := this.DbProiver.CreateTable(tableName,sqlstr,primaryKey)
   if err != nil {
      log.Error(err)
      panic("create table error")
   }
}


func (this *Moudle) createPrimaryKey(tableName string,keyList []string) {
   _,err := this.DbProiver.CreatePrimaryKey(tableName,keyList)
   if err != nil {
      log.Error(err)
   }
}


func (this *Moudle) createForeignKey(tableName string , keyColumn string, refrenceTableName string, referenceColumnName string) {
   _,err := this.DbProiver.CreateForeignKey(tableName,keyColumn,refrenceTableName,referenceColumnName)
   if err != nil {
      log.Error(err)
      panic("create foreign key error")
   }
}

func InSlice (arr []ForeignKeyInfo, val ForeignKeyInfo) (bool){
    for _, v := range(arr) {
       if v.TableName == val.TableName &&  v.KeyColumnName == val.KeyColumnName  { return true; } 
    }    
    return false; 
}


func (this *Moudle)  createRelation() map[string]interface{} {
   var sqlMap map[string]interface{}  = make(map[string]interface{})
   for _,info := range this.RelationInfoList {
      if info.ColumnName == "" || info.ReferencedColumnName == ""{ 
         continue
      }
      if info.RelationType == ONE_TO_MANY {
         val := sqlMap[info.DbTableName] 
         var keymap map[string][]ForeignKeyInfo 
         if val != nil {
           keymap = val.(map[string][]ForeignKeyInfo)
         }else {
           keymap = make(map[string][]ForeignKeyInfo)
           sqlMap[info.DbTableName] = keymap
         }
         foreignInfo := ForeignKeyInfo{}
         tableInfo := this.DbTableInfoByStructName[info.StructName]
         if tableInfo != nil && info.DbTableName != tableInfo.TableName {
           foreignInfo.TableName = tableInfo.TableName
           foreignInfo.KeyColumnName = info.ReferencedColumnName
           if !InSlice(keymap[info.ReferencedColumnName],foreignInfo) {
              keymap[info.ColumnName] = append(keymap[info.ColumnName],foreignInfo)
           }
         }
      }else if info.RelationType == MANY_TO_ONE {
         tableInfo := this.DbTableInfoByStructName[info.StructName]
         if tableInfo != nil {
             val := sqlMap[tableInfo.TableName]
             var keymap map[string][]ForeignKeyInfo 
             if val != nil {
                keymap = val.(map[string][]ForeignKeyInfo)
             }else {
                keymap = make(map[string][]ForeignKeyInfo)
                sqlMap[tableInfo.TableName] = keymap
             }
             if info.DbTableName != tableInfo.TableName { 
               foreignInfo := ForeignKeyInfo{}
               foreignInfo.TableName = info.DbTableName
               foreignInfo.KeyColumnName = info.ColumnName
               if !InSlice(keymap[info.ReferencedColumnName],foreignInfo) {
                  keymap[info.ReferencedColumnName] = append(keymap[info.ReferencedColumnName],foreignInfo)
               }
             }
         }
      }else if info.RelationType == ONE_TO_ONE {
         tableInfo := this.DbTableInfoByStructName[info.StructName]
         if tableInfo != nil {
             val := sqlMap[tableInfo.TableName]
             var keymap map[string][]ForeignKeyInfo 
             if val != nil {
                keymap = val.(map[string][]ForeignKeyInfo)
             }else {
                keymap = make(map[string][]ForeignKeyInfo)
                sqlMap[tableInfo.TableName] = keymap
             }
             if info.DbTableName != tableInfo.TableName { 
               foreignInfo := ForeignKeyInfo{}
               foreignInfo.TableName = info.DbTableName
               foreignInfo.KeyColumnName = info.ColumnName
               if !InSlice(keymap[info.ReferencedColumnName],foreignInfo) {
                  keymap[info.ReferencedColumnName] = append(keymap[info.ReferencedColumnName],foreignInfo)
               }
             }
         }
      }
   }
   return sqlMap
}


func InSliceStringList (arr []string, val string) (bool){
    for _, v := range(arr) {
       if v == val  { return true; } 
    }    
    return false; 
}

func searchTableInRelation(relationMap map[string]interface{},tableName string, tableList *[]string){
   if val,ok := relationMap[tableName];ok {
      mapforegin :=  val.(map[string][]ForeignKeyInfo) 
      for _,val2 := range mapforegin {
         for _, info := range val2 {
            if !InSliceStringList(*tableList,info.TableName) {
               searchTableInRelation(relationMap,info.TableName,tableList)
            }
         }
      }
   }else {
      if !InSliceStringList(*tableList,tableName) {
        (*tableList) = append(*tableList,tableName)
      }
   }
}

func (this *Moudle) sortTable(relationMap map[string]interface{}) []*DBTableInfo {
   sortMap := make([]*DBTableInfo,0,0)
   tableList := make([]string,0,0)
   for _,Info := range this.DbTableInfoByTableName {
      searchTableInRelation(relationMap,Info.TableName,&tableList)
      if !InSliceStringList(tableList,Info.TableName) {
        tableList = append(tableList,Info.TableName)
      }
   }
   
   for _,table := range tableList {
     sortMap = append(sortMap,this.DbTableInfoByTableName[table])
   }
   return sortMap
}

func (this *Moudle)  TableExistsInDB(tableName string) bool{
  if val,err := this.DbProiver.TableExistsInDB(tableName) ; err != nil {
     panic(err)
  }else {
     return val
  }
}

func (this *Moudle)  InitialDB(create bool) {
   log.Debug("Initial DB start")
   var index int = 0
   defer func() {
      if err := recover(); err != nil {
          log.Error("create database error")
       }
   }() 
   sqlMap := this.createRelation()
   this.RelationMap = sqlMap
   sorttablemap := this.sortTable(sqlMap)
   if create {
      for _,Info := range sorttablemap {
         this.droptable(Info.TableName)
      }
   }
   
   for _,Info := range sorttablemap {
      primaryKey := make([]string,0,0)
      var create_sql string = "create table " + Info.TableName + " ( "
      for _,column := range Info.FiledNameMap {
        if column.RelationStructName != "" {
           continue
        }
        columnstr := column.ColumnName + " " +  this.createSqlTypeByLength(column.AutoIncrement,column.SqlType,column.Length,column.Scale)
        if column.NotNull {
           columnstr += " not null "
        }
        if column.Unique {
           columnstr += " unique "
        }
        columnstr += this.createDefaultValue(column.DefaultValue)
        columnstr += ",\n"
        create_sql += columnstr
        if column.IsId {
           primaryKey = append(primaryKey,column.ColumnName)
        }
        index ++
      }
      if create {
        create_sql = create_sql[0:len(create_sql)-2]
        create_sql += "\n)"
        this.createtable(Info.TableName,create_sql,primaryKey)
      }else if !this.TableExistsInDB(Info.TableName){  
        create_sql = create_sql[0:len(create_sql)-2]
        create_sql += "\n)"
        this.createtable(Info.TableName,create_sql,primaryKey)
      }
   }
   this.createForeignKeyByRelation(sqlMap)
}


func (this *Moudle)  createForeignKeyByRelation(sqlMap map[string]interface{}) {
   for table,val := range sqlMap {
      mapforegin :=  val.(map[string][]ForeignKeyInfo) 
      for columnname,val2 := range mapforegin {
         for _, info := range val2 {
            this.createForeignKey(info.TableName,info.KeyColumnName,table,columnname)
         }
      }
   }
}

func (this *Moudle)  getDBIntType() string {
   return this.DbProiver.GetDBIntType()
}

func (this *Moudle)  getDBInt8Type() string {
   return this.DbProiver.GetDBInt8Type()
}

func (this *Moudle)  getDBInt16Type() string {
    return this.DbProiver.GetDBInt16Type()
}

func (this *Moudle)  getDBInt32Type() string {
   return this.DbProiver.GetDBInt32Type()
}

func (this *Moudle)  getDBInt64Type() string {
  return this.DbProiver.GetDBInt64Type()
}

func (this *Moudle)  getDBUintType() string {
  return this.DbProiver.GetDBUintType()
}

func (this *Moudle)  getDBUint8Type() string {
  return this.DbProiver.GetDBUint8Type()
}

func (this *Moudle)  getDBUint16Type() string {
  return this.DbProiver.GetDBUint16Type()
}

func (this *Moudle)  getDBUint32Type() string {
  return this.DbProiver.GetDBUint32Type()
}

func (this *Moudle)  getDBUint64Type() string {
  return this.DbProiver.GetDBUint64Type()
}

func (this *Moudle)  getDBFloat32Type() string {
  return this.DbProiver.GetDBFloat32Type()
}


func (this *Moudle)  getDBFloat64Type() string {
  return this.DbProiver.GetDBFloat64Type()
}

func (this *Moudle)  getDBComplex64Type() string {
  return this.DbProiver.GetDBComplex64Type()
}

func (this *Moudle)  getDBComplex128Type() string {
  return this.DbProiver.GetDBComplex128Type()
}

func (this *Moudle)  getDBBoolType() string {
  return this.DbProiver.GetDBBoolType()
}

func (this *Moudle)  getDBStringType(length int) string {
  return this.DbProiver.GetDBStringType(length)
}

func (this *Moudle)  getDBTimeType() string {
  return this.DbProiver.GetDBTimeType()
}

func (this *Moudle) getDBByteArrayType(length int ) string {
  return this.DbProiver.GetDBByteArrayType(length)
}

func (this *Moudle) beginTransaction() *sql.Tx {
   tx ,err := this.DbProiver.Begin()
   if err == nil {
     return tx
   }else {
     log.Error(err)
     return nil
   }
}

func (this *Moudle) endTransaction(tx *sql.Tx) {
  if err := this.DbProiver.Commit(tx); err != nil {
    log.Error(err)
  }
}

func (this *Moudle)  createSqlTypeByLength(auto_increment bool,sqltype string, length int,scale int) string {
  return this.DbProiver.CreateSqlTypeByLength(auto_increment,sqltype,length,scale)
}

func (this *Moudle)  createDefaultValue(defaultValue interface{}) string {
  return this.DbProiver.CreateDefaultValue(defaultValue)
}

func (this *Moudle)  getRelationType(typestr string) (DbRelationType,error) {
    if typestr == constvalue.DB_RELATION_ONE_TO_MANY {
       return ONE_TO_MANY,nil
    } else if typestr == constvalue.DB_RELATION_MANY_TO_ONE {
       return MANY_TO_ONE,nil
    }else if typestr == constvalue.DB_RELATION_ONE_TO_ONE {
       return ONE_TO_ONE,nil
    }else {
       return ONE_TO_MANY,errors.New("relation type is must " + constvalue.DB_RELATION_ONE_TO_MANY + " or " + constvalue.DB_RELATION_MANY_TO_ONE)
    }
}


func (this *Moudle) addTable(dbtable interface{},tableName string,schemaname string){
  if !this.connectionStatus {
     return 
  }else {
     tableInfo := DBTableInfo{}
     structName := reflect.TypeOf(dbtable).Name()
     tablenamestr := strings.ToLower(structName)     
     fieldNum := reflect.TypeOf(dbtable).NumField()
     tableInfo.FiledNameMap = make(map[string]ColumnInfo)
     tableInfo.KeyFieldIndex = make([]int,0)
     tableInfo.FieldList = make([]string,0)
     if tableName == "" {
       tableInfo.TableName = tablenamestr
     }else {
       tableInfo.TableName = strings.ToLower(tableName) 
     }
     
     for i:=0;i<fieldNum;i++{
         field := reflect.TypeOf(dbtable).Field(i)
         id := field.Tag.Get(constvalue.DB_ID)
         columnInfo := ColumnInfo{}
         column_name := field.Tag.Get(constvalue.DB_COLUMN_NAME)
         not_null := field.Tag.Get(constvalue.DB_NOT_NULL)
         auto_increment := field.Tag.Get(constvalue.DB_AUTO_INCREMENT)
         length := field.Tag.Get(constvalue.DB_LENGTH)
         scale := field.Tag.Get(constvalue.DB_SCALE)
         unique_key := field.Tag.Get(constvalue.DB_UNIQUE_KEY)
         default_value := field.Tag.Get(constvalue.DB_DEFAULT_VALUE)
         referenced_column_name := field.Tag.Get(constvalue.DB_REFERENCED_COLUMN_NAME)
         referenced_type := field.Tag.Get(constvalue.DB_RELATION_TYPE)
         referenced_column_name = strings.Trim(referenced_column_name," ")
         referenced_type = strings.Trim(referenced_type," ")
         
         if len,err := util.GetIntValue(length); err == nil {
             columnInfo.Length = int(len)
         }
         
         switch field.Type.Kind() {
            case reflect.Int:
              columnInfo.SqlType = this.getDBIntType()
              if val,err := util.GetIntValue(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Int8:
              columnInfo.SqlType = this.getDBInt8Type()
              if val,err := util.GetInt8Value(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Int16:
              columnInfo.SqlType = this.getDBInt16Type()
              if val,err := util.GetInt16Value(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Int32:
              columnInfo.SqlType = this.getDBInt32Type()
              if val,err := util.GetInt32Value(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Int64:
              columnInfo.SqlType = this.getDBInt64Type()
              if val,err := util.GetInt64Value(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Uint:
              columnInfo.SqlType = this.getDBUintType()
              if val,err := util.GetUintValue(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Uint8:
              columnInfo.SqlType = this.getDBUint8Type()
              if val,err := util.GetUint8Value(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Uint16:
              columnInfo.SqlType = this.getDBUint16Type()
              if val,err := util.GetUint16Value(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Uint32:
              columnInfo.SqlType = this.getDBUint32Type()
              if val,err := util.GetUint32Value(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Uint64:
              columnInfo.SqlType = this.getDBUint64Type()
              if val,err := util.GetUint64Value(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Uintptr:
              continue
            case reflect.Float32:
              columnInfo.SqlType = this.getDBFloat32Type()
              if val,err := util.GetFloat32Value(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Float64:
              columnInfo.SqlType = this.getDBFloat64Type()
              if val,err := util.GetFloat64Value(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Complex64:
              columnInfo.SqlType = this.getDBComplex64Type()
              if val,err := util.GetComplex64Value(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Complex128:
              columnInfo.SqlType = this.getDBComplex128Type()
              if val,err := util.GetComplex128Value(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Array:
              continue
            case reflect.Bool:
              columnInfo.SqlType = this.getDBBoolType()
              if val,err := util.GetBoolValue(default_value); err == nil {
                 columnInfo.DefaultValue = val
              }
            case reflect.Ptr:
              continue
            case reflect.Struct:
              if field.Type.Name() == "Time" {
                columnInfo.SqlType = this.getDBTimeType()
                columnInfo.DefaultValue = default_value
              }else {
                columnInfo.RelationStructName = field.Type.Name()
                columnInfo.IsArray = false
              }
            case reflect.Slice:
              elemtype := field.Type.Elem().Name()
              if elemtype == "uint8" {
                columnInfo.SqlType = this.getDBByteArrayType(columnInfo.Length)
              }else {
                // foreign key
                columnInfo.RelationStructName = field.Type.Elem().Name()
                columnInfo.IsArray = true
              }
            case reflect.Map:
              continue
            case reflect.String:
              columnInfo.SqlType = this.getDBStringType(columnInfo.Length)
              columnInfo.DefaultValue = default_value
         }
         columnInfo.FieldType = field.Type
         columnInfo.FieldName = field.Name
         if !isEmpty(column_name) {
            columnInfo.ColumnName = column_name
         }else{
            columnInfo.ColumnName = strings.ToLower(field.Name)
         }
         
         if strings.ToLower(strings.Trim(id," ")) == "true" {
            columnInfo.IsId = true
            tableInfo.KeyFieldIndex = append(tableInfo.KeyFieldIndex,i)
         }
         
         if strings.ToLower(strings.Trim(auto_increment," ")) == "true" {
            columnInfo.AutoIncrement = true
         }
         
         if strings.ToLower(strings.Trim(not_null," ")) == "true" {
            columnInfo.NotNull = true
         }
         
         if len,err := util.GetIntValue(scale); err == nil {
             columnInfo.Scale = int(len)
         }
         
         if strings.ToLower(strings.Trim(unique_key," ")) == "true" {
            columnInfo.Unique = true
         }
         tableInfo.FiledNameMap[strings.ToLower(field.Name)] = columnInfo
         tableInfo.FieldList = append(tableInfo.FieldList,strings.ToLower(field.Name))
         if referenced_type  != "" {
            relation := RelationInfo{ReferencedColumnName:referenced_column_name,DbTableName: tableInfo.TableName,ColumnName:column_name,StructName:columnInfo.RelationStructName}
            re_type,err := this.getRelationType(referenced_type)
            if err == nil {
               relation.RelationType = re_type
               this.RelationInfoList = append(this.RelationInfoList,relation) 
            }else {
               log.Error(err)
            }
         }
     }
     
     if schemaname == "" {
       tableInfo.DbSchema = ""
     }else {
       tableInfo.DbSchema = schemaname
     }
     tableInfo.DbStuct = dbtable
     tableInfo.StructName = structName
     this.DbTableInfoByTableName[tableInfo.TableName] = &tableInfo
     this.DbTableInfoByStructName[tableInfo.StructName] = &tableInfo
  }
}