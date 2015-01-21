package main

import (
	"regexp"
	"strconv"
	"time"
)

func init() {
	TagDecoder("ali-access", new_nginx_decoder())
	WordDecoder([]string{"GET", "HTTP/1.1"}, new_nginx_decoder())
}

//118.254.176.199 - - [21/Jan/2015:13:27:22 +0800] "GET /firstGame/Android/funs0004/0.1.0.60_0.0.0.0.1/211/2a130f6b41ca86564b8428880733a399c/43ec1859-893d-45b9-942f-327236410a5e/70/24331/MissionBegin/3_9 HTTP/1.1" 200 151 "-" "-"
//$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"
const nginx_format = `(\d+\.\d+\.\d+\.\d+)\s\S+\s\S+\s\[([^\]]+)\]\s"(\S+)\s(\S+)\sHTTP/(\d+\.\d+)"\s(\d+)\s(\d+)\s"([^"]+)"\s"([^"]+)"`

//var nginx_format_re *regexp.Regexp //= regexp.MustCompile()

//var nginx_format_re

type nginx_decoder struct {
	re *regexp.Regexp
}

func new_nginx_decoder() decoder {
	return nginx_decoder{re: regexp.MustCompile(nginx_format)}.decode
}

func (this nginx_decoder) decode(content string) (v packet) {
	fields := this.re.FindStringSubmatch(content)
	//remote_addr, time_local, request_method, request_url, request_ver, status, bytes_sent, http_refer, http_user_agent
	if len(fields) >= 10 {
		const time_layout = `02/Jan/2006:15:04:05 -0700`
		v = make(packet)
		v["remote"] = fields[1]
		v["time"], _ = time.Parse(time_layout, fields[2])
		v["verb"] = fields[3]
		v["url"] = fields[4]
		v["status"], _ = strconv.Atoi(fields[6])
		v["sent"], _ = strconv.Atoi(fields[7])
		v["http_refer"] = fields[8]
		v["http_ua"] = fields[9]
		v["type"] = "eaccess"
	}
	return
}
