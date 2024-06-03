package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projekat/model"
	"projekat/service"
	"strconv"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"
)

// Klasa
type ConfigGroupHandler struct {
	configGroupservice service.ConfigGroupService
	configService      service.ConfigService
	tracer             trace.Tracer
}

// Konstruktor
func NewConfigGroupHandler(configGroupservice service.ConfigGroupService, configService service.ConfigService, tracer trace.Tracer) ConfigGroupHandler {
	return ConfigGroupHandler{
		configGroupservice: configGroupservice,
		configService:      configService,
		tracer:             tracer,
	}
}

// GetConfigGroup retrieves a configuration group by name and version.
// @Summary Get a configuration group
// @Description Retrieves a configuration group by name and version
// @Tags configGroups
// @Produce json
// @Param name path string true "Name of the configuration group"
// @Param version path int true "Version of the configuration group"
// @Success 200 {object} model.ConfigGroup
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Configuration group not found"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/{name}/{version} [get]
func (c ConfigGroupHandler) GetConfigGroup(writer http.ResponseWriter, request *http.Request) {
	ctx, span := c.tracer.Start(request.Context(), "GetConfigGroup")
	defer span.End()

	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := c.configGroupservice.GetConfigGroup(ctx, name, versionInt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	response, err := json.Marshal(config)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content−Type", "application/json")
	writer.Write(response)
}

// @Summary Add a new configuration group
// @Description Adds a new configuration group
// @Tags configGroups
// @Accept json
// @Produce json
// @Param configGroup body model.ConfigGroup true "Configuration group to add"
// @Success 200 {string} string "Configuration group added successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups [post]
func (c ConfigGroupHandler) AddConfigGroup(writer http.ResponseWriter, request *http.Request) {
	ctx, span := c.tracer.Start(request.Context(), "AddConfigGroup")
	defer span.End()

	defer request.Body.Close()

	var configGroup model.ConfigGroup

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&configGroup)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}
	if configGroup.Name == "" {
		http.Error(writer, "Error: 'name' field is required and cannot be empty", http.StatusBadRequest)
		return
	}
	if configGroup.Version == 0 {
		http.Error(writer, "Error: 'version' field is required and cannot be zero", http.StatusBadRequest)
		return
	}
	if len(configGroup.Configurations) == 0 {
		http.Error(writer, "Error: 'config' field is required and cannot be empty", http.StatusBadRequest)
		return
	}

	configList := configGroup.Configurations
	for i := 0; i < len(configList); i++ {
		_, err := c.configService.GetConfig(ctx, configList[i].Name, configList[i].Version)
		if err != nil {
			c.configService.AddConfig(ctx, configList[i])
		}
	}
	err = c.configGroupservice.AddConfigGroup(ctx, configGroup)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(configGroup)
}

// @Summary Delete a configuration group
// @Description Deletes a configuration group by name and version
// @Tags configGroups
// @Param name path string true "Name of the configuration group"
// @Param version path int true "Version of the configuration group"
// @Success 200 {string} string "Configuration group deleted successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Configuration group not found"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/{name}/{version} [delete]
func (c ConfigGroupHandler) DeleteConfigGroup(writer http.ResponseWriter, request *http.Request) {
	ctx, span := c.tracer.Start(request.Context(), "DeleteConfigGroup")
	defer span.End()

	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, "Invalid version format", http.StatusBadRequest)
		return
	}

	err = c.configGroupservice.DeleteConfigGroup(ctx, name, versionInt)
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

