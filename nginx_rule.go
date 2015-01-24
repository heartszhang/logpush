package main

import (
	"bufio"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	TagDecoder("ali-access", new_nginx_decoder())
	WordDecoder([]string{"GET", "HTTP/1.1"}, new_nginx_decoder())
}

//118.254.176.199 - - [21/Jan/2015:13:27:22 +0800] "GET /firstGame/Android/funs0004/0.1.0.60_0.0.0.0.1/211/2a130f6b41ca86564b8428880733a399c/43ec1859-893d-45b9-942f-327236410a5e/70/24331/MissionBegin/3_9 HTTP/1.1" 200 151 "-" "-"
//$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"
const nginx_format = `(\d+\.\d+\.\d+\.\d+)\s\S+\s\S+\s\[([^\]]+)\]\s"(\S+)\s(\S+)\sHTTP/(\d+\.\d+)"\s(\d+)\s(\d+)\s"([^"]+)"\s"([^"]+)"`

type nginx_decoder struct {
	re *regexp.Regexp
}

var decoders map[string]func([]string) doc //= make(map[string]func([]string) doc)

func new_nginx_decoder() decoder {
	return nginx_decoder{
		re: regexp.MustCompile(nginx_format),
	}.decode
}

func (this nginx_decoder) decode(content string) (v doc) {
	fields := this.re.FindStringSubmatch(content)
	//remote_addr, time_local, request_method, request_url, request_ver, status, bytes_sent, http_refer, http_user_agent
	if len(fields) >= 10 {
		const time_layout = `02/Jan/2006:15:04:05 -0700`
		v = make(doc)
		v["remote"] = fields[1]
		v["time"], _ = time.Parse(time_layout, fields[2])
		v["verb"] = fields[3]
		v["url"] = fields[4]
		v["status"], _ = strconv.Atoi(fields[6])
		v["sent"], _ = strconv.Atoi(fields[7])
		v["http_refer"] = fields[8]
		v["http_ua"] = fields[9]
		v["type"] = "access"
		v.merge(decode_func(fields[4]))
	}
	return
}
func decode_version(v string) doc {
	return nil
}
func decode_func(line string) (v doc) {
	fields := split_fields(line) // fields[10] = keyword
	if len(fields) < 11 {
		log.Println(line)
		return
	}
	v = make(doc)
	var names = []string{"game", "client", "distributor", "version", "server", "dev", "uid", "seq", "player"}
	for idx, name := range names {
		v[name] = fields[idx+1] // ignore the first empty field
	}
	if fields[5] == "serverUnknown" {
		v["server"] = "-"
	}
	if fields[7] == "userUnknown" {
		v["uid"] = "-"
	}
	v.merge(decode_version(fields[4]))
	v["seq"], _ = iconvert(fields[8])
	if p, ok := iconvert(fields[9]); ok {
		v["player"] = p
	} else {
		log.Println(line)
		return
	}
	key := strings.ToLower(fields[10])
	v["type"] = "access/" + key

	if decoder, ok := decoders[key]; ok {
		v.merge(decoder(fields[11:]))
	}
	return
}

func split_fields(line string) (v []string) {
	scaner := bufio.NewScanner(strings.NewReader(line))
	scaner.Split(by_slash_vertical)
	for scaner.Scan() {
		v = append(v, scaner.Text())
	}
	return
}

