package vote

import (
	"errors"
	// "log"
	"sync"
)

// 存储单位
type codeSto struct {
	userids map[uint32]byte
	z_num   int
	d_num   int
	rw      *sync.RWMutex
}

// 结果拼装
type Result struct {
	Theme   string    `json:"theme"`
	Options [2]option `json:"options"`
	Voted   byte      `json:"voted,omitempty"`
}

// 子结果拼装
type option struct {
	Optionid int    `json:"optionid"`
	Name     string `json:"name"`
	Hot      int    `json:"hot"`
}

// 过票
func (self *codeSto) vote(userid uint32, stat bool) {
	self.rw.Lock()
	defer self.rw.Unlock()

	if stat {
		self.userids[userid] = 1
		self.z_num++
	} else {
		self.userids[userid] = 2
		self.d_num++
	}
}

// 过票
func (self *codeSto) getResult(userid uint32) *Result {
	r := &Result{
		Theme:   "猜涨跌",
		Options: [2]option{option{1, "看涨", self.z_num}, option{2, "看跌", self.d_num}},
	}
	if userid > 0 {
		r.Voted = self.isVote(userid)
	}

	return r
}

// 获取投票结果
func GetVoteResult(code string, userid uint32) *Result {
	return getVote(code).getResult(userid)
}

// 判断是否投过票
func (self *codeSto) isVote(userid uint32) byte {
	self.rw.RLock()
	defer self.rw.RUnlock()
	if stat, isVote := self.userids[userid]; isVote {
		return stat
	}
	return 0
}

func New() *codeSto {
	return &codeSto{userids: make(map[uint32]byte), z_num: 0, d_num: 0, rw: new(sync.RWMutex)}
}

// 获取投票数据
func getVote(code string) *codeSto {
	return con.get(code)
}

// 投票
func Vote(code string, userid uint32, stat bool) (*Result, error) {
	// 检查代码是否已存在
	v := getVote(code)
	// log.Fatalln(v)
	r := new(Result)

	// 判断是否已经投过票
	if int(v.isVote(userid)) > 0 {
		return r, errors.New("已投过票")
	}
	v.vote(userid, stat)
	voted := 1
	if false == stat {
		voted = 2
	}

	// 拼装投票后的返回结果
	r.Theme = "猜涨跌"
	r.Options = [2]option{option{1, "看涨", v.z_num}, option{2, "看跌", v.d_num}}
	r.Voted = byte(voted)

	return r, nil
}
