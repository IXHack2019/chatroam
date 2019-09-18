package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8888", Path: "/connect"}
	fmt.Printf("connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Printf("dial: %s", err)
	}
	defer c.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Device ID: ")
	deviceID, _ := reader.ReadString('\n')
	deviceID = strings.Replace(deviceID, "\n", "", -1)

	connectMsg := fmt.Sprintf(`{"type": 0, "data": {"deviceId": "%s", "lon": -79.380688, "lat":43.652112 } }`, deviceID)
	fmt.Printf("Sending message: %s\n", connectMsg)

	err = c.WriteMessage(websocket.TextMessage, []byte(connectMsg))
	if err != nil {
		fmt.Println("write:", err)
		return
	}

	_, message, err := c.ReadMessage()
	if err != nil {
		fmt.Println("read:", err)
		return
	}
	fmt.Printf("Connect response: %s\n", message)

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				return
			}

			fmt.Printf("\nReceived message: %s\nEnter JSON to send: ", message)
		}
	}()

	for {
		fmt.Print("Enter JSON to send: ")
		json, _ := reader.ReadString('\n')
		json = strings.Replace(json, "\n", "", -1)
		fmt.Printf("Sending message: %s\n", json)

		err = c.WriteMessage(websocket.TextMessage, []byte(json))
		if err != nil {
			fmt.Println("write:", err)
			return
		}
	}
}
