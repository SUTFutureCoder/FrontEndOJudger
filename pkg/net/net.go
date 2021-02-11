package net

import (
	"encoding/json"
	"net/http"
)

const (
	SUCCESS           = 200
	ERROR             = 500
	UNAUTHORIZED      = 401
	NOT_LOGINED       = 402
	INVALID_PARAMS    = 406
	TOO_MANY_REQUESTS = 429
)


type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Writer(w http.ResponseWriter, code int, msg string, data interface{}) {
	var resp Response
	resp.Code = code
	resp.Msg = msg
	resp.Data = data
	jsonResp, _ := json.Marshal(resp)
 	w.Write(jsonResp)
}

