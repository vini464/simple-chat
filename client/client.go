package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
  "github.com/vini464/simple-chat/communication"
)

const (
	SERVER_PATH = "localhost:7070"
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
	fmt.Print("> ")
	scanner.Scan()
	USERNAME = scanner.Text()
	end := make(chan bool)
	go handleConnection(conn, end)
	<-end
}

func handleConnection(conn net.Conn, end_chan chan bool) {
	received_data := make(chan string)
	data_to_send := make(chan string)
	keyboard_input := make(chan string)

	go handleKeyboard(keyboard_input)
	// go handleReceive(conn, received_data)
  go communication.handleReceive(conn, received_data)
	// go handleSend(conn, data_to_send)

	can_send := true
	for {
		select {
		case data := <-received_data:
			fmt.Println(data)

		case data2 := <-keyboard_input:
			if can_send {
				data_to_send <- data2

			}
		}
	}
	end_chan <- true
}

func handleKeyboard(keyboard_input chan string) {
	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print(USERNAME + ": ")
		scanner.Scan()
		input := scanner.Text()
		fmt.Println(input)
		keyboard_input <- input
	}
}

