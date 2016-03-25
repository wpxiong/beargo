package moudle

import (
  "github.com/wpxiong/beargo/log"
  "database/sql"
)

func init() {
  log.InitLog()
}

type DbProviderInterface interface {
   ConnectionDb(dburl string) error 
   Query(sql string) (*sql.Rows ,error)
   Insert(sql string) (sql.Result ,error)
   CreateTable(tableName string , sqlstr string,primaryKey []string) (sql.Result ,error)
   DropTable(tableName string) (sql.Result ,error)
   ExecuteSQL(sql string) (sql.Result ,error)
   CreatePrimaryKey(tableName string,keyList []string)  (sql.Result ,error)
   Begin() (*sql.Tx,error)
   Close() error
   Commit(tx *sql.Tx) error
   GetDBIntType() string
   GetDBInt8Type() string
   GetDBInt16Type() string
   GetDBInt32Type() string
   GetDBInt64Type() string
   GetDBUintType() string
   GetDBUint8Type() string
   GetDBUint16Type() string
   GetDBUint32Type() string
   GetDBUint64Type() string
   GetDBFloat32Type() string
   GetDBFloat64Type() string
   GetDBComplex64Type() string
   GetDBComplex128Type() string
   GetDBBoolType() string
   GetDBStringType() string
   GetDBTimeType() string
   CreateSqlTypeByLength(auto_increment bool , sqlType string,length int, scale int) string
   CreateDefaultValue(defaultValue interface{}) string
}

