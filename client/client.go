package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
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
	go handleReceive(conn, received_data)
	go handleSend(conn, data_to_send)

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

func handleReceive(conn net.Conn, received_data chan string) {
	header := make([]byte, 4)
	msg := ""
	for {
		_, err := conn.Read(header)
		if err != nil {
			fmt.Println("[error] - algo deu errado")
			return
		}
		size := binary.BigEndian.Uint32(header)
		data := make([]byte, size)
		readed := 0
		for strLen, err := conn.Read(data); err != nil || readed < int(size); {
			if err != nil {
				fmt.Println("[error] - connection lost")
				return
			}
			msg += string(data)
			readed += strLen
			strLen = 0
		}
		received_data <- msg
	}
}

func handleSend(conn net.Conn, msg chan string) {
	for {
		data := <-msg
		size := uint32(len(data))
		header := make([]byte, 4)

		binary.BigEndian.PutUint32(header, size)
		_, err := conn.Write(header)
		if err != nil {
			fmt.Println("[error] - connection lost")
			panic(err)
		}
		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Println("[error] - connection lost")
			panic(err)
		}
	}
}
