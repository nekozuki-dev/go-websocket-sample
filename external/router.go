package external

import (
	"fmt"
	"go-websocket-sample/app"
	"go-websocket-sample/external/wsservice"
)

type Router interface {
	Run(port int)
}

func NewRouter() Router {
	return &router{}
}

type router struct {
}

func (r *router) Run(port int) {
	wsListener := wsservice.NewListener(port)
	wsListener.RegisterAcceptHandler(r.OnAccept)
	wsListener.RegisterCloseHandler(r.OnClose)
	wsListener.Run()
}

func (r *router) OnAccept(c wsservice.Conn) {
	fmt.Println("OnAccept")
	u := app.NewUser(c)
	u.Run()
}

func (r *router) OnClose(c wsservice.Conn) {
	fmt.Println("OnClose")
}
