package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	serverHost, serverPort, serverHTML := args()
	showHTML := ""
	if serverHTML != "./" {
		showHTML = "html:" + serverHTML
	}
	log.Println("ws video test ->", serverHost+":"+serverPort, showHTML)
	if serverHTML != "" {
		fs := http.FileServer(http.Dir(serverHTML))
		http.Handle("/", fs)
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
		log.Println("connection", r.RemoteAddr)
		go func(conn *websocket.Conn) {
			for {
				mt, data, connErr := conn.ReadMessage()
				if connErr != nil {
					log.Println("error", connErr)
					return
				}
				if mt == 1 {
					log.Println(data)
				}
				if mt == 2 {
					if err := conn.WriteMessage(2, data); err != nil {
						log.Println(err)
					}
				}
			}
		}(conn)
	})
	log.Fatal(http.ListenAndServe(serverHost+":"+serverPort, nil))
}
