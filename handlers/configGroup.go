package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projekat/model"
	"projekat/service"
	"strconv"

	"github.com/gorilla/mux"
)

// Klasa
type ConfigGroupHandler struct {
	configGroupservice service.ConfigGroupService
	configService      service.ConfigService
}

// Konstruktor
func NewConfigGroupHandler(configGroupservice service.ConfigGroupService, configService service.ConfigService) ConfigGroupHandler {
	return ConfigGroupHandler{
		configGroupservice: configGroupservice,
		configService:      configService,
	}
}

// GET /configs/{name}/{version}
func (c ConfigGroupHandler) GetConfigGroup(writer http.ResponseWriter, request *http.Request) {
	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// pozovi servis metodu
	config, err := c.configGroupservice.GetConfigGroup(name, versionInt)
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

func (c ConfigGroupHandler) AddConfigGroup(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var configGroup model.ConfigGroup

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&configGroup)
	if err != nil {
		fmt.Fprintf(writer, "Error parsing JSON: %v", err)
		return
	}
	if configGroup.Name == "" {
		fmt.Fprintf(writer, "Error: 'name' field is required and cannot be empty")
		return
	}
	if configGroup.Version == 0 {
		fmt.Fprintf(writer, "Error: 'version' field is required and cannot be zero")
		return
	}
	if len(configGroup.Configurations) == 0 {
		fmt.Fprintf(writer, "Error: 'config' field is required and cannot be empty")
		return
	}

	configList := configGroup.Configurations
	for i := 0; i < len(configList); i++ {
		_, err := c.configService.GetConfig(configList[i].Name, configList[i].Version)
		if err != nil {
			c.configService.AddConfig(configList[i])
		}
	}
	c.configGroupservice.AddConfigGroup(configGroup)
	fmt.Fprintf(writer, "Received config: %+v", configGroup)
}

func (c ConfigGroupHandler) DeleteConfigGroup(writer http.ResponseWriter, request *http.Request) {
	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, "Invalid version format", http.StatusBadRequest)
		return
	}

	err = c.configGroupservice.DeleteConfigGroup(name, versionInt)
	if err != nil {
		if err.Error() == "config group not found" {
			http.Error(writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{"message": "Configuration group successfully deleted"}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "Failed to write response", http.StatusInternalServerError)
	}
}

func (c ConfigGroupHandler) AddConfigToGroup(writer http.ResponseWriter, request *http.Request) {

	defer request.Body.Close()

	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.configGroupservice.GetConfigGroup(name, versionInt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	var config model.Config

	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Fprintf(writer, "Error parsing JSON: %v", err)
		return
	}
	if config.Name == "" {
		fmt.Fprintf(writer, "Error: 'name' field is required and cannot be empty")
		return
	}
	if config.Version == 0 {
		fmt.Fprintf(writer, "Error: 'version' field is required and cannot be zero")
		return
	}
	if len(config.Parameters) == 0 {
		fmt.Fprintf(writer, "Error: 'params' field is required and cannot be empty")
		return
	}
	if len(config.Labels) == 0 {
		fmt.Fprintf(writer, "Error: 'labels' field is required and cannot be empty")
		return
	}

	group, err := c.configGroupservice.GetConfigGroup(name, versionInt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, c := range group.Configurations {
		if c.Name == config.Name && c.Version == config.Version {
			http.Error(writer, "Configuration already exists in the group", http.StatusBadRequest)
			return
		}
	}

	c.configService.AddConfig(config)
	c.configGroupservice.AddConfigToGroup(name, versionInt, config)
	fmt.Fprintf(writer, "Received config: %+v", config)
}

func (c ConfigGroupHandler) DeleteConfigFromGroup(writer http.ResponseWriter, request *http.Request) {
	groupName := mux.Vars(request)["groupName"]
	groupVersion := mux.Vars(request)["groupVersion"]
	configName := mux.Vars(request)["configName"]
	configVersion := mux.Vars(request)["configVersion"]

	groupVersionInt, err := strconv.Atoi(groupVersion)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	configVersionInt, err := strconv.Atoi(configVersion)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.configGroupservice.DeleteConfigFromGroup(groupName, groupVersionInt, configName, configVersionInt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	response := map[string]string{"message": "Configuration successfully deleted"}
	jsonResponse, _ := json.Marshal(response)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResponse)
}

func (c ConfigGroupHandler) GetConfigsFromGroupByLabels(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	groupName := params["groupName"]
	groupVersionStr := params["groupVersion"]
	labels := params["labels"]

	groupVersion, err := strconv.Atoi(groupVersionStr)
	if err != nil {
		http.Error(w, "Invalid group version", http.StatusBadRequest)
		return
	}

	configs, err := c.configGroupservice.GetConfigsFromGroupByLabel(groupName, groupVersion, labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(configs)
}

func (c ConfigGroupHandler) DeleteConfigsFromGroupByLabels(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	groupName := params["groupName"]
	groupVersionStr := params["groupVersion"]
	labels := params["labels"]

	groupVersion, err := strconv.Atoi(groupVersionStr)
	if err != nil {
		http.Error(w, "Invalid group version", http.StatusBadRequest)
		return
	}

	err = c.configGroupservice.DeleteConfigsFromGroupByLabel(groupName, groupVersion, labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprint(w, "Configuration deleted successfully")
}
