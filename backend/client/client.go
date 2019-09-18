package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
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

	fmt.Print("Enter Latitude: ")
	lat, _ := reader.ReadString('\n')
	lat = strings.Replace(lat, "\n", "", -1)

	fmt.Print("Enter Longitude: ")
	lon, _ := reader.ReadString('\n')
	lon = strings.Replace(lon, "\n", "", -1)

	latFloat, _ := strconv.ParseFloat(lat, 64)
	lonFloat, _ := strconv.ParseFloat(lon, 64)

	connectMsg := fmt.Sprintf(`{"type": 0, "data": {"deviceId": "%s", "lat": %f, "lon":%f } }`, deviceID, latFloat, lonFloat)
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
		fmt.Print("Enter msg to send: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.Replace(msg, "\n", "", -1)
		json := fmt.Sprintf(`{"type": 1, "data": {"deviceId": "%s", "msg":"%s" } }`, deviceID, msg)
		fmt.Printf(`Sending message: %s`, json)

		err = c.WriteMessage(websocket.TextMessage, []byte(json))
		if err != nil {
			fmt.Println("write:", err)
			return
		}
	}
}
