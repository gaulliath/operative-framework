package core

type Core struct {
	Host    string
	Port    string
	Verbose string
}

type ReturnMessage struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   bool        `json:"error"`
}

func (c *Core) PrintData(mess string, e bool, data interface{}) ReturnMessage {
	return ReturnMessage{
		Message: mess,
		Error:   e,
		Data:    data,
	}
}

func (c *Core) PrintMessage(mess string, e bool) ReturnMessage {
	return ReturnMessage{
		Message: mess,
		Error:   e,
	}
}
