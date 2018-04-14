package client

import (
	"guess/backend"
	"log"
	"net/rpc"
)

type Client struct{}

func (self *Client) Vote(address string, code string, userid uint32, stat bool) bool {
	// 捕获有可能从call方法中抛出的panic
	defer func() {
		if err := recover(); err != nil {
			log.Println("rpc call fail:", err)
		}
	}()
	client, err := rpc.DialHTTP("tcp", address)
	if nil != err {
		log.Println("rpc dial fail:", err)
		return false
	}
	defer client.Close()
	args := &backend.Args{code, userid, stat}
	var reply backend.Reply
	err = client.Call("Server.Vote", args, &reply)
	if nil != err {
		log.Println("rpc vote fail:", err)
	}
	return bool(reply)
}

func (self *Client) Reset(address string) {
	// 捕获有可能从call方法中抛出的panic
	defer func() {
		if err := recover(); err != nil {
			log.Println("rpc call fail:", err)
		}
	}()

	client, err := rpc.DialHTTP("tcp", address)
	if nil != err {
		log.Println("rpc dial fail:", err)
		return
	}
	defer client.Close()
	var reply, args bool
	err = client.Call("Server.ClearVote", &reply, &args)
	if nil != err {
		log.Println("rpc clearvote fail:", err)
	}
}
