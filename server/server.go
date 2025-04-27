package main

import (
	"os"
	"bufio"
	"fmt"
	"log"
	"text/template"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{}


func home(c *gin.Context) {
	homeTemplate.Execute(c.Writer, "ws://"+c.Request.Host+"/echo")
}






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


    for {

    	// continuously read messages from connection
    	_, message, err := conn.ReadMessage()

    	if err != nil {
    		log.Println("read:", err)
    		break
    	}
    	log.Printf("recv: %v", string(message))

    }


    // infinite loop for user input
	reader := bufio.NewReader(os.Stdin)
	msgChan := make(chan string)

	done := make(chan struct{}) //unbuffered channel; senders block until receivers can receive


	for {

		fmt.Print(">")

		msg, err := reader.ReadString(byte('\n'))

		if err != nil {
			log.Print(err)
			return
		}

		msgChan <- msg
		
	}


	go sendMessage(conn, msgChan, done)


}




func main() {

	fmt.Println("This is the server speaking.")

	r := gin.Default() // instantiate a router
	
	// endpoints
	r.GET("/", home)
	r.GET("/send", sendEndpoint)


	log.Fatal(r.Run(":8080"))

}


var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))