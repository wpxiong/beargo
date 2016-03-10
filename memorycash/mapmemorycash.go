package memorycash

import (
  "github.com/wpxiong/beargo/log"
  "reflect"
  "errors"
)

func init() {
  log.InitLog()
}


type MapMemoryCash struct {
   valueMap  map[string]interface{}
   count     int64
   maxcount  int64
}


func (this *MapMemoryCash) AddStringObject(key string, val string) error {
   if this.count < this.maxcount {
      this.valueMap[key] = reflect.ValueOf(val).Interface()
      this.count ++
      return nil
   }else {
      return errors.New("memory capacity is full")
   }
}

func (this *MapMemoryCash) GetStringObject(key string) (string,error) {
  var val interface{} = this.valueMap[key] 
  switch reflect.TypeOf(val).Kind() {
     case reflect.String:
        if val != nil {
           return val.(string),nil
        }else {
           return "",errors.New("has no object")
        }
     default:
        return "",errors.New("object is not string")
  }
}

func (this *MapMemoryCash) GetObject(key string) (interface{}, error) {
  var val interface{} = this.valueMap[key] 
  if val != nil {
      return val.(string),nil
  }else {
      return "",errors.New("has no object")
  }
}

func (this *MapMemoryCash) AddObject(key string, val interface{}) error {
  if this.count < this.maxcount {
     this.valueMap[key] = val
     this.count ++
     return nil
  } else {
     return errors.New("memory capacity is full")
  }
}

func (this *MapMemoryCash) DeleteObject(key string) error {
  delete(this.valueMap,key)
  this.count --
  return nil
}

func (this *MapMemoryCash) InitMemoryCash(size int64) error {
  this.valueMap = make(map[string]interface{},size)
  this.maxcount = size 
  this.count = 0
  return nil
}