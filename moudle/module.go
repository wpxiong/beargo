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
  Unique        bool
  DefaultValue  interface{}
  FieldName     string
  FieldType     reflect.Type
  RelationTable *DBTableInfo
  RelationStructName  string
  IsArray       bool
}


type DBTableInfo struct {
  TableName        string
  DbSchema      string
  DbStuct       interface{}
  FiledNameMap  map[string] ColumnInfo
  DbTableExist  bool
  StructName    string
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


func (this *Moudle) droptable(tablename string) {
   _,err := this.DbProiver.DropTable(tablename)
   if err != nil {
      log.Error(err)
   }
}

func (this *Moudle) createtable(sqlstr string) {
   _,err := this.DbProiver.CreateTable(sqlstr)
   if err != nil {
      log.Error(err)
   }
}

func (this *Moudle)  InitialDB(create bool) {
  log.Debug("Initial DB start")
  if create {
    //create Table
    var index int = 0
    for key,Info := range this.DbTableInfo {
      this.droptable(key)
      var create_sql string = "create table " + key + " ( "
      for _,column := range Info.FiledNameMap {
        if column.RelationStructName != "" {
           continue
        }
        columnstr := column.ColumnName + " " +  this.createSqlTypeByLength(column.SqlType,column.Length,column.Scale)
        if column.NotNull {
           columnstr += " not null "
        }
        if column.Unique {
           columnstr += " unique "
        }
        columnstr += this.createDefaultValue(column.DefaultValue)
        columnstr += ",\n"
        create_sql += columnstr
        index ++
      }
      create_sql = create_sql[0:len(create_sql)-2]
      create_sql += "\n)"
      this.createtable(create_sql)
    }
  }else {
    //check Table is exist in DB
    
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

func (this *Moudle)  getDBStringType() string {
  return this.DbProiver.GetDBStringType()
}

func (this *Moudle)  getDBTimeType() string {
  return this.DbProiver.GetDBTimeType()
}



func (this *Moudle)  createSqlTypeByLength(sqltype string, length int,scale int) string {
  return this.DbProiver.CreateSqlTypeByLength(sqltype,length,scale)
}

func (this *Moudle)  createDefaultValue(defaultValue interface{}) string {
  return this.DbProiver.CreateDefaultValue(defaultValue)
}

func (this *Moudle) addTable(dbtable interface{},tableName string,schemaname string){
  if !this.connectionStatus {
     return 
  }else {
     tableInfo := DBTableInfo{}
     tablenamestr := strings.ToLower(reflect.TypeOf(dbtable).Name())
     fieldNum := reflect.TypeOf(dbtable).NumField()
     tableInfo.FiledNameMap = make(map[string]ColumnInfo)
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
              // foreign key
              columnInfo.RelationStructName = field.Type.Elem().Name()
              columnInfo.IsArray = true
            case reflect.Map:
              continue
            case reflect.String:
              columnInfo.SqlType = this.getDBStringType()
              columnInfo.DefaultValue = default_value
         }
         columnInfo.FieldType = field.Type
         columnInfo.FieldName = field.Name
         if !isEmpty(column_name) {
            columnInfo.ColumnName = column_name
         }else{
            columnInfo.ColumnName = field.Name
         }
         
         if strings.ToLower(strings.Trim(id," ")) == "true" {
            columnInfo.IsId = true
         }
         
         if strings.ToLower(strings.Trim(not_null," ")) == "true" {
            columnInfo.NotNull = true
         }
         
         if len,err := util.GetIntValue(scale); err == nil {
             columnInfo.Scale = int(len)
         }
         
         if len,err := util.GetIntValue(length); err == nil {
             columnInfo.Length = int(len)
         }
         
         if strings.ToLower(strings.Trim(unique_key," ")) == "true" {
            columnInfo.Unique = true
         }
         tableInfo.FiledNameMap[strings.ToLower(field.Name)] = columnInfo
     }
     if tableName == "" {
       tableInfo.TableName = tablenamestr
     }else {
       tableInfo.TableName = tableName
     }
     if schemaname == "" {
       tableInfo.DbSchema = ""
     }else {
       tableInfo.DbSchema = schemaname
     }
     tableInfo.DbStuct = dbtable
     tableInfo.StructName = tablenamestr
     this.DbTableInfo[tableInfo.TableName] = tableInfo
  }
}