package moudle

import (
  "github.com/wpxiong/beargo/log"
  "github.com/wpxiong/beargo/constvalue"
  "database/sql"
  "strconv"
  "strings"
  "reflect"
  "time"
  "encoding/hex"
  _ "github.com/go-sql-driver/mysql"
)

func init() {
  log.InitLog()
}

type MysqlDBProvider struct {
   db *sql.DB
   
}

func (this *MysqlDBProvider )  SetMinConnection(max int) {
   this.db.SetMaxIdleConns(max)
}

func (this *MysqlDBProvider )  SetMaxConnection(max int) {
   this.db.SetMaxOpenConns(max)
}

func (this *MysqlDBProvider ) Begin() (*sql.Tx,error) {
   return this.db.Begin()
}

func (this *MysqlDBProvider )  Commit(tx *sql.Tx) error {
   return tx.Commit()
}

func (this *MysqlDBProvider )  Close() error {
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

func (this *MysqlDBProvider ) Query(sql string,ts *Trans) (*sql.Rows ,error){
   log.Info(sql)
   if ts == nil {
     return this.db.Query(sql)
   }else {
     return ts.tx.Query(sql)
   }
}

func (this *MysqlDBProvider ) Insert(sql string,ts *Trans) (sql.Result ,error){
   log.Info(sql)
   if ts == nil {
     return this.db.Exec(sql)
   }else {
     return ts.tx.Exec(sql)
   }
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

func (this *MysqlDBProvider ) ExecuteSQL(sql string,ts *Trans) (sql.Result ,error){
   log.Info(sql)
   if ts == nil {
     return this.db.Exec(sql)
   }else {
     return ts.tx.Exec(sql)
   }
}

func (this *MysqlDBProvider ) Update(sql string,ts *Trans) (sql.Result ,error){
   log.Info(sql)
   if ts == nil {
     return this.db.Exec(sql)
   }else {
     return ts.tx.Exec(sql)
   }
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

func (this *MysqlDBProvider ) GetInsertByteDataSql(byteval []byte ) string {
  return  "0x" + hex.EncodeToString(byteval)
}

func (this *MysqlDBProvider ) GetInsertDBComplex64Sql(val complex128 ) string {
  return "'" + strconv.FormatFloat(real(val),'f', -1, 32) +"," + strconv.FormatFloat(imag(val),'f', -1, 32) + "'"
}

func (this *MysqlDBProvider ) GetInsertDBComplex128Sql(val complex128 ) string {
   return "'" + strconv.FormatFloat(real(val),'f', -1, 64) +"," + strconv.FormatFloat(imag(val),'f', -1, 64) + "'"
}

func (this *MysqlDBProvider )   GetInsertDBTimeSql(ti time.Time) string {
   return "'" + ti.Format(constvalue.DEFAULT_TIME_FORMATE) + "'"
}

func (this *MysqlDBProvider )   AppendScanComplexField(list *[]interface{}) {
    var complexField string 
    (*list) = append(*list,&complexField)
}

func (this *MysqlDBProvider )  PrepareExecuteSQL(sql string ,parameter []interface{},ts *Trans) {
 log.Debug(sql)
 if ts == nil {
    if smt,err := this.db.Prepare(sql);err == nil {
       if _, err1 := smt.Exec(parameter...);err1 != nil {
          panic(err)
       }
    }else {
       panic(err)
    }
 }else {
    if smt,err := this.db.Prepare(sql);err == nil {
       smt = ts.tx.Stmt(smt)
       if _, err1 := smt.Exec(parameter...);err1 != nil {
          panic(err)
       }
    }else {
       panic(err)
    }
 }
}

   
func (this *MysqlDBProvider )  TableExistsInDB(tableName string) (bool,error) {
   var sql string = "show tables like '" + tableName + "';"
   log.Info(sql)
   rows,err := this.db.Query(sql)
   if err != nil {
      return false,err
   }else {
      return rows.Next(),err
   }
}

func (this *MysqlDBProvider ) GetCurrentSerialValue(tableName string, columnName string,ts *Trans) int64 {
   sql := "select last_insert_id()"
   log.Info(sql)
   if ts == nil {
     rows ,err := this.db.Query(sql)
     if err == nil {
       if rows.Next() {
          var res int64 
          err = rows.Scan(&res)
          if err == nil {
            return res
          }else {
            panic(err)
          }
       }
     }else {
        panic(err)
     }
   }else {
     rows ,err := ts.tx.Query(sql)
     if err == nil {
       if rows.Next() {
          var res int64 
          err = rows.Scan(&res)
          if err == nil {
            return res
          }else {
            panic(err)
          }
       }
     }else {
        panic(err)
     }
   }
   panic("get current serial value failture")
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

