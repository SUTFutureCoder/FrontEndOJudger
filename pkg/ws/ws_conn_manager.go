package ws

import (
	"FrontEndOJudger/pkg/setting"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
)

type wsConnManager struct {
	c *websocket.Conn
}

func (c *wsConnManager) connAndTickTacker() {
	for {
		err := c.conn()
		if err != nil {
			log.Printf("Dail Front Websocket Error : %v. Retrying...", err)
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}

		// start tick tacker
		c.tickTacker()
	}
}

func (c *wsConnManager) conn() error {
	u := url.URL{Scheme: setting.FrontEndSetting.WebsocketScheme, Host: setting.FrontEndSetting.WebsocketHost, Path: fmt.Sprintf("/%s", setting.FrontEndSetting.WebsocketPath)}
	var dialer *websocket.Dialer
	var err error
	c.c, _, err = dialer.Dial(u.String(), http.Header{"session_token": []string{setting.FrontEndSetting.WebsocketToken}})
	return err
}

func (c *wsConnManager) tickTacker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <- ticker.C:
			req := &WsJsonReq{}
			req.Cmd = Ticker
			err := c.c.WriteJSON(req)
			//fmt.Println("send ticktacker")
			if err != nil {
				fmt.Printf("Websocket ticktacker Err: %v", err)
				return
			}
			continue
		}
	}
}