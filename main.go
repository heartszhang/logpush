package main

import (
	"bufio"
	"flag"
	"log"
	"net"

	"github.com/heartszhang/rulematch"
	"github.com/jeromer/syslogparser/rfc3164"
)

type packet map[string]interface{}
type decoder func(string) packet

var rule_manager = struct {
	tag_decoders  map[string]decoder
	word_decoders []decoder
	word_rules    [][]string
	_matcher      rulematch.Matcher
}{tag_decoders: map[string]decoder{}}

func TagDecoder(tag string, decoder decoder) {
	rule_manager.tag_decoders[tag] = decoder
}
func WordDecoder(rule []string, decoder decoder) {
	rule_manager.word_rules = append(rule_manager.word_rules, rule)
	rule_manager.word_decoders = append(rule_manager.word_decoders, decoder)
}
func build_rule_manager() {
	rule_manager._matcher = rulematch.NewMatcher(rule_manager.word_rules...)
}

var option = struct {
	sock       string
	redis_addr string
	redis_port uint
	verbose    bool
}{
	sock:       "localhost:4514",
	redis_addr: "127.0.0.1",
	redis_port: 6379,
}

func init() {
	flag.StringVar(&option.sock, "sock", option.sock, "rsyslog upstream socket")
	flag.BoolVar(&option.verbose, "verbose", option.verbose, "verbose mode")
	build_rule_manager()
}

func main() {
	flag.Parse()
	ln, err := net.Listen("tcp", option.sock)
	if err != nil {
		log.Fatal(err)
	}
	for {
		if conn, err := ln.Accept(); err == nil {
			go handle_connection(conn)
		} else {
			log.Fatal(err)
		}
	}
}

func handle_connection(conn net.Conn) {
	log.Println("client-start", conn.RemoteAddr())
	defer conn.Close()
	packet_chan := make(chan packet, 32)
	defer close(packet_chan)
	go handle_syslog(packet_chan, conn)
	var err error
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() && err == nil {
		err = rsyslog_publish(scanner.Bytes(), packet_chan)
	}
	log.Println("client-close", conn.RemoteAddr(), err)
}

func rsyslog_publish(data []byte, channel chan packet) (v error) {
	dec := rfc3164.NewParser(data) // why we use 3164?
	var body packet
	if err := dec.Parse(); err == nil {
		body = packet(dec.Dump())
		if body["tag"].(string) == "" {
			body["tag"] = "default"
		}
		body = decode_content(body)
	} else {
		body = packet{
			"tag":     "unknown",
			"content": data,
			"error":   err,
		}
	}
	channel <- body
	return v
}

func handle_syslog(channel chan packet, conn net.Conn) {
	//	rediscli := redis.New()
	//	err := rediscli.Connect(option.redis_addr, option.redis_port)
	//	if err != nil {
	//		conn.Close()
	//		log.Println(err)
	//		return
	//	}
	//	defer rediscli.Close()
	for body := range channel {
		log_packet(body)
		//		channel := fmt.Sprintf("%v", body["tag"])
		//		jbody, _ := json.Marshal(body)
		//		if _, err = rediscli.Publish(channel, jbody); err != nil {
		//			conn.Close()
		//			log.Println(err)
		//			break
		//		}
	}
}

/*
   "timestamp": p.header.timestamp,
   "hostname":  p.header.hostname,
   "tag":       p.message.tag,
   "content":   p.message.content,
   "priority":  p.priority.P,
   "facility":  p.priority.F.Value,
   "severity":  p.priority.S.Value,
*/

func decode_content(body packet) (v packet) {
	v = decode_content_by_tag(body)
	if v == nil {
		v = decode_content_by_words(body)
	}
	if v == nil {
		v = packet{"content": body["content"]}
	}
	v["tag"] = body["tag"]
	v["hostname"] = body["hostname"]
	return
}

func decode_content_by_tag(body packet) (v packet) {
	if decoder, ok := rule_manager.tag_decoders[body["tag"].(string)]; ok {
		v = decoder(body["content"].(string))
	}
	return
}

func decode_content_by_words(body packet) (v packet) {
	content := body["content"].(string)
	rules := rule_manager._matcher.Match(content)
	for _, idx := range rules {
		if v = rule_manager.word_decoders[idx](content); v != nil {
			break
		}
	}

	return v
}

func log_packet(p packet) {
	if !option.verbose {
		return
	}
	for k, v := range p {
		log.Println(k, "\t=>", v)
	}
}
