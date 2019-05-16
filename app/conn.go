package app

type Conn interface {
	Run(readCh chan []byte, closeCh chan bool)
	Write([]byte)
	Close()
}
