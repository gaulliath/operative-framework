package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/graniet/operative-framework/session"
	"net/http"
)

type TargetsResponse struct{
	TargetId string `json:"target_id"`
	TargetName string `json:"target_name"`
	TargetType string `json:"target_type"`
	TargetLinked []TargetLink `json:"target_linked"`
}

type TargetLink struct{
	TargetId string `json:"target_id"`
	TargetName string `json:"target_name"`
	TargetType string `json:"target_type"`
	TargetResultId string `json:"target_result_id"`
}

type TargetInformationResponse struct{
	TargetId string `json:"target_id"`
	TargetName string `json:"target_name"`
	TargetType string `json:"target_type"`
	TargetResults map[string][]session.TargetResults `json:"target_results"`

}

func (api *ARestFul) Targets(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var targets []TargetsResponse
	for _, target := range api.sess.Targets{
		response := TargetsResponse{
			TargetId: target.GetId(),
			TargetName: target.GetName(),
			TargetType: target.GetType(),
		}
		if len(target.GetLinked()) > 0{
			for _, element := range target.GetLinked(){
				response.TargetLinked = append(response.TargetLinked, TargetLink{
					TargetId: element.TargetId,
					TargetName: element.TargetName,
					TargetType: element.TargetType,
					TargetResultId: element.TargetResultId,
				})
			}
		}
		targets = append(targets, response)
	}
	message := api.Core.PrintData("request executed", false, targets)
	_ = json.NewEncoder(w).Encode(message)
	return
}

func (api *ARestFul) Target(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	targetId := params["target_id"]
	t , err := api.sess.GetTarget(targetId)
	if err != nil{
		message := api.Core.PrintMessage("We can't found selected target", true)
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	targetInformationR := TargetInformationResponse{
		TargetId: t.GetId(),
		TargetName: t.GetName(),
		TargetType: t.GetType(),
		TargetResults: t.GetResults(),

	}
	message := api.Core.PrintData("request executed", false, targetInformationR)
	_ = json.NewEncoder(w).Encode(message)

}
