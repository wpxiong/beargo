package moudle

import (
  "github.com/wpxiong/beargo/log"
  "database/sql"
  "errors"
)

func init() {
  log.InitLog()
}


type Trans struct {
  tx *sql.Tx
  moudle *Moudle
}

func (this *Moudle) Beign() *Trans {
  var tr *Trans = &Trans{moudle:this}
  if tx,err := this.DbProiver.Begin();err == nil {
    tr.tx = tx
    return tr
  }else {
    panic(err)
  }
}

func (this *Trans) Commit() error {
  if this.tx != nil {
     return this.tx.Commit()
  }
  return errors.New("No DB transaction")
}

func (this *Moudle) Close() error {
  return this.DbProiver.Close()
}

func (this *Trans) Rollback() error {
  if this.tx != nil {
     return this.tx.Rollback()
  }
  return errors.New("No DB transaction")
}


func (this *DbSqlBuilder) DeleteExecute() {
   this.deleteExecute(this.ts)
}


func (this *Trans)  DeleteWithSql(sql string,parameter []interface{}) {
   this.moudle.deleteWithSql(sql,parameter,this)
}


func (this *Trans) InsertWithSql(sql string,parameter []interface{})  {
  this.moudle.insertWithSql(sql,parameter,this)
}


func (this *Trans) UpdateWithSql(sql string,parameter []interface{}) {
   this.moudle.prepareExecuteSQL(sql,parameter,this)
}

func (this *Trans) ExecuteWithSql(sql string,parameter []interface{}) {
   this.moudle.executeWithSql(sql,parameter,this)
}

func (this *Trans) DeleteWhere(structVal interface{}) *DbSqlBuilder {
   sqlBuilder := this.moudle.DeleteWhere(structVal)
   sqlBuilder.ts = this
   return sqlBuilder
}

func (this *Trans) UpdateWithWhere(structVal interface{},fieldName[]string)  *DbSqlBuilder {
   sqlBuilder := this.moudle.UpdateWithWhere(structVal,fieldName)
   sqlBuilder.ts = this
   return sqlBuilder
}

func (this *DbSqlBuilder) InsertExecute() {
   this.insertExecute(this.ts)
}

func (this *DbSqlBuilder) SaveExecute() {
   this.saveExecute(this.ts)
}

func (this *Trans) SimpleQuery(tableObj interface{}) *DbSqlBuilder {
   sqlBuilder := this.moudle.SimpleQuery(tableObj)
   sqlBuilder.ts = this
   return sqlBuilder
}

func (this *Trans) Query(tableObj interface{},fetchType FetchType,structName []string) *DbSqlBuilder {
   sqlBuilder := this.moudle.Query(tableObj,fetchType,structName)
   sqlBuilder.ts = this
   return sqlBuilder
}

func (this *Trans) Insert(structVal interface{} ) *DbSqlBuilder {
   sqlBuilder := this.moudle.Insert(structVal)
   sqlBuilder.ts = this
   return sqlBuilder
}

func (this *Trans)  Delete(structVal interface{} ) {
  this.moudle.delete(structVal,this)
}