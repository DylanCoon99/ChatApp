package main

import (
	"os"
	"log"
	"time"
	"os/signal"
	//"text/template"
	//"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)


func main() {

	log.Println("This is the client speaking.")


	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// conn, resp, err
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/echo", nil)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	done := make(chan struct{})

	go func() {

		defer close(done)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()


	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {

		select {
		case <-done:
			// if done can be read from --> return
			return
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}

	}

}


