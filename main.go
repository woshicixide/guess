package main

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"guess/client"
	"guess/peer"
	"guess/server"
	"guess/vote"
	"io/ioutil"
	"log"
	"net/http"
	// _ "net/http/pprof"
	"runtime"
	"strconv"
	"strings"
	"ths/cookie"
	"ths/format"
)

var pubkey = ""
var bakpubkey = ""
var peers peer.Peers

// 投票
func makeVote(w http.ResponseWriter, r *http.Request) {
	// 解析return
	rt := r.FormValue("return")
	callback := r.FormValue("callback")
	if "json" != rt && "jsonp" != rt {
		fmt.Fprintln(w, format.New(-11, "缺少参数", nil).Fmt(""))
		return
	} else if "jsonp" == rt && "" == callback {
		fmt.Fprintln(w, format.New(-12, "缺少参数", nil).Fmt(""))
		return
	}

	// 解析fields
	fields := r.FormValue("fields")
	code, err := parseFields(fields)
	if nil != err {
		fmt.Fprintln(w, format.New(-13, "fields参数错误", nil).Fmt(callback))
		return
	}

	// 解析选项
	seqid, err := strconv.Atoi(r.FormValue("seqid"))
	if nil != err && !checkSeqid(seqid) {
		fmt.Fprintln(w, format.New(-14, "seqid参数错误", nil).Fmt(callback))
		return
	}
	stat := false
	if 1 == seqid {
		stat = true
	}

	// 解析cookie
	c := cookie.New(pubkey, bakpubkey)
	user, err := c.CheckCookieFromHead(r)
	if nil != err {
		fmt.Fprintln(w, format.New(-15, "请登录后投票", nil).Fmt(callback))
		return
	}
	userid := user.GetUserid()

	// vote
	vr, err := vote.Vote(code, userid, stat)
	if nil != err {
		fmt.Fprintln(w, format.New(-16, err.Error(), nil).Fmt(callback))
		return
	}

	// 返回数据
	fmt.Fprintln(w, format.New(0, "success", map[string]*vote.Result{fields: vr}).Fmt(callback))

	go sync(code, userid, stat)
}

// 获取投票
func getVote(w http.ResponseWriter, r *http.Request) {
	// 解析return
	rt := r.FormValue("return")
	callback := r.FormValue("callback")
	if "json" != rt && "jsonp" != rt {
		fmt.Fprintln(w, format.New(-11, "缺少参数", nil).Fmt(""))
		return
	} else if "jsonp" == rt && "" == callback {
		fmt.Fprintln(w, format.New(-12, "缺少参数", nil).Fmt(""))
		return
	}

	// 解析fields
	fields := r.FormValue("fields")
	code, err := parseFields(fields)
	if nil != err {
		fmt.Fprintln(w, format.New(-13, "fields参数错误", nil).Fmt(callback))
		return
	}

	// 解析cookie
	c := cookie.New(pubkey, bakpubkey)
	user, err := c.CheckCookieFromHead(r)
	var userid uint32
	if nil == err {
		userid = user.GetUserid()
	} else {
		userid = 0
	}

	fmt.Fprintln(w, format.New(0, "success", map[string]*vote.Result{fields: vote.GetVoteResult(code, userid)}).Fmt(callback))
}

// 清空投票
func clearVote(w http.ResponseWriter, r *http.Request) {
	salt := r.FormValue("salt")
	if pubkey != salt {
		fmt.Fprintln(w, format.New(-12, "fail", nil).Fmt(""))
		return
	}
	vote.Reset()
	fmt.Fprintln(w, format.New(0, "success", nil).Fmt(""))
	go reset()
}

// 解析fields
func parseFields(fields string) (string, error) {
	fs := strings.Split(fields, "_")
	if 2 != len(fs) {
		return "", errors.New("param fields is wrong")
	}
	if "stock" != fs[0] {
		return "", errors.New("param fields prefix is wrong")
	}
	return fs[1], nil
}

// 检查seqid是否合法
func checkSeqid(seqid int) bool {
	if 1 != seqid && 2 != seqid {
		return false
	}
	return true
}

func main() {
	http.HandleFunc("/clearVote", clearVote)
	http.HandleFunc("/makeVote", makeVote)
	http.HandleFunc("/getVote", getVote)
	err := http.ListenAndServe(":9090", nil)
	if nil != err {
		log.Fatal("ListenAndServe fail:", err)
	}
}

// 同步
func sync(code string, userid uint32, stat bool) {
	if len(peers) == 0 {
		return
	}
	client := new(client.Client)
	for _, v := range peers {
		client.Vote(v, code, userid, stat)
	}
}

// 同步
func reset() {
	if len(peers) == 0 {
		return
	}
	client := new(client.Client)
	for _, v := range peers {
		client.Reset(v)
	}
}

var c = new(conf)

type conf struct {
	Pubkey    string   `yaml:"pubkey"`
	Bakpubkey string   `yaml:"bakpubkey"`
	Peers     []string `yaml:"peers"`
}

// 初始化rpc服务端和配置文件
func init() {
	f := flag.String("conf", "./conf/guess.yaml", "conf file")
	flag.Parse()
	c.getConf(*f)
	pubkey = c.Pubkey
	bakpubkey = c.Bakpubkey
	peers = c.Peers

	go start()
}

func (self *conf) getConf(file string) {
	buffer, err := ioutil.ReadFile(file)
	if nil != err {
		log.Fatal("conf file does not exist")
	}

	yaml.Unmarshal(buffer, self)
	if nil != err {
		log.Fatal("parse conf fail")
	}
}

func start() {
	s := server.New()
	s.Start(":9091")
}
