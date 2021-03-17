package ws

import (
	"FrontEndOJudger/pkg/setting"
	"log"
)

var connManager wsConnManager

func Setup() {
	if !setting.FrontEndSetting.EnableWebsocket {
		return
	}

	// ensure alive
	go connManager.connAndTickTacker()
	go connManager.handleSend()
}


type WsJsonReq struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

var SendChan chan *WsJsonReq

const (
	Ticker = "Ticker"
	JudgerResultCallBack = "JudgerResultCallBack"
)

func (c *wsConnManager)handleSend() {
	SendChan = make(chan *WsJsonReq, 128)
	for {
		select {
		case sendData := <- SendChan:
			if c.c == nil {
				log.Printf("Websocket connection not ready. Reset data")
				SendChan <- sendData
				continue
			}
			err := c.c.WriteJSON(sendData)
			if err != nil {
				log.Printf("SEND Msg to Websocket Server Error err:%v", err)
				// wait reconnect
				SendChan <- sendData
			}
		}
	}
}