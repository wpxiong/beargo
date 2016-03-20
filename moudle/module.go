package moudle

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/constvalue"
  "github.com/wpxiong/beargo/util"
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

type ColumnInfo struct {
  ColumnName    string
  NotNull       bool
  SqlType       string
  IsId          bool
  Length        int
  Scale         int
  UniqueKey     string
  DefaultValue  interface{}
  FieldName     string
  FieldType     reflect.Type
}


type DBTableInfo struct {
  DbName        string
  DbSchema      string
  DbStuct       interface{}
  FiledNameMap  map[string] ColumnInfo
  DbTableExist  bool
}

type Moudle struct {
  DbDialect        DbDialectType
  DbName           string
  DbConnectionUrl  string
  DbUserName       string
  DbPassword       string
  DbTableInfo      map[string]DBTableInfo
  DbProiver        DbProviderInterface
  connectionStatus bool
}


func CreateModuleInstance(DbDialect DbDialectType,DbName string,DbConnectionUrl string, DbUserName string,DbPassword string) *Moudle {
   module :=  &Moudle{DbDialect:DbDialect,DbName:DbName,DbConnectionUrl:DbConnectionUrl,DbUserName:DbUserName,DbPassword:DbPassword}
   module.initModuleInstance()
   return module
}

func (this *Moudle) initModuleInstance(){
   this.DbTableInfo = make(map[string]DBTableInfo)
   var connectionUrl string  
   switch this.DbDialect {
      case MYSQL :
        connectionUrl = this.DbUserName  + ":" + this.DbPassword +   "@" + this.DbConnectionUrl + "/" + this.DbName;
        this.DbProiver = &MysqlDBProvider{}
      case POSTGRESQL :
   }
   err := this.DbProiver.ConnectionDb(connectionUrl)
   if err != nil {
      log.Error("DB Connection Error!")
   }else {
      this.connectionStatus = true
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


func (this *Moudle) addTable(dbtable interface{},tablename string,schemaname string){
  if !this.connectionStatus {
     return 
  }else {
     dbInfo := DBTableInfo{}
     dbname := strings.ToLower(reflect.TypeOf(dbtable).Name())
     fieldNum := reflect.TypeOf(dbtable).NumField()
     dbInfo.FiledNameMap = make(map[string]ColumnInfo)
     for i:=0;i<fieldNum;i++{
         field := reflect.TypeOf(dbtable).Field(i)
         id := field.Tag.Get(constvalue.DB_ID)
         columnInfo := ColumnInfo{}
         column_name := field.Tag.Get(constvalue.DB_COLUMN_NAME)
         not_null := field.Tag.Get(constvalue.DB_NOT_NULL)
         length := field.Tag.Get(constvalue.DB_LENGTH)
         scale := field.Tag.Get(constvalue.DB_SCALE)
         unique_key := field.Tag.Get(constvalue.DB_UNIQUE_KEY)
         default_value := field.Tag.Get(constvalue.DB_DEFAULT_VALUE)
         
         log.Debug(id)
         log.Debug(not_null)
         
         log.Debug(length)
         log.Debug(scale)
         
         log.Debug(unique_key)
         
         log.Debug(default_value)
         
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
            case reflect.Float32:
            case reflect.Float64:
            case reflect.Complex64:
            case reflect.Complex128:
            case reflect.Array:
            case reflect.Bool:
            case reflect.Ptr:
            case reflect.Struct:
            case reflect.Slice:
            case reflect.Map:
         }
         columnInfo.FieldType = field.Type
         columnInfo.FieldName = field.Name
         if !isEmpty(column_name) {
            columnInfo.ColumnName = column_name
         }else{
            columnInfo.ColumnName = field.Name
         }
         dbInfo.FiledNameMap[strings.ToLower(field.Name)] = columnInfo
     }
     if tablename == "" {
       dbInfo.DbName = dbname
     }else {
       dbInfo.DbName = tablename
     }
     if schemaname == "" {
       dbInfo.DbSchema = ""
     }else {
       dbInfo.DbSchema = schemaname
     }
     dbInfo.DbStuct = dbtable
     this.DbTableInfo[dbname] = dbInfo
     log.Debug(reflect.TypeOf(dbtable).Elem().Kind()) 
     log.Debug(this.DbTableInfo)
  }
}
