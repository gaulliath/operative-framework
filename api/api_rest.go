package api

import (
	"github.com/graniet/operative-framework/session"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"github.com/graniet/operative-framework/api/core"
)

type ARestFul struct{
	sess *session.Session
	Server *http.Server
	Core *core.Core
}

func PushARestFul(s *session.Session) *ARestFul{
	mod := ARestFul{
		sess: s,
		Core: core.PushCore(),
	}
	mod.Server = &http.Server{
		Addr: mod.Core.Host + ":" + mod.Core.Port,
	}
	return &mod
}

func (api *ARestFul) Start(){
	r := mux.NewRouter()
	r.HandleFunc("/api/modules", api.Modules).Methods("GET")
	r.HandleFunc("/api/module/{module}", api.Module).Methods("GET")
	r.HandleFunc("/api/module/{module}", api.RunModule).Methods("POST")

	r.HandleFunc("/api/targets", api.Targets).Methods("GET")
	r.HandleFunc("/api/target/{target_id}", api.Target).Methods("GET")

	api.Server.Handler = r
	err := api.Server.ListenAndServe()
	if err != nil{
		fmt.Println(err.Error())
		return
	}
}
