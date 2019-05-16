package app

type (
	Packet struct {
		ID   uint16      `json:"id"`
		Body interface{} `json:"body"`
	}

	MessagePacket struct {
		Msg string `json:"msg"`
	}
)
