package main

import (
	"fmt"
	"net"
	"sync"
	"vini464/simple-chat/utils"
)

var CLIENTS = make(map[string]chan utils.Message)
var QUEUE = make([]string, 0)
var ROOMS = make(map[string]string)

const (
	SERVER_HOST = "server"
	SERVER_PORT = "7070"
	SERVER_TYPE = "tcp"
	SERVER_PATH = SERVER_HOST + ":" + SERVER_PORT
)

func main() {

	fmt.Println("[log] Iniciando servidor...")
	listener, err := net.Listen(SERVER_TYPE, SERVER_PATH)
	if err != nil {
		fmt.Println("[error] Ocorreu um erro:", err.Error())
	}
	fmt.Println("[log] Escutando em:", listener.Addr())
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[error] Não foi possível aceitar a conexão:", err.Error())
			continue
		}
		go handleClient(conn)
	}

}

func handleClient(conn net.Conn) {
	var username string
	var wg_clients sync.WaitGroup
	receive_channel := make(chan []byte)
	data_to_send := make(chan []byte)
	self_channel := make(chan utils.Message)

	defer conn.Close()
	defer wg_clients.Wait()
	defer delete(CLIENTS, username)

	fmt.Println("[log] Novo cliente:", conn.RemoteAddr().String())

	wg_clients.Add(1)
	go utils.ReceiveHandler(conn, receive_channel, &wg_clients)
	wg_clients.Add(1)
	go utils.SendHandler(conn, data_to_send, &wg_clients)

	for {
		select {
		case income := <-receive_channel:
      var message utils.Message
      utils.DeserializeToJson(income, &message)
      switch message.Cmd {
      case "quit":
        fmt.Println("[debug] - case quit")
        delete(CLIENTS, username)
        tmp, ok := ROOMS[username]
        if ok {
          delete(ROOMS, username)
          delete(ROOMS, tmp)
        }
      case "message":
        fmt.Println("[debug] - case message")
        _, ok := ROOMS[username]
        if ok {
          msg := utils.Message{Cmd: "message", Data: "[" + username + "]:" + message.Data}
          CLIENTS[ROOMS[username]] <- msg
        }
      case "set_user":
        fmt.Println("[debug] - case set_user")
        var response utils.Message
        _, ok := CLIENTS[message.Data]
        if ok {
          fmt.Println("[debug] - ok")
          response = utils.Message{Cmd: "error", Data: "invalid user"}
          sendResponse(response, data_to_send)
        } else {
          fmt.Println("[debug] - else")
          response = utils.Message{Cmd: "set_user", Data: "ok"}
          CLIENTS[message.Data] = self_channel
          username = message.Data
          sendResponse(response, data_to_send)
          if len(QUEUE) == 0 {
            QUEUE = utils.Enqueue(QUEUE, username)
            fmt.Println("[debug] - in queue. queue size:", len(QUEUE))
          } else {
            fmt.Println("[debug] -  dequeue")
            var other string
            other, QUEUE = utils.Dequeue(QUEUE)
            ROOMS[username] = other
            ROOMS[other] = username
            fmt.Println("[debug] -  before response")
            response := utils.Message{Cmd: "allocated", Data: username}
            fmt.Println("[debug] - other:", other)
            CLIENTS[other] <- response
            fmt.Println("[debug] -  after")
            response = utils.Message{Cmd: "allocated", Data: other}
            sendResponse(response, data_to_send)
            fmt.Println("[debug] - in room, queue size:", len(QUEUE))
          }
        }
      default:
        fmt.Println("[debug]: unknow command")
      }
		case trasmition := <-CLIENTS[username]:
      sendResponse(trasmition, data_to_send)
		}

  }
}

func sendResponse(msg utils.Message, send_channel chan []byte) {
	serialized, err := utils.SerializeJson(msg)
	for err != nil {
		fmt.Println("[error] - error while serializing:", err)
		serialized, err = utils.SerializeJson(msg)
	}
	send_channel <- serialized
}
