package screws

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func TestWebsocketServer(t *testing.T) {

	var WSM = NewWSManager()

	r := gin.Default()
	r.GET("/websocket", func(c *gin.Context) {
		user := c.Query("user")
		if user == "" {
			c.JSON(http.StatusForbidden, nil)
			return
		}
		// change the request to websocket model
		conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
		if error != nil {
			http.NotFound(c.Writer, c.Request)
			return
		}
		// websocket connect
		wsClient := NewWSClient(user, conn, WSM)
		go wsClient.Reading()
		go wsClient.Writing()

		WSM.Broadcast([]byte(fmt.Sprintf("Broadcast: webcome %s.", user)))
		WSM.Notice([]byte(fmt.Sprintf("Hello %s, this is server.", user)), wsClient)
	})

	// websocket manager start
	go func() {
		WSM.Start()
	}()

	// go func() {
	// 	for {
	// 		time.Sleep(time.Second * 5)
	// 		fmt.Println(WSM.GetClients())
	// 		WSM.Broadcast([]byte(fmt.Sprintf("%v %d", time.Now(), len(WSM.GetClients()))))
	// 	}
	// }()

	r.Run("127.0.0.1:8040")
}

func TestWebsoketClient(t *testing.T) {
	user := uuid.NewV4().String()
	loadWebsocketClient(user)
}

func TestConcurrentWebsoketClient(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			loadWebsocketClient(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func loadWebsocketClient(user interface{}) {
	var dialer *websocket.Dialer
	conn, resp, err := dialer.Dial(fmt.Sprintf("ws://127.0.0.1:8040/websocket?user=%v", user), nil)
	if err != nil {
		log.Println(resp, err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hi, im client %v.", user)))
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		fmt.Printf("received: %s\n", message)
	}
}

/*HTML Test


<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
</head>

<body>

</body>

</html>
<script>
    var ws = new WebSocket("ws://127.0.0.1:8040/ws?user=1001");

    ws.onopen = (evt) => {
        console.log("Connection open ...");
        ws.send("hello server!")
    }

    ws.onmessage = (evt) => {
        console.log("Received Message:" + evt.data);
    }

    ws.onerror = (evt) => {
        console.log("Error: " + evt.data);
    }

    ws.onclose = (evt) => {
        console.log("Connection closed.");
        ws = null;
    }
</script>

*/
