package main

import (
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


func echo(c *gin.Context) {
	
	
	w,r := c.Writer, c.Request

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
        log.Println(err)
        return
    }

    defer conn.Close()

    for {
    	// continuously read from connection
    	msgType, message, err := conn.ReadMessage()
    	if err != nil {
    		log.Println("read:", err)
    		break
    	}
    	log.Printf("recv:%s", message)

    	// echo the message back to the client
    	err = conn.WriteMessage(msgType, message)
    	if err != nil {
    		log.Println("write:", err)
    		break
    	}
    }


}


func main() {

	fmt.Println("This is the server speaking.")

	r := gin.Default() // instantiate a router
	
	// endpoints
	r.GET("/", home)
	r.GET("/echo", echo)


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