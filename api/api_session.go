package api

import (
	"encoding/json"
	"net/http"
)

type OneSession struct{
	Id int `json:"id"`
	SessionName string `json:"session_id"`
}

func (OneSession) TableName() string{
	return "sessions"
}

func (api *ARestFul) Sessions(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var sessions []OneSession
	api.sess.Connection.ORM.Find(&sessions)
	message := api.Core.PrintData("request executed", false, sessions)
	_= json.NewEncoder(w).Encode(message)
}
