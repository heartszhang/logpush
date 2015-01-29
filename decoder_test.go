package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mattbaird/elastigo/lib"
)

//const sample = `/firstGame/Android/anzhi002/0.1.0.60_0.0.0.0.2/210/3ff4e3f1c3a71cd99bceeb891577b2fc2/anzhi_201501171417022I8S29P7dT/17/21537/MissionCompleted/12_3`

//const sample = `/firstGame/Android/funs0012/0.1.0.41/serverUnknown/24b983d3eb4fe1dfcd0b47ae6f8a6145/userUnknown/1/Start`
//const sample = `/firstGame/Android/funs0012/0.1.0.41/serverUnknown/22f9654c5ed59f08ffd91118b92cf/userUnknown/1/Start`
const sample = `C8287130/B3E2BC42 [REL 2015-01-27 18:35:38] 玩家[23101],通过[装备升级],[消耗][银两],数量为[2480],获得[0],数量为[1], 玩家银两数为[172275-4]`

func TestDecode_func(t *testing.T) {
	t.Skip()
	v := decode_func(sample)
	t.Log(v)
}

func TestElsPost(t *testing.T) {
	t.Skip()
	channel := make(chan doc, 12)
	go els_post(channel, nil)
	v := decode_document(doc{
		"timestamp": "p.header.timestamp",
		"hostname":  "p.header.hostname",
		"tag":       "p.message.tag",
		"priority":  "p.priority.P",
		"facility":  "p.priority.F.Value",
		"severity":  "p.priority.S.Value-x",
		"content":   sample,
	})
	if len(v) > 0 {
		channel <- v
	}
	time.Sleep(time.Second * 10)
}

func els_postx(channel chan doc, t *testing.T) {
	var els = elastigo.NewConn()

	for doc := range channel {
		dt, ok := doc["type"].(string)
		if !ok {
			dt = "generic"
		}
		if resp, err := els.Index(option.index_prefix, dt, "", nil, doc); err != nil {
			t.Log("index", doc, err)
		} else {
			x, _ := json.Marshal(resp)
			t.Log("resp.id", string(x))
		}
	}
}
