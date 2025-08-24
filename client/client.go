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
	scanner := bufio.NewScanner(os.Stdin)
	conn, err := net.Dial(SERVER_TYPE, SERVER_PATH)
	if err != nil {
		fmt.Println("[error] - algo deu errado!")
		panic(err)
	}
	fmt.Println("Insert your username: ")
	scanner.Scan()
	USERNAME = scanner.Text()

	wg_main.Add(1)
	go handleConnection(conn, &wg_main, send_channel, receive_channel, input_channel)
	reading_state := "paused"
	wg_main.Add(1)
	go handleKeyboard(input_channel, &reading_state, &wg_main)
}

func handleConnection(conn net.Conn, wg_main *sync.WaitGroup, send_channel chan []byte, receive_channel chan []byte, input_channel chan string) {
	var wg_server sync.WaitGroup
	defer wg_server.Wait()
	defer wg_main.Done()

	wg_server.Add(1)
	go utils.ReceiveHandler(conn, receive_channel, &wg_server)
	wg_server.Add(1)
	go utils.SendHandler(conn, send_channel, &wg_server)

	can_send := true
	for {
		select {
		case data := <-receive_channel:
      var received_data utils.Message
      err := utils.DeserializeToJson(data, &received_data)
      if (err != nil) {
        fmt.Println("[error] - error while deserializing:", err)
      } else {
        switch received_data.Cmd {
        case "message":
          fmt.Println(received_data.Data)
        default:
          fmt.Println("[error] - dont know what to do")
        }
      }

		case data2 := <-input_channel:
			if can_send && len(data2) > 0 {
				msg := utils.Message{Cmd: "message", Data: data2}
				serialized, err := utils.SerializeJson(msg)
				if err != nil {
					fmt.Println("[error] - error while serializing\n", err)
				}
				send_channel <- serialized
			}
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
