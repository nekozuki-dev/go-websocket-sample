package wsservice

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type AcceptHandler func(Conn)
type CloseHandler func(Conn)

type Listener interface {
	Run()
	RegisterAcceptHandler(AcceptHandler)
	RegisterCloseHandler(CloseHandler)
}

type listener struct {
	listenerAcync
	port          int
	upgrader      websocket.Upgrader
	acceptHandler AcceptHandler
	closeHandler  CloseHandler
}

type listenerAcync struct {
	m     sync.Mutex
	conns map[*websocket.Conn]Conn
}

func NewListener(port int) Listener {
	lis := &listener{
		port:     port,
		upgrader: websocket.Upgrader{},
	}
	lis.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	lis.conns = make(map[*websocket.Conn]Conn)
	return lis
}

func (lis *listener) Run() {
	http.HandleFunc("/", pageHandle)
	http.HandleFunc("/ws", lis.handleConnection)

	servAddr := fmt.Sprintf(":%d", lis.port)
	fmt.Println("BeginListener", servAddr)
	if err := http.ListenAndServe(servAddr, nil); err != nil {
		panic(err)
	}
}

const page = `
<html lang="ja">
    <head>
        <meta charset="UTF-8">
        <title>WebSocketSample</title>
    </head>
    <script>
        var sock = new WebSocket('ws://127.0.0.1:9080/ws');

        var send = function(msgid, body) {
            var packet = {
                'id': msgid,
                'body': body
            };
            var json = JSON.stringify(packet)
            sock.send(json)
        };

        sock.addEventListener('open', function(e) {
            console.log('Connect success.')
            document.getElementById('banana').addEventListener('click',function(e) {
                var msg = {
                    'msg': 'ばななをあげる'
                };
                send(1, msg)
            });
        });

        sock.addEventListener('close', function(e) {
            console.log('Connect close.')
        });

        sock.addEventListener('message', function(e) {
            var json = JSON.parse(e.data)
            var msgid = json.id;
            var body = json.body;
            if (msgid == 1) {
                console.log(body.msg);
            }
        });
    </script>
    <body>
        <input type="button" id="banana" value="バナナを送る" />
    </body>
</html>
`

func pageHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, page)
}

func (lis *listener) RegisterAcceptHandler(handler AcceptHandler) {
	lis.acceptHandler = handler
}

func (lis *listener) RegisterCloseHandler(handler CloseHandler) {
	lis.closeHandler = handler
}

func (lis *listener) handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := lis.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error", err.Error())
		return
	}
	defer lis.closeConnection(ws)

	addr := ws.RemoteAddr().String()
	fmt.Println("NewConnection", addr)

	c := NewConn(ws)
	lis.m.Lock()
	lis.conns[ws] = c
	lis.m.Unlock()

	if lis.acceptHandler != nil {
		lis.acceptHandler(c)
	}
}

func (lis *listener) closeConnection(ws *websocket.Conn) {
	addr := ws.RemoteAddr().String()
	fmt.Println("CloseConnection", addr)

	lis.m.Lock()
	c := lis.conns[ws]
	delete(lis.conns, ws)
	lis.m.Unlock()

	ws.Close()
	if lis.closeHandler != nil {
		lis.closeHandler(c)
	}
}
