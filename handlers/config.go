package handlers

import (
	"encoding/json"
	"net/http"
	"projekat/service"
	"strconv"

	"github.com/gorilla/mux"
)

//Klasa
type ConfigHandler struct {
	service service.ConfigService
}

//Konstruktor
func NewConfigHandler(service service.ConfigService) ConfigHandler {
	return ConfigHandler{
		service: service,
	}
}


func(c ConfigHandler) AddConfig (writer http.ResponseWriter, request *http.Request){
	
}


// GET /configs/{name}/{version}
func (c ConfigHandler) GetConfig (writer http.ResponseWriter, request *http.Request) {
	name := mux.Vars(request)["name"];
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// pozovi servis metodu
	config, err := c.service.GetConfig(name, versionInt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	// vrati odgovor
	response, err := json.Marshal(config)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content−Type", "application/json")
	writer.Write(response)
}