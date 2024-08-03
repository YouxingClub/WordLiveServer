package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"log"
	"net/http"
	"slices"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	m := melody.New()
	r := gin.Default()
	r.Static("/echo-live", "static/echo-live/echo-live")
	//fmt.Println(echoliveStaticFS.ReadDir("static/echo-live/echo-live"))
	r.LoadHTMLGlob("static/echo-live/*.html")
	var devicesList []Device

	r.GET("/live", func(c *gin.Context) {
		c.HTML(http.StatusOK, "live.html", gin.H{})
	})

	r.GET("/history", func(c *gin.Context) {
		c.HTML(http.StatusOK, "history.html", gin.H{})
	})

	r.GET("/settings", func(c *gin.Context) {
		c.HTML(http.StatusOK, "settings.html", gin.H{})
	})

	r.GET("/editor", func(c *gin.Context) {
		c.HTML(http.StatusOK, "editor.html", gin.H{})
	})

	r.GET("/ws", func(c *gin.Context) {
		err := m.HandleRequest(c.Writer, c.Request)
		if err != nil {
			log.Println("HandleRequest error:", err)
			return
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		log.Println("存在设备接入：", s.RemoteAddr())
		//newDeviceRegData := NewEchoAPI()
		//json, _ := newDeviceRegData.Marshal()
		//err := s.Write(json)
		//if err != nil {
		//	log.Println("发送注册信息失败：", err)
		//	return
		//}
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Println("设备断开连接", s.RemoteAddr())
		for index, device := range devicesList {
			if device.Addr == s.RemoteAddr() {
				log.Println("设备删除：", device.Addr, device.UUID)
				slices.Delete(devicesList, index, index+1)
				break
			}
		}
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		log.Println("收到新WS消息：", string(msg))
		echoAPIData := NewEchoAPI()
		err := echoAPIData.Unmarshal(msg)
		if err != nil {
			log.Println("解析消息失败：", err)
			return
		}
		switch echoAPIData.Action {
		case "hello":
			log.Println("收到设备注册信息：", s.RemoteAddr(), echoAPIData.From.UUID)
			devicesList = append(devicesList, Device{
				Name:      echoAPIData.From.Name,
				UUID:      echoAPIData.From.UUID,
				Type:      echoAPIData.From.Type,
				Timestamp: echoAPIData.From.Timestamp,
				Addr:      s.RemoteAddr(),
			})
		}
		err = m.Broadcast(msg)
		if err != nil {
			log.Println("广播消息失败：", err)
			return
		}
	})

	err := r.Run(":3000")
	if err != nil {

		log.Fatal("ListenAndServe: ", err)
		return
	}
}
