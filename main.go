package main

import (
	"encoding/json"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/graniet/operative-framework/api"
	"github.com/graniet/operative-framework/cron"
	"github.com/graniet/operative-framework/engine"
	"github.com/graniet/operative-framework/export"
	"github.com/graniet/operative-framework/session"
	"github.com/graniet/operative-framework/supervisor"
	"github.com/segmentio/ksuid"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var sess *session.Session
	var sp *supervisor.Supervisor
	var files engine.Files

	files, err := engine.Preload()
	if err != nil {
		log.Println(err)
		return
	}

	// Argument parser
	parser := argparse.NewParser("operative-framework", "digital investigation framework")
	rApi := parser.Flag("a", "api", &argparse.Options{
		Required: false,
		Help:     "Load instantly operative framework restful API",
	})
	rSupervisor := parser.Flag("", "cron", &argparse.Options{
		Required: false,
		Help:     "Running supervised cron job(s).",
	})
	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Help:     "Do not show modules messages response",
	})
	cli := parser.Flag("n", "no-cli", &argparse.Options{
		Required: false,
		Help:     "Do not run framework cli",
	})
	execute := parser.String("e", "execute", &argparse.Options{
		Required: false,
		Help:     "Execute a single module",
	})
	target := parser.String("t", "target", &argparse.Options{
		Required: false,
		Help:     "Set target to '-e/--execute' argument",
	})
	parameters := parser.String("p", "parameters", &argparse.Options{
		Required: false,
		Help:     "Set parameters to '-e/--execute' argument",
	})
	loadSession := parser.Int("", "session-file", &argparse.Options{
		Required: false,
		Help:     "Load session cache file",
	})
	onlyModuleOutput := parser.Flag("", "no-banner", &argparse.Options{
		Required: false,
		Help:     "Do not print a banner information",
	})

	help := parser.Flag("h", "help", &argparse.Options{
		Required: false,
		Help:     "Print help",
	})

	file := parser.String("f", "file", &argparse.Options{
		Required: false,
		Help:     "Source file",
	})

	mode := parser.String("m", "mode", &argparse.Options{
		Required: false,
		Help:     "Start in specific mode: (perception, tracking, api, console): default (console)",
	})

	wh := parser.String("", "webhook", &argparse.Options{
		Required: false,
		Help:     "Autostart webHook by name",
	})

	quiet := parser.Flag("q", "quiet", &argparse.Options{
		Required: false,
		Help:     "Don't prompt operative shell",
	})

	modules := parser.Flag("l", "list", &argparse.Options{
		Required: false,
		Help:     "List available modules",
	})

	jsonExport := parser.Flag("", "json", &argparse.Options{
		Required: false,
		Help:     "Print result with a JSON format",
	})

	csvExport := parser.Flag("", "csv", &argparse.Options{
		Required: false,
		Help:     "Print result with a CSV format",
	})

	autoloadWH := parser.Flag("", "autoload-webhooks", &argparse.Options{
		Required: false,
		Help:     "Set active all 'web hooks' loaded in session",
	})

	eval := parser.String("", "eval", &argparse.Options{
		Required: false,
		Help:     "Execute commands while framework boot",
	})

	err = parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	// Checking if session as been specified
	if *loadSession > 0 {
		sess = engine.Load(*loadSession)
	} else {
		sess = engine.New()
	}

	sess.PushPrompt()
	sess.Config.Common.ConfigurationFile = files.Config
	sess.Config.Common.ConfigurationJobs = files.Job
	sess.Config.Common.BaseDirectory = files.Base
	sess.Config.Common.ExportDirectory = files.Export
	sess.ParseModuleConfig()
	apiRest := api.PushARestFul(sess)

	// Load supervised cron job.
	sp = supervisor.GetNewSupervisor(sess)
	cron.Load(sp)

	if *modules {
		sess.ListModules()
		return
	}

	if *rSupervisor {
		// Reading loaded cron job
		sp.Read()
		return
	}

	if *help {
		fmt.Print(parser.Usage(""))
		return
	}

	if *verbose {
		sess.Stream.Verbose = false
	}

	if *execute != "" {
		if *target == "" {
			sess.Stream.Error("'-e/--execute' argument need a target argument '-t/--target'")
			return
		}
		module, err := sess.SearchModule(*execute)
		if err != nil {
			sess.Stream.Error(err.Error())
			return
		}

		target, err := sess.AddTarget(module.GetType()[0], *target)
		if err != nil {
			sess.Stream.Error(err.Error())
			return
		}
		_, _ = module.SetParameter("TARGET", target)

		if *parameters != "" {

			if !strings.Contains(*parameters, "=") {
				sess.Stream.Error("Please use a correct format. example: limit=50;id=1")
				return
			}

			if strings.Contains(*parameters, ";") {
				lists := strings.Split(*parameters, ";")
				for _, parameter := range lists {
					template := strings.Split(parameter, "=")
					_, err := module.SetParameter(template[0], template[1])
					if err != nil {
						sess.Stream.Error(err.Error())
						return
					}
				}
			} else {
				template := strings.Split(*parameters, "=")
				_, err := module.SetParameter(template[0], template[1])
				if err != nil {
					sess.Stream.Error(err.Error())
					return
				}
			}
		}
		if *csvExport {
			sess.Stream.CSV = true
		}
		sess.NewInstance(module.Name())
		module.Start()

		if *jsonExport {
			e := export.JSON(sess)
			j, err := json.Marshal(e)
			if err == nil {
				print(string(j))
			}
			return
		}
		return
	}

	if *rApi {
		if *cli {
			sess.Stream.Standard("running operative framework api...")
			sess.Stream.Standard("available at : " + apiRest.Server.Addr)
			sess.Information.SetApi(true)
			apiRest.Start()
		} else {
			sess.Stream.Standard("running operative framework api...")
			go apiRest.Start()
			sess.Stream.Standard("available at : " + apiRest.Server.Addr)
			sess.Information.SetApi(true)
		}
	}

	// Checking if source file exists in argv
	if *file != "" {
		sess.SetSourceFile(*file)
		_ = sess.FromSourceFile()
	}

	// Load Webhooks configuration
	sess.LoadWebHook()

	if *wh != "" {
		webHook, err := sess.GetWebHookByName(*wh)
		if err != nil {
			sess.Stream.Error(err.Error())
			return
		}
		webHook.SetStatus(true)
	}

	if *autoloadWH {
		for _, wh := range sess.WebHooks {
			sess.Stream.Standard("Starting '" + wh.GetName() + "' web hooks")
			wh.Up()
		}
	}

	if *mode != "" {
		switch strings.ToLower(*mode) {
		case "perception":
			if *file == "" {
				sess.Stream.Error("Please select a source file (-f) with interval commands.")
				return
			}

			sess.Stream.Standard("Mode '" + strings.ToLower(*mode) + "' as started now...")
			select {}
		case "tracking":
			sess.Stream.Standard("running operative framework api...")
			sess.Stream.Standard("[API] available : " + apiRest.Server.Addr)
			go apiRest.Start()
			sess.Stream.Standard("[GUI] Tracking : " + sess.GetTrackingUrlWithParam())
			sess.ServeTrackerGUI()
			break
		case "console":
			break
		case "api":
			sess.Stream.Standard("running operative framework api...")
			sess.Stream.Standard("available at : " + apiRest.Server.Addr)
			sess.Information.SetApi(true)
			apiRest.Start()
			break
		default:
			sess.Stream.Warning("Mode '" + *mode + "' as unknown, running 'console' mode...")
			break
		}
	}

	if *quiet {
		return
	}

	if !*onlyModuleOutput && sess.Stream.Verbose {
		c := color.New(color.FgYellow)
		_, _ = c.Println("OPERATIVE FRAMEWORK - DIGITAL INVESTIGATION FRAMEWORK")
		sess.Stream.WithoutDate("Loading a configuration file '" + files.Config + "'")
		sess.Stream.WithoutDate("Loading a cron job configuration '" + sess.Config.Common.ConfigurationJobs + "'")
		sess.Stream.WithoutDate("Loading '" + strconv.Itoa(len(sess.Config.Modules)) + "' module(s) configuration(s)")
	}

	admin, err := sess.GetUser("admin")
	if admin == nil {
		password := ksuid.New().String()
		_, err := sess.NewUser("admin", password)
		if err != nil {
			sess.Stream.Error(err.Error())
			return
		}
	}

	l, errP := readline.NewEx(sess.Prompt)
	if errP != nil {
		panic(errP)
	}
	defer l.Close()

	// Checking in background available interval
	go sess.WaitInterval()

	// Checking in background available monitor
	go sess.WaitMonitor()

	if *eval != "" {
		sess.ParseCommands(*eval)
	}

	// Run Operative Framework Menu
	for {

		sess.UpdatePrompt()
		l, errP := readline.NewEx(sess.Prompt)
		if errP != nil {
			panic(errP)
		}

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
		if line == "api run" {
			sess.Stream.Success("API Rest as been started at http://" + sess.Config.Api.Host + ":" + sess.Config.Api.Port)
			go apiRest.Start()
			sess.Information.SetApi(true)
		} else if line == "api stop" {
			_ = apiRest.Server.Close()
			sess.Information.SetApi(false)
		} else if line == "tracker run" {
			sess.Stream.Success("[GUI] Tracking : " + sess.GetTrackingUrlWithParam())
			go sess.ServeTrackerGUI()
			sess.Information.SetTracker(true)
		} else if line == "tracker stop" {
			_ = sess.Tracker.Server.Close()
			sess.Information.SetTracker(false)
		} else {
			if !engine.CommandBase(line, sess) {
				sess.ParseCommand(line)
			}
		}
	}
}
