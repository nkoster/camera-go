package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Client do a shit
type Client struct {
	conn     *websocket.Conn
	localID  string
	remoteID string
}

// MicClient do a shit
type MicClient struct {
	conn     *websocket.Conn
	localID  string
	remoteID string
}

// type MicClient struct {
// 	active bool
// }

var clients = make(map[string]Client)

// MicClients do a shit
// var MicClients = make(map[*websocket.Conn]MicClient)
var micClients = make(map[string]MicClient)

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
	http.HandleFunc("/cam", func(w http.ResponseWriter, r *http.Request) {
		// not safe, only for dev:
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("error", err)
			return
		}
		ID := RandomString(3)
		clients[ID] = Client{conn, ID, ""}
		mu.Lock()
		if err := conn.WriteMessage(1, []byte(ID)); err != nil {
			log.Println(err)
			conn.Close()
			delete(clients, ID)
		}
		mu.Unlock()
		log.Println("cam connection", r.RemoteAddr, ID)
		go func(conn *websocket.Conn) {
			for {
				mt, data, connErr := conn.ReadMessage()
				if connErr != nil {
					log.Println("error", connErr)
					return
				}
				if mt == 1 {
					clients[ID] = Client{conn: conn, localID: ID, remoteID: string(data)}
				}
				if mt == 2 {
					for id, client := range clients {
						if id == clients[ID].remoteID && client.conn != nil {
							mu.Lock()
							if err := client.conn.WriteMessage(2, data); err != nil {
								log.Println(err)
								if err := client.conn.Close(); err != nil {
									log.Println(err)
								}
								delete(clients, id)
							}
							mu.Unlock()
						}
					}
				}
			}
		}(conn)
	})

	http.HandleFunc("/mic", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("error", err)
			return
		}
		ID := RandomString(3)
		micClients[ID] = MicClient{conn, ID, ""}
		mu.Lock()
		if err := conn.WriteMessage(1, []byte(ID)); err != nil {
			log.Println(err)
			conn.Close()
			delete(clients, ID)
		}
		mu.Unlock()
		log.Println("mic connection", r.RemoteAddr, ID)

		go func(conn *websocket.Conn) {
			for {
				mt, data, connErr := conn.ReadMessage()
				if connErr != nil {
					log.Println("error", connErr)
					return
				}
				if mt == 1 {
					micClients[ID] = MicClient{conn: conn, localID: ID, remoteID: string(data)}
				}
				if mt == 2 {
					for id, client := range micClients {
						if id == micClients[ID].remoteID && client.conn != nil {
							mu.Lock()
							if err := client.conn.WriteMessage(2, data); err != nil {
								log.Println(err)
								if err := client.conn.Close(); err != nil {
									log.Println(err)
								}
								delete(micClients, id)
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
