package handle

import (
	"encoding/xml"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"youtube/domain"
	"youtube/store"
)

type Controller struct {
	Store *store.Store
}

var once sync.Once
var instance *http.Client

func NewController() *Controller {
	return &Controller{Store: store.NewStore()}
}

func GetClient() *http.Client {
	once.Do(func() {
		instance = &http.Client{}
	})
	return instance
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 1024
)

var upgrader = &websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

type Psocket struct {
	Socket *websocket.Conn
	mu     sync.Mutex
}

func (p *Psocket) SendToClient(v []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Socket.WriteMessage(websocket.TextMessage, v)
}

func (c *Controller) AddVideoFromYoutube(w http.ResponseWriter, r *http.Request) {
	listchannel := []string{"UCPtmylrUkGoDkAAWMaUH91A", "UC_x5XG1OV2P6uZZ5FSM9Ttw", "UCsT0YIqwnpJCM-mx7-gSA4Q", "UCC3L8QaxqEGUiBC252GHy3w", "UCiaHN4XRxNHjP84tlSAmDyQ", "UCAuUUnT6oDeKwE6v1NGQxug", "UCgXbCVJ79-JVyHoBIDhpvEQ", "UClYb9NpXnRemxYoWbcYANsA", "UCYxRlFDqcWM4y7FfpiAN3KQ"} //,
	//==================concurrent====================
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := new(sync.WaitGroup)

	//=========================socket==============================
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		return
	}
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("WebSoc:", err)
		return
	}

	p := &Psocket{}
	p.Socket = socket

	defer p.Socket.Close()
	//================================================
	var IsFound bool
	for {
		_, msg, err := p.Socket.ReadMessage()
		IsFound = false
		log.Println(string(msg))

		if err != nil {
			log.Println("nil")
			return
		} else {
			client := GetClient()
			for i := 0; i < len(listchannel); i++ {

				wg.Add(1)
				go c.Goconc(msg, listchannel[i], wg, client, p, &IsFound)

			}
		}
		wg.Wait()
		log.Println(IsFound)
		if !IsFound {
			log.Println("NOT FOUND")
			err1 := p.SendToClient([]byte("NOT FOUND"))
			if err1 != nil {
				return
			}
		}
	}

}

func (c *Controller) Goconc(msg []byte, str string, wg *sync.WaitGroup, client *http.Client, psocket *Psocket, IsFound *bool) {
	defer wg.Done()

	result1, err1 := c.Store.GetResultVideos(str)
	if err1 != nil {
		panic(err1)
	}
	log.Println("a")
	for _, result := range result1 {
		wg.Add(1)
		go EachVideo(msg, result, wg, client, psocket, IsFound)
	}
}

func EachVideo(msg []byte, result *domain.YoutubeSearchResult, wg *sync.WaitGroup, client *http.Client, psocket *Psocket, IsFound *bool) {
	defer wg.Done()
	url := "http://video.google.com/timedtext?lang=en&v=" + result.YoutubeId

	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(req)

	a := domain.Video{}
	a.VideoId = result.YoutubeId

	data, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(data, &a)
	log.Println("0")
	for i, _ := range a.Content {
		log.Println("1")
		matched, _ := regexp.MatchString(" "+strings.ToLower(string(msg))+" ", strings.ToLower(a.Content[i].Text))
		if matched {
			*IsFound = true

			err1 := psocket.SendToClient([]byte(fmt.Sprintf("%v", a.VideoId+" "+a.Content[i].Start)))
			if err1 != nil {
				return
			}
			break
		}
	}
}
