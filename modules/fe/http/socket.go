package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"strings"

	"github.com/Cepave/open-falcon-backend/modules/fe/http/base"
	event "github.com/Cepave/open-falcon-backend/modules/fe/model/falcon_portal"
	"github.com/Cepave/open-falcon-backend/modules/fe/model/uic"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type ScoketController struct {
	base.BaseController
}

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) (err error) {
	log.Println("echo 12345")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	currentTime := time.Now().Unix()
	// control := make(chan string)
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		messageStr := string(message)
		if messageStr == "0" || err != nil {
			break
		}
		sigMsg := strings.Split(messageStr, ";")
		if len(sigMsg) != 2 {
			log.Printf("error")
			break
		}
		sig := sigMsg[1]
		name := sigMsg[0]
		uic.ReadSessionBySig(sig)
		session := uic.ReadSessionBySig(sig)
		if session.Uid != uic.SelectUserIdByName(name) {
			log.Println("can not find this kind of session")
			break
		}
		log.Printf("recv: %s", message)

		go func() {
			for {
				res, err := event.GetEventCases(currentTime, time.Now().Unix(), -1, "ALL", "ALL", 10, 0, "root", "ALL", "")
				if err != nil {
					break
				}
				err = c.WriteMessage(mt, []byte(strconv.Itoa(int(currentTime))))
				resjsB, _ := json.Marshal(res)
				err = c.WriteMessage(mt, resjsB)
				if err != nil {
					break
				}
				currentTime = time.Now().Unix()
				time.Sleep(30000 * time.Millisecond)
			}
		}()
	}
	return
}

func (this *ScoketController) WebScoket() {
	echo(this.Ctx.ResponseWriter, this.Ctx.Request)
}

func (this *ScoketController) TestInd() {
	this.TplName = "home/socket.html"
}

func MyRoutesing() {
	beego.Router("/echo", &ScoketController{}, "get:WebScoket")
	beego.Router("/sss", &ScoketController{}, "get:TestInd")
}
