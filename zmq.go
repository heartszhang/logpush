package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/pebbe/zmq4"
)

func zmq_push(channel chan doc, conn net.Conn) {
	defer conn.Close()
	pusher, err := zmq4.NewSocket(zmq4.PUSH)
	if err == nil {
		defer pusher.Close()
	}
	err = pusher.Bind(option.zmq_addr)
	for doc := range channel {
		if data, err := json.Marshal(doc); err == nil {
			if _, err = pusher.SendBytes(data, 0); err != nil {
				log.Println("zmq-push", err)
				break
			}
		}
	}
}
