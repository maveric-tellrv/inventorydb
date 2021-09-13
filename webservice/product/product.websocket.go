package product

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

func productSocket(ws *websocket.Conn) {
	log.Println("Registering WebSocket Connection Established...")
	done := make(chan struct{})

	// Function to listen incocming data on webscoket
	// because we dont want to block we use go routine
	go func(c *websocket.Conn) {
		log.Println("Outside For loop:....")
		for {
			var msg message
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				log.Println(err)
				break
			}
			log.Println("INside For loop:....")
			fmt.Printf("Received Message %s\n", msg.Data)
		}
		close(done)
	}(ws)
loop:

	for {
		select {
		case <-done:
			fmt.Println("Connection closed: ...")
			break loop

		default:
			products, err := GetTopTenProducts()
			if err != nil {
				log.Println(err)
				break
			}
			if err := websocket.JSON.Send(ws, products); err != nil {
				log.Println(err)
				break
			}
			time.Sleep(10 * time.Second)

		}

	}
	fmt.Println("Closing websocket connection....")
	defer ws.Close()

}
