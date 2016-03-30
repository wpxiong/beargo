package moudle

import (
  "github.com/wpxiong/beargo/log"
  "database/sql"
  "time"
)

func init() {
  log.InitLog()
}

type DbProviderInterface interface {
   ConnectionDb(dburl string) error 
   Query(sql string) (*sql.Rows ,error)
   Insert(sql string) (sql.Result ,error)
   Update(sql string) (sql.Result ,error)
   CreateTable(tableName string , sqlstr string,primaryKey []string) (sql.Result ,error)
   DropTable(tableName string) (sql.Result ,error)
   ExecuteSQL(sql string) (sql.Result ,error)
   CreatePrimaryKey(tableName string,keyList []string)  (sql.Result ,error)
   CreateForeignKey(tableName string ,  keyColumn string, refrenceTableName string, referenceColumnName string) (sql.Result ,error)
   Begin() (*sql.Tx,error)
   Close() error
   Commit(tx *sql.Tx) error
   LimitSql(limit int) string
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
   GetDBStringType(length int ) string
   GetDBTimeType() string
   GetDBByteArrayType(length int ) string
   CreateSqlTypeByLength(auto_increment bool , sqlType string,length int, scale int) string
   CreateDefaultValue(defaultValue interface{}) string
   GetInsertDBComplex64Sql(val complex128 ) string
   GetInsertDBComplex128Sql(val complex128 ) string
   GetInsertDBTimeSql(time.Time) string
}