// @Summary Add a configuration to a group
// @Description Adds a configuration to a specified group
// @Tags configGroups
// @Accept json
// @Produce json
// @Param name path string true "Name of the configuration group"
// @Param version path int true "Version of the configuration group"
// @Param config body model.Config true "Configuration to add to the group"
// @Success 200 {string} string "Configuration added to group successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Configuration group not found"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/{name}/{version}/configs [post]
func (c ConfigGroupHandler) AddConfigToGroup(writer http.ResponseWriter, request *http.Request) {
	ctx, span := c.tracer.Start(request.Context(), "AddConfigToGroup")
	defer span.End()

	defer request.Body.Close()

	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, "Invalid version format", http.StatusBadRequest)
		return
	}

	_, err = c.configGroupservice.GetConfigGroup(ctx, name, versionInt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	var config model.Config

	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&config)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}
	if config.Name == "" {
		http.Error(writer, "Error: 'name' field is required and cannot be empty", http.StatusBadRequest)
		return
	}
	if config.Version == 0 {
		http.Error(writer, "Error: 'version' field is required and cannot be zero", http.StatusBadRequest)
		return
	}
	if len(config.Parameters) == 0 {
		http.Error(writer, "Error: 'params' field is required and cannot be empty", http.StatusBadRequest)
		return
	}
	if len(config.Labels) == 0 {
		http.Error(writer, "Error: 'labels' field is required and cannot be empty", http.StatusBadRequest)
		return
	}

	group, err := c.configGroupservice.GetConfigGroup(ctx, name, versionInt)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	for _, c := range group.Configurations {
		if c.Name == config.Name && c.Version == config.Version {
			http.Error(writer, "Configuration already exists in the group", http.StatusBadRequest)
			return
		}
	}

	err = c.configService.AddConfig(ctx, config)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = c.configGroupservice.AddConfigToGroup(ctx, name, versionInt, config)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(config)
}

// @Summary Remove a configuration from a group
// @Description Removes a specific configuration from a specified group
// @Tags configGroups
// @Param groupName path string true "Name of the configuration group"
// @Param groupVersion path int true "Version of the configuration group"
// @Param configName path string true "Name of the configuration to remove"
// @Param configVersion path int true "Version of the configuration to remove"
// @Success 200 {string} string "Configuration removed from group successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Configuration or group not found"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/{groupName}/{groupVersion}/{configName}/{configVersion} [delete]
func (c ConfigGroupHandler) DeleteConfigFromGroup(writer http.ResponseWriter, request *http.Request) {
	ctx, span := c.tracer.Start(request.Context(), "DeleteConfigFromGroup")
	defer span.End()

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

	err = c.configGroupservice.DeleteConfigFromGroup(ctx, groupName, groupVersionInt, configName, configVersionInt)
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

// @Summary Retrieve configurations by labels
// @Description Retrieves all configurations within a specific group that match the specified labels
// @Tags configGroups
// @Produce json
// @Param groupName path string true "Name of the configuration group"
// @Param groupVersion path int true "Version of the configuration group"
// @Param labels path string true "Labels to filter the configurations"
// @Success 200 {array} model.Config "List of configurations"
// @Failure 400 {string} string "Invalid group version"
// @Failure 404 {string} string "Configurations not found"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/{groupName}/{groupVersion}/{labels} [get]
func (c ConfigGroupHandler) GetConfigsFromGroupByLabels(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.tracer.Start(r.Context(), "GetConfigsFromGroupByLabels")
	defer span.End()

	params := mux.Vars(r)
	groupName := params["groupName"]
	groupVersionStr := params["groupVersion"]
	labels := params["labels"]

	groupVersion, err := strconv.Atoi(groupVersionStr)
	if err != nil {
		http.Error(w, "Invalid group version", http.StatusBadRequest)
		return
	}

	configs, err := c.configGroupservice.GetConfigsFromGroupByLabel(ctx, groupName, groupVersion, labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(configs)
}

// @Summary Delete configurations by labels
// @Description Deletes all configurations within a specific group that match the specified labels
// @Tags configGroups
// @Param groupName path string true "Name of the configuration group"
// @Param groupVersion path int true "Version of the configuration group"
// @Param labels path string true "Labels to filter the configurations for deletion"
// @Success 200 {string} string "Configurations deleted successfully"
// @Failure 400 {string} string "Invalid group version"
// @Failure 404 {string} string "Configurations not found for deletion"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/{groupName}/{groupVersion}/{labels} [delete]
func (c ConfigGroupHandler) DeleteConfigsFromGroupByLabels(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.tracer.Start(r.Context(), "DeleteConfigsFromGroupByLabels")
	defer span.End()

	params := mux.Vars(r)
	groupName := params["groupName"]
	groupVersionStr := params["groupVersion"]
	labels := params["labels"]

	groupVersion, err := strconv.Atoi(groupVersionStr)
	if err != nil {
		http.Error(w, "Invalid group version", http.StatusBadRequest)
		return
	}

	err = c.configGroupservice.DeleteConfigsFromGroupByLabel(ctx, groupName, groupVersion, labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Configuration deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
