package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/mattbaird/elastigo/lib"
)

const sample = `/firstGame/Android/anzhi002/0.1.0.60_0.0.0.0.2/210/3ff4e3f1c3a71cd99bceeb891577b2fc2/anzhi_201501171417022I8S29P7dT/17/21537/MissionCompleted/12_3`

//const sample = `/firstGame/Android/funs0012/0.1.0.41/serverUnknown/24b983d3eb4fe1dfcd0b47ae6f8a6145/userUnknown/1/Start`
//const sample = `/firstGame/Android/funs0012/0.1.0.41/serverUnknown/22f9654c5ed59f08ffd91118b92cf/userUnknown/1/Start`

func TestDecode_func(t *testing.T) {
	//	t.Skip()
	v := decode_func(sample)
	t.Log(v)
}

func TestElsPost(t *testing.T) {
	t.Skip()
	channel := make(chan doc, 12)
	go els_postx(channel, t)
	doc := decode_func(sample)
	channel <- doc
	time.Sleep(time.Second * 4)
}

func els_postx(channel chan doc, t *testing.T) {
	var els = elastigo.NewConn()
	var hosts []string
	for i := 1; i < option.els_pool_num; i++ {
		hosts = append(hosts, fmt.Sprintf("127.0.0.%v:%v", i, option.els_port))
	}
	els.SetHosts(hosts)

	for doc := range channel {
		if dt, ok := doc["dtype"].(string); ok {
			if resp, err := els.Index(option.index_prefix, dt, "", nil, doc); err != nil {
				t.Log("index", doc, err)
			} else {
				x, _ := json.Marshal(resp)
				t.Log("resp.id", string(x))
			}
		}
	}
}
