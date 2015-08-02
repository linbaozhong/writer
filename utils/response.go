package utils

import (
//"errors"
//"fmt"
//"reflect"
)

type ResultData struct {
	Key  interface{} `json:"key"`
	Data interface{} `json:"data"`
}
type Result struct {
	Status interface{} `json:"ok"`
	ResultData
}

//返回JSON格式消息
func JsonResult(ok bool, key string, data interface{}) *Result {
	var r = new(Result)
	r.Status = ok
	r.Key = key
	r.Data = data
	return r
}

//返回JSON格式消息
func ActionResult(args ...interface{}) *Result {
	var r = new(Result)
	var n = len(args)

	if n > 0 {
		r.Status = args[0]
		if n > 1 {
			r.Data = args[1]
		} else {
			r.Data = ""
		}
	} else {
		r.Status = false
		r.Data = ""
	}

	return r
}
