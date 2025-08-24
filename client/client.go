package main

import (
	"fmt"
	"bufio"
	"net"
	"os"
	"sync"
	"vini464/simple-chat/utils"
)

const (
	SERVER_PATH = "server:7070"
	SERVER_TYPE = "tcp"
)

var USERNAME string = "you"

func main() {
  scanner := bufio.NewScanner(os.Stdin)
	conn, err := net.Dial(SERVER_TYPE, SERVER_PATH)
	if err != nil {
		fmt.Println("[error] - algo deu errado!")
		panic(err)
	}
	fmt.Println("Insert your username: ")
  scanner.Scan()
	USERNAME = scanner.Text()

	end := make(chan bool)
	go handleConnection(conn, end)
	<-end
}

func handleConnection(conn net.Conn, end_chan chan bool) {
	var wg_server sync.WaitGroup
	defer wg_server.Wait()
	received_data := make(chan []byte)
	data_to_send := make(chan []byte)
	keyboard_input := make(chan string)

	go handleKeyboard(keyboard_input)
	wg_server.Add(1)
	go utils.ReceiveHandler(conn, received_data, &wg_server)
	wg_server.Add(1)
	go utils.SendHandler(conn, data_to_send, &wg_server)

	can_send := true
	for {
		select {
		case data := <-received_data:
			fmt.Println(string(data))

		case data2 := <-keyboard_input:
			if can_send && len(data2) > 0 {
				data_to_send <- []byte(data2)
			}
		}
	}
	end_chan <- true
}

func handleKeyboard(keyboard_input chan string) {
		scanner := bufio.NewScanner(os.Stdin)

	for {

    scanner.Scan()
		input := scanner.Text()
		keyboard_input <- input
	}
}
