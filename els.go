package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mattbaird/elastigo/lib"
)

//doc.tag is required
//doc.type is required
func els_post(channel chan doc, conn net.Conn) {
	var els = elastigo.NewConn()
	var hosts []string
	for i := 1; i < option.els_pool_num; i++ {
		hosts = append(hosts, fmt.Sprintf("127.0.0.%v:%v", i, option.els_port))
	}
	els.SetHosts(hosts)

	for doc := range channel {
		if dtype, ok := doc["type"].(string); ok {
			if _, err := els.Index(option.index_prefix, dtype, "", nil, doc); err != nil {
				log.Println("index", doc, err)
			}
		}
	}
}

func verbose(v ...interface{}) {
	if option.verbose {
		log.Println(v...)
	}
}
