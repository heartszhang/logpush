package main

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

func init() {
	WordDecoder([]string{"msg-sent"}, db_msg_sent)
	WordDecoder([]string{"msg-rcv"}, db_msg_recv)
	WordDecoder([]string{"award-msg-num"}, db_msg_award)
	WordDecoder([]string{"_tbl"}, db_msg_tableupdate)
	//msg_rcv msg_id received_count max_handle_time_us
	WordDecoder([]string{"msg_rcv msg_id received_count"}, db_msg_ignore)
	WordDecoder([]string{"msg_snt msg_id msg_sent_count"}, db_msg_ignore)
	WordDecoder([]string{"db-status"}, db_msg_ignore)
	WordDecoder([]string{"db-status"}, db_msg_ignore)
	WordDecoder([]string{"game-status"}, db_msg_ignore)
	WordDecoder([]string{"msg-snt"}, db_msg_ignore)
	WordDecoder([]string{"S2CPipe-info-Used-size"}, db_msg_pipe)
	WordDecoder([]string{"C2SPipe-info-Used-size"}, db_msg_pipe)
	WordDecoder([]string{"AuthTCPConn-SendBuf-Used-size"}, db_msg_pipe)
	WordDecoder([]string{"AuthTCPConn-RecvBuf-Used-size"}, db_msg_pipe)
	WordDecoder([]string{"LogTCPConn-SendBuf-Used-size"}, db_msg_pipe)
	WordDecoder([]string{"LogTCPConn-RecvBuf-Used-size"}, db_msg_pipe)
	WordDecoder([]string{"DBTCPConn-RecvBuf-Used-size"}, db_msg_pipe)
	WordDecoder([]string{"DBTCPConn-SendBuf-Used-size"}, db_msg_pipe)
	WordDecoder([]string{"---------------------------"}, db_msg_ignore)
	WordDecoder([]string{"total socket"}, db_msg_socket)
	WordDecoder([]string{"LOGIN", "data", "post url"}, db_msg_login)
}

func db_msg_ignore(line string) doc {
	return make(doc)
}

func db_msg_generic(line, typ string, cnt int) doc {
	fields := strings.Fields(line)
	if len(fields) < cnt {
		return nil
	}
	v := doc{
		"path": fields[0],
		"time": firstgame_time(fields[1]),
		"type": typ,
	}
	for i := 2; i < cnt; i += 2 {
		v[strings.ToLower(fields[i])], _ = iconvert(fields[i+1])
	}
	return v
}

//82D47F4C/9EEA46A5 currurnt total socket num = 100,free socket num=97
func db_msg_socket(line string) (v doc) {
	fields := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' ' || r == '=' || r == ','
	})
	if len(fields) < 10 {
		return
	}
	return doc{
		"path": fields[0],
		"type": "socket",
		"cnt":  iconvert2(fields[5]),
		"free": iconvert2(fields[9]),
	}
}

//82D47F4C/9EEA46A5 ------------------------------------------------------------------

//`33EDAE8/8027C03D [ERR 2015-01-27 16:00:38] LOGIN :{"id":1422345637,"state":{"code":1,"msg":"操作成功"},"data":{"ucid":789994237,"nickName":"九游玩家789994237"}}, post url[http://sdk.g.uc.cn/ss/] data[789994237]`
//(path0) [... (time1)] LOGIN :(login2), post url[(url3)] data[(data4)]
//path0: [a-zA-Z0-9]+/[a-zA-Z0-9]+
//time1: \d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}
//login2: \{.*\}
//url3: https?://[^\]]+
//data4: \d+
var login_re = regexp.MustCompile(`([a-zA-Z0-9]+/[a-zA-Z0-9]+) \[... (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\] LOGIN :(\{.*\}), post url\[(https?://[^\]]+)\] data\[(\d+)\]`)

func db_msg_login(line string) (v doc) {
	fields := login_re.FindStringSubmatch(line)
	if len(fields) < 6 {
		return
	}
	v = doc{
		"path":  fields[1],
		"time":  firstgame_time(fields[2]),
		"url":   fields[4],
		"data":  iconvert2(fields[5]),
		"login": &doc{},
	}

	dec := json.NewDecoder(strings.NewReader(fields[3]))
	dec.UseNumber()
	// ignore the error
	if err := dec.Decode(v["login"]); err != nil {
		v["login2"] = fields[3]
	}

	return
}

