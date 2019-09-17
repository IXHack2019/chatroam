package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	log.SetFlags(log.LstdFlags)

	u := url.URL{Scheme: "ws", Host: "localhost:8888", Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Device ID: ")
	deviceID, _ := reader.ReadString('\n')
	deviceID = strings.Replace(deviceID, "\n", "", -1)

	connectMsg := fmt.Sprintf(`{"type": 0, "data": {"deviceID": "%s"} }`, deviceID)
	log.Printf("Sending message: %s\n", connectMsg)

	err = c.WriteMessage(websocket.TextMessage, []byte(connectMsg))
	if err != nil {
		log.Println("write:", err)
		return
	}

	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Printf("Connect response: %s\n", message)

	for {
		fmt.Print("Enter JSON to send: ")
		json, _ := reader.ReadString('\n')
		// convert CRLF to LF
		json = strings.Replace(json, "\n", "", -1)

		err = c.WriteMessage(websocket.TextMessage, []byte(json))
		if err != nil {
			log.Println("write:", err)
			return
		}

		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("Received response: %s\n", message)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				if strings.Contains(string(message), `"type":1`) {
					log.Printf("Received broadcast: %s", message)
				} else {
					err := c.WriteMessage(websocket.TextMessage, message)
					if err != nil {
						log.Println("write:", err)
						return
					}
				}

			}
		}
	}()
}
