package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// CamClient do your shit
type CamClient struct {
	conn     *websocket.Conn
	remoteID string
}

// MicClient do your shit
type MicClient struct {
	conn     *websocket.Conn
	remoteID string
}

var camClients = make(map[string]CamClient)

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
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("error", err)
			return
		}
		ID := RandomString(3)
		camClients[ID] = CamClient{conn, ID}
		mu.Lock()
		if err := conn.WriteMessage(1, []byte(ID)); err != nil {
			log.Println(err)
			conn.Close()
			delete(camClients, ID)
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
					camClients[ID] = CamClient{conn, string(data)}
				}
				if mt == 2 {
					for _, client := range camClients {
						if client.conn == conn {
							if _, ok := camClients[client.remoteID]; ok {
								mu.Lock()
								if err := camClients[client.remoteID].conn.WriteMessage(2, data); err != nil {
									log.Println(err)
									if err := camClients[client.remoteID].conn.Close(); err != nil {
										log.Println(err)
									}
									delete(camClients, client.remoteID)
								}
								mu.Unlock()
							}
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
		micClients[ID] = MicClient{conn, ID}
		mu.Lock()
		if err := conn.WriteMessage(1, []byte(ID)); err != nil {
			log.Println(err)
			conn.Close()
			delete(micClients, ID)
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
					micClients[ID] = MicClient{conn, string(data)}
				}
				if mt == 2 {
					for _, client := range micClients {
						mu.Lock()
						if client.conn == conn {
							if _, ok := micClients[client.remoteID]; ok {
								if err := micClients[client.remoteID].conn.WriteMessage(2, data); err != nil {
									log.Println(err)
									if err := micClients[client.remoteID].conn.Close(); err != nil {
										log.Println(err)
									}
									delete(micClients, client.remoteID)
								}
							}
						}
						mu.Unlock()
					}
				}
			}
		}(conn)
	})

	log.Fatal(http.ListenAndServe(serverHost+":"+serverPort, nil))
}
