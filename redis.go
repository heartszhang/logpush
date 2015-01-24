package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/gosexy/redis"
)

func redis_pub(channel chan doc, conn net.Conn) {
	defer conn.Close()
	rediscli := redis.New()
	err := rediscli.Connect(option.redis_addr, option.redis_port)
	if err != nil {
		conn.Close()
		log.Println("redis-conn", err)
		return
	}
	defer rediscli.Close()

	for doc := range channel {
		log_document(doc)
		jbody, _ := json.Marshal(doc)
		if _, err = rediscli.Publish(doc["type"].(string), jbody); err != nil {
			log.Println("redis-publish", err)
			break
		}
	}
}