//33EDAE8/8027C03D [ERR 2015-01-27 16:00:38] LOGIN :{"id":1422345637,"state":{"code":1,"msg":"操作成功"},"data":{"ucid":789994237,"nickName":"九游玩家789994237"}}, post url[http://sdk.g.uc.cn/ss/] data[789994237]
//(path0) [... (time1)] LOGIN :(login2), post url[(url3)] data[(data4)]
//path0: [a-zA-Z0-9]+/[a-zA-Z0-9]+
//time1: \d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}
//login2: \{.*\}
//url3: https?://[^\]]+
//data4: \d+
const login_rexp = `([a-zA-Z0-9]+/[a-zA-Z0-9]+) [... (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})] LOGIN :(\{.*\}), post url[(https?://[^\]]+)] data[(\d+)]`

var login_re = regexp.MustCompile(`([a-zA-Z0-9]+/[a-zA-Z0-9]+) [... (\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})] LOGIN :(\{.*\}), post url[(https?://[^\]]+)] data[(\d+)]`)

func db_msg_login(line string) (v doc) {
	fields := login_re.FindStringSubmatch(line)
	fields := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' ' || r == '=' || r == ','
	})
	if len(fields) < 10 {
		return
	}
	return doc{
		"path": fields[0],
		"type": "socket",
		"cnt":  iconvert2(fields[5]),
		"free", iconvert2(fields[9]),
	}
}

//24CCE4E8/EAA5104B 2015-01-26-180114 award-msg-num 123 click-msg-num 3 tips-msg-num 0 chat-msg-num 0

func db_msg_award(line string) (v doc) {
	return db_msg_generic(line, "message", 10)
}

//DBTCPConn-RecvBuf-Used-size 0 Surplus-size 10485760 MaxUsed-size 711724
//AuthTCPConn-SendBuf-Used-size 0 Surplus-size 10485760 MaxUsed-size 401
//LogTCPConn-SendBuf-Used-size 0 Surplus-size 10485760 MaxUsed-size 14230
//24CCE4E8/EAA5104B 2015-01-27-103904 S2CPipe-info-Used-size 0 Surplus-size 805306368 MaxUsed-size 26046
//24CCE4E8/EAA5104B 2015-01-26-180114 C2SPipe-info-Used-size 0 Surplus-size 805306368 MaxUsed-size 408
func db_msg_pipe(line string) (v doc) {
	return db_msg_generic(line, "pipe", 8)
}

func db_msg_0(line string) (v doc) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return nil
	}
	v = doc{
		"path": fields[0],
		"time": firstgame_time(fields[1]),
		"type": fields[2],
	}
	return v
}

//DE504002/BFBED4A22015 01-26-180034 Treasure_tbl UPDATE 71 0 11 0
func db_msg_tableupdate(line string) doc {
	return db_msg_x(line, 8)
}

func db_msg_x(line string, cnt int) doc {
	fields := strings.Fields(line)
	if len(fields) < cnt {
		return nil
	}
	v := doc{
		"path":  fields[0],
		"time":  firstgame_time(fields[1]),
		"type":  fields[2],
		"mtype": fields[3],
	}
	for i := 4; i < cnt; i++ {
		v[string('a'+i)] = fields[i]
	}
	return v
}

//24CCE4E8/EAA5104B2015 01-26-180224 msg-rcv id-154 1 30 30 30 30
func db_msg_recv(line string) doc {
	return db_msg_x(line, 9)
}

//24CCE4E8/EAA5104B 2015-01-26-180114 msg-sent id-404 1
//2006-01-02-150405
func db_msg_sent(line string) (v doc) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return nil
	}
	v = doc{
		"path":  fields[0],
		"time":  firstgame_time(fields[1]),
		"type":  fields[2],
		"mtype": fields[3],
		"cnt":   iconvert2(fields[4]),
	}
	return v
}

func firstgame_time(t string) (v time.Time) {
	const layout = "2006-01-02-150405"
	v, _ = time.Parse(layout, t)
	return
}
