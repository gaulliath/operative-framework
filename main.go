package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/graniet/operative-framework/api"
	"github.com/graniet/operative-framework/compiler"
	"github.com/graniet/operative-framework/engine"
	"github.com/graniet/operative-framework/services"
	"github.com/graniet/operative-framework/session"
	"github.com/graniet/operative-framework/supervisor"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
)


func main(){
	var sess *session.Session
	var sp *supervisor.Supervisor
	var configService string
	var opfBaseDirectory string
	var opfExport string

	// Load Configuration File
	configFile := ".env"
	err := godotenv.Load(".env")

	if err != nil {

		// Generate Default .env File
		u, errU := user.Current()
		if errU != nil{
			fmt.Println("Please create a .env file on root path.")
			return
		}
		if _, err := os.Stat(u.HomeDir + "/.opf/.env"); os.IsNotExist(err){
			if _, err := os.Stat(u.HomeDir + "/.opf/"); os.IsNotExist(err){
				_ = os.Mkdir(u.HomeDir + "/.opf/", os.ModePerm)
			}
			log.Println("Generating default .env on '"+u.HomeDir+"' directory...")
			path, errGeneration := engine.GenerateEnv(u.HomeDir + "/.opf/.env")
			if errGeneration != nil{
				log.Println(errGeneration.Error())
				return
			}
			err := godotenv.Load(path)
			if err != nil{
				log.Println(err.Error())
				return
			}
		}
		configFile = u.HomeDir + "/.opf/.env"
		configService = u.HomeDir + "/.opf/services/"
		opfBaseDirectory = u.HomeDir + "/.opf/"
		opfExport = opfBaseDirectory + "export/"
	}

	// Argument parser
	parser := argparse.NewParser("operative-framework", "digital investigation framework")
	rApi := parser.Flag("a", "api", &argparse.Options{
		Required: false,
		Help: "Load instantly operative framework restful API",
	})
	rSupervisor := parser.Flag("", "services", &argparse.Options{
		Required: false,
		Help: "Running supervised services.",
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

	scripts := parser.String("f", "opf", &argparse.Options{
		Required: false,
		Help: "Run script before prompt starting",
	})

	quiet := parser.Flag("q", "quiet", &argparse.Options{
		Required: false,
		Help: "Don't prompt operative shell",
	})

	err = parser.Parse(os.Args)
	if err != nil{
		fmt.Print(parser.Usage(err))
		return
	}

	// Checking if session as been specified
	if *loadSession > 0{
		sess = engine.Load(*loadSession)
	} else{
		sess = engine.New()
	}

	sess.PushPrompt()
	sess.Config.Common.ConfigurationFile = configFile
	sess.Config.Common.ConfigurationService = configService
	sess.Config.Common.BaseDirectory = opfBaseDirectory
	sess.Config.Common.ExportDirectory = opfExport
	apiRest := api.PushARestFul(sess)

	// Load supervised services.
	sp = supervisor.GetNewSupervisor(sess)
	services.Load(sp)

	if *rSupervisor{
		// Reading loaded services
		sp.Read()
		return
	}

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

	if *scripts != ""{
		compiler.Run(sess, *scripts)
	}

	if *quiet{
		return
	}

	if *verbose{
		sess.Stream.Verbose = false
	} else{
		c := color.New(color.FgYellow)
		_, _ = c.Println("OPERATIVE FRAMEWORK - DIGITAL INVESTIGATION FRAMEWORK")
		sess.Stream.WithoutDate("Loading a configuration file '" + configFile + "'")
		sess.Stream.WithoutDate("Loading a services configuration '"+sess.Config.Common.ConfigurationService+"'")
	}


	l, errP := readline.NewEx(sess.Prompt)
	if errP != nil {
		panic(errP)
	}
	defer l.Close()

	// Run Operative Framework Menu
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

		// Get Line With Trim Space
		line = strings.TrimSpace(line)

		// Checking Command
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