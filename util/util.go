package util

import (
  "strconv"
  "strings"
  "github.com/wpxiong/beargo/log"
)

func init() {
  log.InitLog()
}

func StringToInt(valStr string) interface{}  {
   i, err := strconv.Atoi(valStr)
   if err != nil {
      return i
   }else {
      return nil
   }
}

func StringToFloat(valStr string) interface{} {
   i, err := strconv.ParseFloat(valStr, 64)
   if err != nil {
      return i
   }else {
      return nil
   }
}


func StringToDouble(valStr string) interface{} {
   i, err := strconv.ParseFloat(valStr, 64)
   if err != nil {
      return i
   }else {
      return nil
   }
}

func StringToBool(valStr string) interface{} {
  if strings.ToLower(valStr) == "true" {
     return true
  }else if strings.ToLower(valStr) == "false" {
     return false
  }else {
     return nil
  }
}


func Redirect(errorInfo string) {

}


func GetUint64Value(str string) (uint64,error) {
   return strconv.ParseUint(str, 10, 64)
}

func GetUint32Value(str string) (uint64,error) {
  return strconv.ParseUint(str, 10, 32)
}

func GetUint16Value(str string) (uint64,error) {
  return strconv.ParseUint(str, 10, 16)
}

func GetUint8Value(str string) (uint64,error) {
  return strconv.ParseUint(str, 10, 8)
}

func GetUintValue(str string) (uint64,error) {
  return strconv.ParseUint(str,10,32)
}


func GetInt64Value(str string) (int64,error) {
   return strconv.ParseInt(str, 10, 64)
}

func GetInt32Value(str string) (int64,error) {
  return strconv.ParseInt(str, 10, 32)
}

func GetInt16Value(str string) (int64,error) {
  return strconv.ParseInt(str, 10, 16)
}

func GetInt8Value(str string) (int64,error) {
  return strconv.ParseInt(str, 10, 8)
}

func GetIntValue(str string) (int64,error) {
   return strconv.ParseInt(str, 10, 32)
}

