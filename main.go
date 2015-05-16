package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	_ "fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var messageDB []Message
var messages chan Message
var sockets []*websocket.Conn
var socketsLock sync.Mutex

// HandleSocket handles a new websocket
func HandleSocket(ws *websocket.Conn) {
	var err error

	sockets = append(sockets, ws)
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Send previous message
	socketsLock.Lock()
	for _, msg := range messageDB {

		json, err := json.Marshal(msg)
		if err != nil {
			log.Println("Couldn't unmarshal JSON: ", err)
		}

		if err = websocket.Message.Send(ws, string(json)); err != nil {
			log.Println("Can't send: ", err.Error())
			break
		}
	}
	socketsLock.Unlock()

	go func() {
		// Wait for incoming websocket messages
		for {
			var reply string

			if err = websocket.Message.Receive(ws, &reply); err != nil {
				log.Println("Can't receive: ", err.Error())
				break
			}

			log.Println("Received back from client: " + reply)

			messages <- Message{Content: reply, Time: time.Now().Unix()}

		}

		// Remove closed socket from socket list
		for index, closeSocket := range sockets {
			if closeSocket == ws {
				sockets = append(sockets[:index], sockets[index+1:]...)
			}
		}

		wg.Done()

	}()

	wg.Wait()

}

func WaitAndBroadcast() {
	// Wait for new message to broadcast
	for msg := range messages {

		// Save message to disk
		messageDB = append(messageDB, msg)
		err := EncodeFile("./messages.json", &messageDB)
		if err != nil {
			log.Println(err)
		}

		// Marshal JSON
		json, err := json.Marshal(msg)
		if err != nil {
			log.Println("Couldn't unmarshal JSON: ", err)
		}

		socketsLock.Lock()
		for _, ws := range sockets {

			if err = websocket.Message.Send(ws, string(json)); err != nil {
				log.Println("Can't send: ", err.Error())
				break
			}
		}
		socketsLock.Unlock()
	}
}

func main() {

	// Init
	messages = make(chan Message)
	messageDB = make([]Message, 0)
	sockets = make([]*websocket.Conn, 0)
	socketsLock = sync.Mutex{}

	err := DecodeFile("./messages.json", &messageDB)
	if err != nil {
		log.Println(err)
	}

	// Wait for messages to broadcast to all sockets
	go WaitAndBroadcast()

	// Handle socket
	http.HandleFunc("/",
		func(w http.ResponseWriter, req *http.Request) {
			s := websocket.Server{Handler: websocket.Handler(HandleSocket)}
			s.ServeHTTP(w, req)
		})

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
