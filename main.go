package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/graniet/operative-framework/api"
	"github.com/graniet/operative-framework/engine"
	"github.com/graniet/operative-framework/session"
	"github.com/joho/godotenv"
	"io"
	"os"
	"strings"
)


func main(){
	var sess *session.Session
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Please rename/create .env file on root path.")
		return
	}
	parser := argparse.NewParser("operative-framework", "digital investigation framework")
	rApi := parser.Flag("a", "api", &argparse.Options{
		Required: false,
		Help: "Load instantly operative framework restful API",
	})
	verbose := parser.Flag("v","verbose", &argparse.Options{
		Required: false,
		Help: "Do not show modules messages response",
	})
	cli := parser.Flag("n", "no-cli", &argparse.Options{
		Required: false,
		Help: "Do not run framework cli",
	})
	loadSession := parser.Int("s", "session", &argparse.Options{
		Required: false,
		Help: "Load specific session id",
	})

	help := parser.Flag("h", "help", &argparse.Options{
		Required: false,
		Help: "Print help",
	})

	err = parser.Parse(os.Args)
	if err != nil{
		fmt.Print(parser.Usage(err))
		return
	}
	if *loadSession > 0{
		sess = engine.Load(*loadSession)
	} else{
		sess = engine.New()
	}

	sess.PushPrompt()
	apiRest := api.PushARestFul(sess)
	if *help{
		fmt.Print(parser.Usage(""))
		return
	}
	if *rApi{
		if *cli{
			sess.Stream.Standard("running operative framework api...")
			sess.Stream.Standard("available at : " + apiRest.Server.Addr)
			sess.Information.SetApi(true)
			apiRest.Start()
		} else{
			sess.Stream.Standard("running operative framework api...")
			go apiRest.Start()
			sess.Stream.Standard("available at : " + apiRest.Server.Addr)
			sess.Information.SetApi(true)
		}
	}

	if *verbose{
		sess.Stream.Verbose = false
	} else{
		c := color.New(color.FgYellow)
		_, _ = c.Println("OPERATIVE FRAMEWORK - DIGITAL INVESTIGATION FRAMEWORK")
	}


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
		if line == "api run"{
			sess.Stream.Success("API Rest as been started at http://" + sess.Config.Api.Host + ":" + sess.Config.Api.Port)
			go apiRest.Start()
			sess.Information.SetApi(true)
		} else if line == "api stop"{
			_ = apiRest.Server.Close()
			sess.Information.SetApi(false)
		} else {
			if !engine.CommandBase(line, sess) {
				sess.ParseCommand(line)
			}
		}
	}
}