func by_slash_vertical(data []byte, ateof bool) (advance int, token []byte, err error) {
	var f int
	for idx, b := range data {
		if b == '/' || b == '|' {
			f = 1
			advance = idx + f
			break
		}
	}
	if advance == 0 && ateof {
		advance = len(data)
	}
	if advance > 0 {
		token = data[:advance-f]
	}
	return
}
func sconvert(s string) (interface{}, bool) {
	return s, true
}
func iconvert(s string) (interface{}, bool) {
	if v, err := strconv.Atoi(s); err == nil {
		return v, true
	}
	return s, false
}
func fconvert(s string) (interface{}, bool) {
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return v, true
	}
	return s, false
}
func np(_types []int, names ...string) func([]string) doc {
	converters := []func(string) (interface{}, bool){sconvert, iconvert, fconvert}
	return func(fields []string) doc {
		v := make(doc)
		if len(fields) >= len(names) {
			for idx, name := range names {
				v[name], _ = converters[_types[idx]](fields[idx])
			}
		}
		return v
	}
}
func init() {
	var z1, z2, z3, o1, o2 = []int{0}, []int{0, 0}, []int{0, 0, 0}, []int{1}, []int{1, 1}
	decoders = map[string]func([]string) doc{
		"serverlisterror":       np(z3, "dev", "version", "resver"),
		"serverlistsucc":        np(z3, "dev", "version", "resver"),
		"levelupgrade":          np(o1, "to"),
		"login":                 np([]int{0, 1}, "result", "net"), //0: unknown, 1:wifi, 2: mobile
		"registeriospush":       np(z1, "pid"),
		"event":                 np(z1, "eid"),
		"registerpush":          np(z1, "eid"),
		"pushclick":             np(z1, "eid"),
		"exception":             np(z2, "err", "desc"),
		"purchase":              np([]int{0, 1, 2}, "name", "count", "price"),
		"use":                   np([]int{0, 1}, "name", "count"),
		"missionbegin":          np([]int{0}, "mid"),
		"missioncompleted":      np([]int{0}, "mid"),
		"missionfailed":         np([]int{0, 0}, "mid", "reason"),
		"chargerequest":         np([]int{0, 0, 2, 0, 2, 0}, "order", "iap", "price", "currency", "virtual_value", "payment"),
		"chargesuccess":         np(z1, "order"),
		"reward":                np([]int{2, 0}, "count", "reason"),
		"foregound":             np([]int{0, 0, 0, 1}, "platform", "version", "dev", "diff"),
		"starttask":             np([]int{1, 0, 1}, "player", "dev", "task"),
		"finishtask":            np([]int{1, 0, 1}, "player", "dev", "task"),
		"jumptask":              np([]int{1, 0, 1, 1, 1}, "player", "dev", "task", "stepdone", "steps"),
		"guideidhadshow":        np([]int{1, 0}, "player", "guide"),
		"visitorlogin":          np([]int{0, 1}, "result", "net"),
		"message":               np(o2, "mtype", "mid"),
		"banner":                np(o1, "bid"),
		"activity":              np(o2, "atype", "level"),
		"updateversion":         np(o1, "level"),
		"updateversionfinished": np(z1, "version"),
		"downloadfilefail": np([]int{0, 0, 0, 1, 1, 0, 0, 1, 1, 1},
			"dev", "%version", "resver", "uptype",
			"resuptype", "version2", "resver2", "pack_size", "download_size", "threads"),
		"downloadadd":               download_params(),
		"downloadfinish":            download_params(),
		"downloadmergefail":         download_params(),
		"downloadunzipfail":         download_params(),
		"downloadcheckfail":         download_params(),
		"connectupdaterservererror": np(o1, "err"),
		"updateinfo":                np([]int{1, 1, 0, 0}, "uptype", "resuptype", "version2", "resver2"),
		"updatelocalspacenotenough": np(o2, "pack_size", "diskspace"),
		"updateremotefilesizeerror": np([]int{1, 0}, "file_size", "url"),
		"updateprogress":            np(o1, "p"),
		"writefileerror":            np(o1, "err"),
	}
}
func download_params() func([]string) doc {
	return np([]int{0, 0, 0, 1, 1, 0, 0}, "dev", "%version", "resver", "uptype", "resuptype", "version2", "resver2")
}

/*
@game/@client/@distributor/@version(+resver)/@server/@dev/@uid/@sid(int)/@player(int)/
guestLoginClick
forceUpdateCancel
serverListError/@dev/@version(ver)/@resver
serverListSucc/@dev/@version/@resver
accountLoginClick
LevelUpgrade/@to(int)
Start
Login/@result/@net(int unknown=0, wifi=1, mobile=2)
RegisterIOSPush/@pid;
Event/@eid;
RegisterPush/@token;
PushClick/@ptype;
Exception/@err/@desc;
Purchase/%name/%count(int)/@price(num)
Use/%name/%count(int)
MissionBegin/@mid;
MissionCompleted/@mid
MissionFailed/@mid/@reason
ChargeRequest/@order/@iap/@price(num)/@currency/@virtual_value(num)/@payment
ChargeSuccess/@order
Reward/%count(num)/%reason
foregound|%platform|%version|%dev|%diff(\d+)|
startTask|%player(\d+)|%dev|%task(\d+)
finishTask|%player(\d+)|%dev|%task(\d+)
jumpTask|%player(\d+)|%dev|%task(\d+)|%stepdone(\d+)|%steps(\d+)
GuideIDHadShow|%player|%guide
VisitorLogin/Success/%net(\d+)/0
VisitorLogin/Failed/%net/0
Message/%mtype(\d+)/%mid(\d+)/
Banner/%bid(\d+)/
Activity/%atype(\d+)/%level(\d+)
CheckIn
UpdateVersion/%level/0/0
UpdateVersionfinished/%version
DownloadAdd|%dev|%version|%resver|%uptype(\d+)|%resuptype(int)|%version2|%resver2|
DownloadFinish|%dev|%version|%resver|%uptype(\d)|%resuptype(\d+)|%version2|%resver2|
DownloadFileFail|%dev|%version|%resver|%uptype(\d)|%resuptype(\d)|%version2|%resver2|%pack_size(\d+)|%download_size(\d+)|%threads(\d+)|
DownloadMergeFail|%dev|%version|%resver|%uptype(\d)|%resuptype(\d)|%version2|%resver2|
DownloadUnzipFail|%dev|%version|%resver|%uptype(\d)|%resuptype(\d)|%version2|%resver2|
DownloadCheckFail|%dev|%version|%resver|%uptype|%resuptype|%version2|%resver2|"
updaterDNSError/
connectUpdaterServerError/%err(\d+)/
updateInfo/%uptype(\d)/%resuptype(\d)/%version2/%resver2/
updateLocalSpaceNotEnough/%pack_size(\d+)/%diskspace(\d+)/
updateRemoteFileSizeError/%file_size(\d+)/%url/
updateProgress/%p(\d+)/
DownloadPartFilesFinish/
MergePartFilesFinish/
UnzipFileFinish/
CheckFilesFinish/
WriteFileError/%err(\d+)/
GetUpdateInfoFail/
*/
