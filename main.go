package main

import (
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/graniet/operative-framework/api"
	"github.com/graniet/operative-framework/modules"
	"github.com/graniet/operative-framework/session"
	"github.com/joho/godotenv"
	"io"
	"strings"
)


func main(){
	c := color.New(color.BgYellow).Add(color.FgBlack)
	_, _ = c.Println("operative framework - digital investigation framework")
	sess := session.New()
	err := godotenv.Load(".env")
	if err != nil {
		sess.Stream.Error("Please rename/create .env file on root path.")
		return
	}

	sess.PushPrompt()
	apiRest := api.PushARestFul(sess)
	modules.LoadModules(sess)

	l, err := readline.NewEx(sess.Prompt)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for{
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "api.run"{
			go apiRest.Start()
		} else {
			sess.ParseCommand(line)
		}
	}
}