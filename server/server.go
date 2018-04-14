package server

import (
	"guess/backend"
	"guess/vote"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Server struct{}

func New() *Server {
	return new(Server)
}

func (self *Server) Start(port string) {
	rpc.Register(self)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", port)
	if nil != e {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

func (self *Server) Vote(args *backend.Args, rep *backend.Reply) error {
	if _, err := vote.Vote(args.Code, args.Userid, args.Stat); nil != err {
		*rep = false
		return err
	}
	*rep = true
	return nil
}

func (self *Server) ClearVote(args *bool, rep *bool) error {
	vote.Reset()
	*rep = true
	return nil
}

// func main() {
// 	server := newServer()
// 	server.start(":1234")
// }
