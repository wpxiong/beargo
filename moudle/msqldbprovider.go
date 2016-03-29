package moudle

import (
  "github.com/wpxiong/beargo/log"
  "database/sql"
  "strconv"
  "strings"
  "reflect"
  _ "github.com/go-sql-driver/mysql"
)

func init() {
  log.InitLog()
}

type MysqlDBProvider struct {
   db *sql.DB
}


func (this *MysqlDBProvider ) Begin() (*sql.Tx,error) {
   return this.db.Begin()
}

func (this *MysqlDBProvider )  Commit(tx *sql.Tx) error {
   return tx.Commit()
}

func (this *MysqlDBProvider )  Close() error{
   return this.db.Close()
}

func (this *MysqlDBProvider ) ConnectionDb(dburl string) error {
  var err error = nil
  this.db, err = sql.Open("mysql", dburl)
  if err != nil {
    log.Error("Mysql Connection Error ")
  }
  return err
}

func (this *MysqlDBProvider ) Query(sql string) (*sql.Rows ,error){
   log.Info(sql)
   return this.db.Query(sql)
}

func (this *MysqlDBProvider ) Insert(sql string) (sql.Result ,error){
   log.Info(sql)
   return this.db.Exec(sql)
}

func (this *MysqlDBProvider ) CreateTable(tableName string,sqlstr string,primaryKey []string) (sql.Result ,error) {
   if len(primaryKey) != 0 {
      sqlstr = sqlstr[0:len(sqlstr)-2] +  ",\n PRIMARY KEY("  + strings.Join(primaryKey,",") + ") );"  
   }
   log.Info(sqlstr)
   return this.db.Exec(sqlstr)
}

func  (this *MysqlDBProvider ) CreateForeignKey(tableName string ,  keyColumn string, refrenceTableName string, referenceColumnName string) (sql.Result ,error) {
   foreignkeyname := keyColumn + "_" + refrenceTableName + "_" + referenceColumnName
   sqlstr := "alter table `" + tableName + "` add constraint `" + foreignkeyname  + "` foreign key (`"  + keyColumn + "`) references `" + refrenceTableName + "` (`" + referenceColumnName + "`) on delete cascade on update cascade "
   log.Info(sqlstr)
   return this.db.Exec(sqlstr)
}

func  (this *MysqlDBProvider ) CreatePrimaryKey(tableName string,keyList []string)  (sql.Result ,error) {
   if len(keyList) != 0 {
      sqlstr := "ALTER TABLE " + tableName  + " ADD PRIMARY KEY("  + strings.Join(keyList,",") + ")" 
      log.Info(sqlstr)
      return this.db.Exec(sqlstr)  
   }else {
      return nil,nil
   }
}

func (this *MysqlDBProvider ) ExecuteSQL(sql string) (sql.Result ,error){
   log.Info(sql)
   return this.db.Exec(sql)
}

func (this *MysqlDBProvider ) DropTable(tableName string) (sql.Result ,error){
   var sql string = "drop table if exists " + tableName + ";"
   log.Info(sql)
   return this.db.Exec(sql)
}

func (this *MysqlDBProvider ) GetDBIntType() string {
  return "INT"
}

func (this *MysqlDBProvider )  GetDBInt8Type() string {
   return "TINYINT"
}

func (this *MysqlDBProvider ) GetDBInt16Type() string {
   return "SMALLINT"
}

func (this *MysqlDBProvider )  GetDBInt32Type() string {
  return "SMALLINT"
}

func (this *MysqlDBProvider ) GetDBInt64Type() string {
  return "BIGINT"
}

func (this *MysqlDBProvider )  GetDBUintType() string {
  return "INT"
}

func (this *MysqlDBProvider )  GetDBUint8Type() string {
   return "TINYINT"
}

func (this *MysqlDBProvider )  GetDBUint16Type() string {
  return "SMALLINT"
}

func (this *MysqlDBProvider )  GetDBUint32Type() string {
    return "INT"
}

