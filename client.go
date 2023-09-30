package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

// weiredest implementation of an async chat thingie
type Message struct {
	Name    string
	Message string
}

var peer string
var reader = bufio.NewReader(os.Stdin)
var port string
var name string

//var recv = make(chan string)

func main() {
	var wg sync.WaitGroup
	fmt.Print("Enter Port to host:")
	val, _ := reader.ReadString('\n')
	port = strings.Replace(val, "\n", "", -1)
	wg.Add(2)
	go ginInit(&wg)
	time.Sleep(time.Second * 2)
	go userMessage(&wg)

	wg.Wait()
}

func sendMessage(msg Message) {
	//send message to peer
	jsonValue, _ := json.Marshal(msg)
	_, err := http.Post("http://"+peer+"/recv", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatalln(err)
	}
}

func userMessage(wg *sync.WaitGroup) {
	fmt.Println("******Welcome*****")
	fmt.Print("Enter your name:")
	idk, _ := reader.ReadString('\n')
	name = strings.Replace(idk, "\n", "", -1)

	fmt.Print("Enter Peer Address eg(localhost:8080) :")
	val, _ := reader.ReadString('\n')
	peer = strings.Replace(val, "\n", "", -1)

	newmesg := Message{Name: name, Message: "Has Joined"}
	sendMessage(newmesg)

	exit := 0
	for i := 0; exit != 1; i++ {
		fmt.Print("\n" + "(" + name + "):")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if text == "exit" {
			exit = 1
			break
		}
		if text != "" {
			newmesg.Message = text
			sendMessage(newmesg)
		}

	}

	wg.Done()

}

func ginInit(wg *sync.WaitGroup) {

	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.POST("/recv", func(c *gin.Context) {
		//print json value
		bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
		jsonBody := string(bodyAsByteArray)
		herename := gjson.Get(jsonBody, "Name")
		message := gjson.Get(jsonBody, "Message")
		fmt.Println("\n" + "(" + herename.String() + "):" + message.String())
		fmt.Print("(" + name + "):")
		c.JSON(http.StatusOK, "ok")
	})
	r.Run(":" + port)
	wg.Done()
}
