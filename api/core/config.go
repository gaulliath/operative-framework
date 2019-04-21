package core

import "os"

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

func PushCore() *Core{
	return &Core{
		Host: os.Getenv("API_HOST"),
		Port: os.Getenv("API_PORT"),
		Verbose: os.Getenv("API_VERBOSE"),
	}
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