func (this *MysqlDBProvider )  GetDBUint64Type() string {
  return "BIGINT"
}

func (this *MysqlDBProvider )  GetDBFloat32Type() string {
  return "FLOAT"
}

func (this *MysqlDBProvider )  GetDBFloat64Type() string {
  return "DOUBLE"
}

func (this *MysqlDBProvider )  GetDBComplex64Type() string {
  return "VARCHAR(64)"
}

func (this *MysqlDBProvider )  GetDBComplex128Type() string {
  return "VARCHAR(128)"
}

func (this *MysqlDBProvider )  GetDBBoolType() string {
  return "CHAR(1)"
}

func (this *MysqlDBProvider )  GetDBStringType(length int ) string {
  if length < 65535 {
     return "VARCHAR"
  }else if length < 16777215 {
     return "MEDIUMTEXT"
  }else {
     return "LONGTEXT"
  }
}

func (this *MysqlDBProvider )  GetDBTimeType() string {
  return "TIMESTAMP"
}

func (this *MysqlDBProvider )  GetDBByteArrayType(length int) string {
  if length < 256 {
     return "TINYBLOB"
  }else if length < 65536 {
     return "BLOB"
  }else if length < 16777216 {
     return "MEDIUMBLOB"
  }else {
     return "LONGBLOB"
  }
}

func (this *MysqlDBProvider )  LimitSql( limit int ) string {
   return " LIMIT "  + strconv.Itoa(limit)
}

func (this *MysqlDBProvider )  CreateDefaultValue(defaultValue interface{}) string {
  if defaultValue == nil {
     return ""
  }
  t := reflect.TypeOf(defaultValue)
  switch t.Name() {
     case "int64":
       val := defaultValue.(int64) 
       return " DEFAULT " + strconv.FormatInt(val,10) + " "
     case "uint64":
       val := defaultValue.(uint64)
       return " DEFAULT " + strconv.FormatUint(val,10) + " "
     case "float64":
       val := defaultValue.(float64) 
       return " DEFAULT " + strconv.FormatFloat(val, 'f', -1, 64) + " "
     case "complex64":   
       val := defaultValue.(complex64)
       return " DEFAULT '" + strconv.FormatFloat(float64(real(val)), 'f', -1, 32) + "," +  strconv.FormatFloat(float64(imag(val)), 'f', -1, 32) +  "' "
     case "complex128":
       val := defaultValue.(complex128) 
       return " DEFAULT '" + strconv.FormatFloat(real(val), 'f', -1, 64) + "," +  strconv.FormatFloat(imag(val), 'f', -1, 64) +  "' "
     case "bool":
       val := defaultValue.(bool) 
       if val {
          return " DEFAULT '" + "1" + "' "
       }else {
          return " DEFAULT '" + "0" + "' "
       }
     case "string":
       val := defaultValue.(string)
       if val == "" {
          return ""
       }else {
          return " DEFAULT '" + defaultValue.(string) + "' "
       }
     default :
       return ""
  }
  return ""
}


func (this *MysqlDBProvider ) CreateSqlTypeByLength(auto_increment bool ,sqlType string,length int, scale int) string {
   result := sqlType
   switch sqlType {
      case "INT","TINYINT","SMALLINT","BIGINT","VARCHAR" :
        if length != 0 {
          result = sqlType + "(" + strconv.Itoa(length)  + ")"
        }
      case "FLOAT","DOUBLE":
        if length != 0 && scale != 0 {
          result = sqlType + "(" + strconv.Itoa(length)  + ","  +  strconv.Itoa(scale) + ")"
        } else if scale != 0 {
          result = sqlType + "(" + strconv.Itoa(scale)  + ","  +  strconv.Itoa(scale) + ")"
        } else if length != 0 {
          result = sqlType + "(" + strconv.Itoa(length)  + ","  +  strconv.Itoa(length) + ")"
        }
      default:
        result = sqlType
   }
   if auto_increment {
     result += " AUTO_INCREMENT "
   }
   return result
}

