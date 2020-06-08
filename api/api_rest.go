package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graniet/operative-framework/api/core"
	"github.com/graniet/operative-framework/session"
)

type ARestFul struct {
	sess   *session.Session
	Server *http.Server
	Core   *core.Core
}

func PushARestFul(s *session.Session) *ARestFul {
	c := core.Core{
		Host:    s.Config.Api.Host,
		Port:    s.Config.Api.Port,
		Verbose: s.Config.Api.Verbose,
	}
	mod := ARestFul{
		sess: s,
		Core: &c,
	}
	mod.Server = &http.Server{
		Addr: mod.Core.Host + ":" + mod.Core.Port,
	}
	return &mod
}

func (api *ARestFul) LoadRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/sessions", api.Sessions).Methods("GET")

	r.HandleFunc("/api/modules", api.Modules).Methods("GET")
	r.HandleFunc("/api/modules/{module}", api.Module).Methods("GET")
	r.HandleFunc("/api/modules", api.RunModule).Methods("POST")

	r.HandleFunc("/api/targets", api.Targets).Methods("GET")
	r.HandleFunc("/api/targets/type/{target_type}", api.TargetByType).Methods("GET")
	r.HandleFunc("/api/targets/{target_id}", api.Target).Methods("GET")
	r.HandleFunc("/api/targets/{target_id}/results", api.Results).Methods("GET")
	r.HandleFunc("/api/targets/{target_id}/results/{result_id}", api.Result).Methods("GET")

	r.HandleFunc("/api/intervals", api.Intervals).Methods("GET")

	r.HandleFunc("/api/trackers", api.Trackers).Methods("GET")
	r.HandleFunc("/api/trackers", api.PutTracker).Methods("POST")
	r.HandleFunc("/api/tracker/{identifier}", api.Tracker).Methods("GET")
	return r
}

func (api *ARestFul) Start() {
	r := api.LoadRouter()
	api.Server.Handler = r
	err := api.Server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
