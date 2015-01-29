package main

import (
	"fmt"
	"log"
	"net"
	"time"

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
		dtype, ok := doc["type"].(string)
		if !ok {
			dtype = option.dft_type
			verbose(doc)
		}
		if _, ok := doc["time"].(time.Time); !ok {
			doc["time"] = time.Now()
		}
		if _, err := els.Index(option.index_prefix, dtype, "", nil, doc); err != nil {
			log.Println("index", doc, err)
		}
	}
}

func verbose(v ...interface{}) {
	if option.verbose {
		log.Println(v...)
	}
}
