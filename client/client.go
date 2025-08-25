package main

import (
	"bufio"
	"fmt"
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
	receive_channel := make(chan []byte)
	send_channel := make(chan []byte)
	input_channel := make(chan string)
	var wg_main sync.WaitGroup
	defer wg_main.Wait()
	conn, err := net.Dial(SERVER_TYPE, SERVER_PATH)
	if err != nil {
		fmt.Println("[error] - algo deu errado!")
		panic(err)
	}

	wg_main.Add(1)
	//	go handleConnection(conn, &wg_main, send_channel, receive_channel, input_channel)
	reading_state := "paused"
	wg_main.Add(1)
	go handleKeyboard(input_channel, &reading_state, &wg_main)
	wg_main.Add(1)
	go utils.ReceiveHandler(conn, receive_channel, &wg_main)
	wg_main.Add(1)
	go utils.SendHandler(conn, send_channel, &wg_main)
  
  tmp := false
	for !tmp {
    tmp = set_username(input_channel, send_channel, receive_channel, &reading_state) 
	}

	reading_state = "ready"
	in_room := false

main_loop:
	for {
		select {
		case data := <-input_channel:
			if data == ":q" {
				reading_state = "stopped"
				msg := utils.Message{Cmd: "quit", Data: USERNAME}
				serialized, err := utils.SerializeJson(msg)
				if err != nil {
					fmt.Println("[debug] - error while serializing:", err)
				} else if data != "\n" {
					send_channel <- serialized
				}
				wg_main.Done()
				wg_main.Done()
				wg_main.Done()
				break main_loop
			} else if in_room {
				msg := utils.Message{Cmd: "message", Data: data}
				serialized, err := utils.SerializeJson(msg)
				if err != nil {
					fmt.Println("[error] - error while serializing:", err)
				} else {
					send_channel <- serialized
				}
			}
		case received_data := <-receive_channel:
			var msg utils.Message
			err := utils.DeserializeToJson(received_data, &msg)
			if err != nil {
				fmt.Println("[error] - error while deserializing:", err)
			} else {
				switch msg.Cmd{
        case "message":
          fmt.Println(msg.Data)
        case "allocated":
          fmt.Println("You are now in a room with:", msg.Data)
          in_room = true
        default:
          fmt.Println("dont know what to do yet")
				}
			}
		}
	}
}

func set_username(input_channel chan string, send_channel chan []byte, receive_channel chan []byte, reading_state *string) bool {
	fmt.Println("username:")
	*reading_state = "ready"
	USERNAME = <-input_channel
	*reading_state = "paused"
	msg := utils.Message{Cmd: "set_user", Data: USERNAME}
	serialized, err := utils.SerializeJson(msg)
	if err != nil {
		fmt.Println("[error] - error while serializing\n", err)
	}
	send_channel <- serialized

	received := <-receive_channel
	var received_data utils.Message
	err = utils.DeserializeToJson(received, &received_data)
	if err != nil {
		fmt.Println("[error] - error while deserializing:", err)
		return false
	} else {
		switch received_data.Cmd {
		case "set_user":
			if received_data.Data == "ok" {
        fmt.Println("username saved")
				return true
			}
      fmt.Println("set user error")
			return false
		default:
      fmt.Println("default - data:", received_data.Data)

			return false
		}
	}
}


func handleKeyboard(keyboard_input chan string, reading_state *string, wg_main *sync.WaitGroup) {
	defer wg_main.Done()
	scanner := bufio.NewScanner(os.Stdin)
loop:
	for {
		switch *reading_state {
		case "paused":
			continue
		case "stopped":
			break loop
		default:
			scanner.Scan()
			input := scanner.Text()
			keyboard_input <- input
		}
	}
}
