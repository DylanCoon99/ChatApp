package main

import (
	"os"
	"log"
	"fmt"
	//"time"
	"bufio"
	"os/signal"
	//"text/template"
	//"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)


func sendMessage(conn *websocket.Conn, msgChan chan string, done chan struct{}) {

	go func() {

		// sends messages at it receives from the channel
		for {

			select {
			case <-done:
				return
			case msg := <-msgChan:
				// send the message
				err := conn.WriteMessage(websocket.TextMessage, []byte(msg))

				if err != nil {
					log.Printf("failed to send msg: %v", err)
					return
				}

			}

		}

	}()

}



func main() {

	log.Println("This is the client speaking.")


	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// conn, resp, err
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/send", nil)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()


	go func() {

		//defer close(done)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s\n>", message)
		}
	}()



	// infinite loop for user input
	reader := bufio.NewReader(os.Stdin)
	msgChan := make(chan string)

	done := make(chan struct{}) //unbuffered channel; senders block until receivers can receive

	go sendMessage(conn, msgChan, done)

	for {

		fmt.Print(">")

		msg, err := reader.ReadString(byte('\n'))

		if err != nil {
			log.Print(err)
			return
		}

		msgChan <- msg
		
	}

}


