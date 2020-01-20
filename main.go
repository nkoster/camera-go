package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var clients = make(map[string]*websocket.Conn)

var mu sync.Mutex

func main() {
	serverHost, serverPort, serverHTML := args()
	showHTML := ""
	if serverHTML != "./" {
		showHTML = "html:" + serverHTML
	}
	log.Println("ws video test ->", serverHost+":"+serverPort, showHTML)
	if serverHTML != "" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, serverHTML)
		})
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// not safe, only for dev:
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("error", err)
			return
		}
		ID := RandomString(10)
		clients[ID] = conn
		log.Println("connection", r.RemoteAddr, ID)
		go func(conn *websocket.Conn) {
			for {
				mt, data, connErr := conn.ReadMessage()
				if connErr != nil {
					log.Println("error", connErr)
					return
				}
				if mt == 1 {
					log.Println("type 1")
				}
				if mt == 2 {
					for id, client := range clients {
						if id != ID {
							mu.Lock()
							if err := client.WriteMessage(2, data); err != nil {
								log.Println(err)
								client.Close()
								delete(clients, id)
							}
							mu.Unlock()
						}
					}
				}
			}
		}(conn)
	})
	log.Fatal(http.ListenAndServe(serverHost+":"+serverPort, nil))
}
