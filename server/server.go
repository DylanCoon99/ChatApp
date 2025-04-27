package main

import (
	"os"
	"bufio"
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{}




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


func sendEndpoint(c *gin.Context) {

	w,r := c.Writer, c.Request

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
        log.Println(err)
        return
    }

    defer conn.Close()


    go func() {

    	for {

    	// continuously read messages from connection
    	_, msg, err := conn.ReadMessage()

    	if err != nil {
    		log.Println("read:", err)
    		break
    	}
    	log.Printf("recv: %v\n>", string(msg))


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




func main() {

	fmt.Println("This is the server speaking.")

	r := gin.Default() // instantiate a router
	
	// endpoints
	r.GET("/send", sendEndpoint)


	log.Fatal(r.Run(":8080"))

}

