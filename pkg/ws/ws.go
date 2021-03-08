package ws

import (
	"FrontEndOJudger/pkg/setting"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

func Setup() {
	if !setting.FrontEndSetting.EnableWebsocket {
		return
	}
	u := url.URL{Scheme: setting.FrontEndSetting.WebsocketScheme, Host: setting.FrontEndSetting.WebsocketHost, Path: fmt.Sprintf("/%s", setting.FrontEndSetting.WebsocketPath)}
	var dialer *websocket.Dialer
	c, _, err := dialer.Dial(u.String(), http.Header{"session_token": []string{setting.FrontEndSetting.WebsocketToken}})
	if err != nil {
		log.Fatalf("Dail Front Websocket Error : %v", err)
	}

	handleSend(c)
}


type WsJsonReq struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

var SendChan chan *WsJsonReq

const (
	HelloWorld = "HelloWorld"
	SayHello   = "SayHello"
	JudgerResultCallBack = "JudgerResultCallBack"
)

func handleSend(c *websocket.Conn) {
	SendChan = make(chan *WsJsonReq, 128)
	for {
		select {
		case sendData := <- SendChan:
			err := c.WriteJSON(sendData)
			if err != nil {
				log.Printf("SEND Msg to Websocket Server Error err:%v", err)
			}
		}
	}
}