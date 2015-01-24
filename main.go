package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"time"

	"github.com/heartszhang/rulematch"
	"github.com/jeromer/syslogparser/rfc3164"
)

type decoder func(string) doc

type doc_decoder interface {
	decode(string) (p doc, doctype string)
}

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
	index_prefix           string
	dft_type               string
	sock                   string
	redis_addr, els_domain string
	redis_port, els_port   uint
	zmq_addr               string
	els_pool_num           int
	verbose                bool
}{
	sock:         "localhost:4514",
	redis_addr:   "127.0.0.1",
	redis_port:   6379,
	els_domain:   "localhost",
	els_port:     9200,
	els_pool_num: 16,
	index_prefix: "logstash",
	dft_type:     "generic",
	zmq_addr:     "ipc://stashlog.zmq.pull",
}

func init() {
	flag.StringVar(&option.sock, "sock", option.sock, "rsyslog upstream socket")
	flag.StringVar(&option.els_domain, "els-domain", option.els_domain, "elasticsearch working domain")
	flag.BoolVar(&option.verbose, "verbose", option.verbose, "verbose mode")
	option.index_prefix = option.index_prefix + "-" + time.Now().Format("20060102")
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
	doc_chan := make(chan doc, 32)
	defer close(doc_chan)
	//	go redis_pub(doc_chan, conn)
	//	go zmq_push(doc_chan, conn)
	go els_post(doc_chan, conn)
	var err error
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() && err == nil {
		err = rsyslog_publish(scanner.Bytes(), doc_chan)
	}
	log.Println("client-close", conn.RemoteAddr(), err)
}

func rsyslog_publish(data []byte, channel chan doc) (v error) {
	dec := rfc3164.NewParser(data) // why we use 3164?
	var body doc
	if err := dec.Parse(); err == nil {
		body = doc(dec.Dump())
		if body["tag"].(string) == "" {
			body["tag"] = "default"
		}
		body = decode_document(body)
	} else {
		body = doc{
			"tag":     "unknown",
			"content": data,
			"error":   err,
		}
	}
	channel <- body
	return v
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

func decode_document(body doc) (v doc) {
	v = decode_document_by_tag(body)
	if v == nil {
		v = decode_document_by_words(body)
	}
	if v == nil {
		v = doc{"content": body["content"]}
	}
	v["tag"] = body["tag"]
	v["hostname"] = body["hostname"]
	if _, ok := v["type"].(string); !ok {
		v["type"] = option.dft_type
	}
	return
}

func decode_document_by_tag(body doc) (v doc) {
	if decoder, ok := rule_manager.tag_decoders[body["tag"].(string)]; ok {
		v = decoder(body["content"].(string))
	}
	return
}

func decode_document_by_words(body doc) (v doc) {
	content := body["content"].(string)
	rules := rule_manager._matcher.Match(content)
	for _, idx := range rules {
		if v = rule_manager.word_decoders[idx](content); v != nil {
			break
		}
	}

	return v
}

func log_document(p doc) {
	if !option.verbose {
		return
	}
	for k, v := range p {
		log.Println(k, "\t=>", v)
	}
}
