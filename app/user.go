package app

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type User interface {
	Run()
	Write(msgid uint16, body interface{})
}

func NewUser(c Conn) User {
	return &user{
		conn:        c,
		msgHandlers: NewMessageHandlers(),
		readCh:      make(chan []byte),
	}
}

type user struct {
	conn        Conn
	msgHandlers MessageHandlers
	readCh      chan []byte
}

func (u *user) Run() {
	u.msgHandlers.Register(1, u.handleMessage)

	readCh := make(chan []byte)
	closeCh := make(chan bool)

	go u.conn.Run(readCh, closeCh)

	for {
		select {
		case bytes := <-readCh:
			u.doHandler(bytes)

		case <-closeCh:
			fmt.Println("CloseUser")
			return
		default:
		}
	}
}

func (u *user) Write(msgid uint16, body interface{}) {
	packet := &Packet{
		ID:   msgid,
		Body: body,
	}
	bytes, err := json.Marshal(packet)
	if err != nil {
		u.conn.Close()
		return
	}
	u.conn.Write(bytes)
}

func (u *user) doHandler(bytes []byte) error {
	packet := &Packet{}
	if err := json.Unmarshal(bytes, packet); err != nil {
		return err
	}

	handler := u.msgHandlers.Get(packet.ID)
	if handler != nil {
		handler(packet.Body)
	}
	return nil
}

func (u *user) handleMessage(body interface{}) {
	req := &MessagePacket{}
	if err := mapstructure.Decode(body, req); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(req.Msg)

	res := &MessagePacket{
		Msg: "ごちそうさまでした！",
	}
	u.Write(1, res)
}
