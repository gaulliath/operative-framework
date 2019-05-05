package engine

import (
	"errors"
	"log"
	"os"
	"os/user"
)

type Environment struct{
	Param []EnvParam
}

type EnvParam struct{
	Name string
	Value string
}

func (e *Environment) Add(name string, value string){

	// Add argument to virtual env
	e.Param = append(e.Param, EnvParam{
		Name: name,
		Value: value,
	})
	return
}

func GenerateEnv(path string) (string, error){
	var Env Environment
	u, err := user.Current()
	database := "./opf.db"
	if err == nil{
		database = u.HomeDir + "/.opf/opf.db"
	}

	// Put .env arguments
	Env.Add("API_HOST","127.0.0.1")
	Env.Add("API_PORT", "8888")
	Env.Add("API_VERBOSE", "true")
	Env.Add("DB_NAME", database)
	Env.Add("DB_DRIVER", "sqlite3")
	Env.Add("DB_HOST", "")
	Env.Add("DB_USER", "")
	Env.Add("DB_PASS", "")
	Env.Add("OPERATIVE_HISTORY", "/tmp/operative_framework.tmp")
	Env.Add("INSTAGRAM_LOGIN", "")
	Env.Add("INSTAGRAM_PASSWORD", "")
	Env.Add("TWITTER_CONSUMER", "")
	Env.Add("TWITTER_CONSUMER_SECRET", "")
	Env.Add("TWITTER_ACCESS_TOKEN", "")
	Env.Add("TWITTER_ACCESS_TOKEN_SECRET", "")

	// Generate a .env
	var file *os.File
	var errPath error

	file, errPath = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if errPath != nil{
		return "", errors.New(errPath.Error())
	}
	defer file.Close()

	// Writing parameters
	for _, param := range Env.Param{
		if param.Value == "" {
			_, _ = file.WriteString(param.Name + "=\n")
		} else{
			_, _ = file.WriteString(param.Name + "=" + "\""+param.Value+"\"\n")
		}
	}

	log.Println("New environment file as been write '" + path + "'")
	return path, nil
}
