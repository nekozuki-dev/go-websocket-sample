package app

type MessageHandleFunc func(interface{})

type MessageHandlers interface {
	Get(msgid uint16) MessageHandleFunc
	Register(msgid uint16, handler MessageHandleFunc)
	Unregister(msgid uint16)
}

func NewMessageHandlers() MessageHandlers {
	return &messageHandlers{
		handlers: make(map[uint16]MessageHandleFunc),
	}
}

type messageHandlers struct {
	handlers map[uint16]MessageHandleFunc
}

func (m *messageHandlers) Get(msgid uint16) MessageHandleFunc {
	return m.handlers[msgid]
}

func (m *messageHandlers) Register(msgid uint16, handler MessageHandleFunc) {
	m.handlers[msgid] = handler
}

func (m *messageHandlers) Unregister(msgid uint16) {
	delete(m.handlers, msgid)
}
