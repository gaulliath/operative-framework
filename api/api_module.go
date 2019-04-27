package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Module struct{
	Name string
	Description string
	Author string
}

func (api *ARestFul) Modules(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var modules []Module
	for _, module := range api.sess.Modules{
		modules = append(modules, Module{
			Name: module.Name(),
			Description: module.Description(),
			Author: module.Author(),
		})
	}
	message := api.Core.PrintData("requests executed", false, modules)
	_ = json.NewEncoder(w).Encode(message)
}

func (api *ARestFul) RunModule(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	_ = r.ParseForm()
	if _, ok := r.Form["module"]; !ok{
		message := api.Core.PrintMessage("Argument 'module' as required.", true)
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	mod, err := api.sess.SearchModule(r.Form.Get("module"))
	if err != nil{
		message := api.Core.PrintMessage("A selected module as been note found", true)
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	moduleInformation := mod.GetInformation()
	for _, parameter := range moduleInformation.Parameters{
		if parameter.IsRequired{
			if _, ok := r.Form[parameter.Name]; !ok{
				message := api.Core.PrintMessage("Argument '" + parameter.Name + "' as required.", true)
				_ = json.NewEncoder(w).Encode(message)
				return
			}

			if r.Form.Get(parameter.Name) == ""{
				message := api.Core.PrintMessage("Argument '" + parameter.Name + "' required value.", true)
				_ = json.NewEncoder(w).Encode(message)
				return
			}
			_, _ = mod.SetParameter(parameter.Name, r.Form.Get(parameter.Name))
		} else{
			if _, ok := r.Form[parameter.Name]; ok{
				_, _ = mod.SetParameter(parameter.Name, r.Form.Get(parameter.Name))
			}
		}
	}
	mod.Start()
	result := mod.GetExport()
	message := api.Core.PrintData("request executed", false, result)
	_ = json.NewEncoder(w).Encode(message)
	return
}

func (api *ARestFul) Module(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	mod, err := api.sess.SearchModule(param["module"])
	if err != nil{
		message := api.Core.PrintMessage("A selected module as been note found", true)
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	api.sess.Stream.Verbose = false
	message := api.Core.PrintData("request executed", false, mod.GetInformation())
	_ = json.NewEncoder(w).Encode(message)
	api.sess.Stream.Verbose = true
	return
}
