package net

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

func GetAndSaveFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	pix, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	written, err := io.Copy(out, bytes.NewReader(pix))
	if err != nil {
		return err
	}
	if written < 100 {
		// 小于100则视为请求失败
		return errors.New(fmt.Sprintf("get empty file from url written[%d]", written))
	}
	return nil
}
