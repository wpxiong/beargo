package util

import (
  "strconv"
  "strings"
  "errors"
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

func GetFloat32Value(str string) (float64,error) {
   return strconv.ParseFloat(str, 32)
}

func GetFloat64Value(str string) (float64,error) {
   return strconv.ParseFloat(str, 64)
}

func GetComplex64Value(str string) (complex64,error) {
   com := strings.Split(str,",")
   var err error
   var value1,value2 float64
   if len(com) == 2 {
      value1,err = strconv.ParseFloat(com[0], 32)
      if err == nil {
         value2,err = strconv.ParseFloat(com[1], 32)
         if err == nil {
            var v1 float32 = float32(value1)
            var v2 float32 = float32(value2)
            return complex(v1,v2),nil
         }
      }
   }else {
     err = errors.New("string to complex64 is error")
   }
   return complex(0,0),err
}

func GetComplex128Value(str string) (complex128,error) {
   com := strings.Split(str,",")
   var err error
   var value1,value2 float64
   if len(com) == 2 {
      value1,err = strconv.ParseFloat(com[0], 64)
      if err == nil {
         value2,err = strconv.ParseFloat(com[1], 64)
         if err == nil {
            return complex(value1, value2),nil
         }
      }
   }else {
     err = errors.New("string to complex64 value is error")
   }
   return complex(0,0),err
}

func GetBoolValue(str string) (bool,error) {
   if str == "1" {
     return true,nil
   }else if str == "0" {
     return false,nil
   }
   return true,errors.New("string to bool value is error")
}


