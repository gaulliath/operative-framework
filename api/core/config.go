package core

type Core struct{
	Host string
	Port string
	Verbose string
}

type ReturnMessage struct{
	Message string
	Data interface{}
	Error bool
}

func (c *Core) PrintData(mess string, e bool, data interface{}) ReturnMessage{
	return ReturnMessage{
		Message: mess,
		Error: e,
		Data: data,
	}
}

func (c *Core) PrintMessage(mess string, e bool) ReturnMessage{
	return ReturnMessage{
		Message: mess,
		Error: e,
	}
}


