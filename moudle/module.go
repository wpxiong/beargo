package moudle

import (
  "github.com/wpxiong/beargo/log"
<<<<<<< HEAD
  "github.com/wpxiong/beargo/constvalue"
=======
  "github.com/wpxiong/beargo/util/dbutil"
>>>>>>> origin/master
  "reflect"
  "strings"
)

func init() {
  log.InitLog()
}


type DbDialectType int
<<<<<<< HEAD

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
  DefaultValue  string
}
=======
>>>>>>> origin/master

const (
    MYSQL   DbDialectType = iota
    SQLITE
    POSTGRESQL
)

type DBTableInfo struct {
  DbName        string
  DbSchema      string
  DbStuct       interface{}
  FiledList     map[string] reflect.Type
<<<<<<< HEAD
  FiledNameMap  map[string] ColumnInfo
=======
  FiledNameMap  map[string] string
>>>>>>> origin/master
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
<<<<<<< HEAD
=======
}

type DBTableInterface interface {
   GetDbTableName() string
>>>>>>> origin/master
}

type DBTableInterface interface {
   GetDbTableName() string
}


<<<<<<< HEAD
=======
type DBTable struct {

}

func (this *DBTable) GetDbTableName() string {
   return reflect.TypeOf(*this).Name()
}


>>>>>>> origin/master
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


<<<<<<< HEAD
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

func (this *Moudle) addTable(dbtable interface{},tablename string,schemaname string){
=======

func (this *Moudle) AddTable(dbtable interface{}){
>>>>>>> origin/master
  if !this.connectionStatus {
     return 
  }else {
     dbInfo := DBTableInfo{}
     dbname := strings.ToLower(reflect.TypeOf(dbtable).Name())
     fieldNum := reflect.TypeOf(dbtable).NumField()
     dbInfo.FiledList = make(map[string]reflect.Type)
<<<<<<< HEAD
     dbInfo.FiledNameMap = make(map[string]ColumnInfo)
     for i:=0;i<fieldNum;i++{
         field := reflect.TypeOf(dbtable).Field(i)
         id := field.Tag.Get(constvalue.DB_ID)
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
            case reflect.Int8:
            case reflect.Int16:
            case reflect.Int32:
            case reflect.Int64:
            case reflect.Uint:
            case reflect.Uint8:
            case reflect.Uint16:
            case reflect.Uint32:
            case reflect.Uint64:
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
         dbInfo.FiledList[field.Name] = field.Type
         if !isEmpty(column_name) {
           field.Name = column_name
         }
         dbInfo.FiledNameMap[strings.ToLower(field.Name)] = ColumnInfo{ColumnName:field.Name}
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
=======
     dbInfo.FiledNameMap = make(map[string]string)
     for i:=0;i<fieldNum;i++{
         field := reflect.TypeOf(dbtable).Field(i)
         dbInfo.FiledList[field.Name] = field.Type
         dbInfo.FiledNameMap[strings.ToLower(field.Name)] = field.Name
     }
     dbInfo.DbName = dbname
     dbInfo.DbSchema = ""
     dbInfo.DbStuct = dbtable
     this.DbTableInfo[dbname] = dbInfo
     dbutil.GetCreateTableSql()
     log.Debug(reflect.TypeOf(dbtable).Elem().Kind()) 
     this.DbTableInfo[dbtable.GetDbTableName()] = dbInfo
     log.Debug(this.DbTableInfo)
>>>>>>> origin/master
  }

}