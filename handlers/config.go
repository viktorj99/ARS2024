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
type ConfigHandler struct {
	service service.ConfigService
	tracer  trace.Tracer
}

// Konstruktor
func NewConfigHandler(service service.ConfigService, tracer trace.Tracer) ConfigHandler {
	return ConfigHandler{
		service: service,
		tracer:  tracer,
	}
}

// @Summary Add a new configuration
// @Description Adds a new configuration
// @Tags configs
// @Accept json
// @Produce json
// @Param config body model.Config true "Configuration to add"
// @Success 200 {object} model.Config
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /configs [post]
func (c ConfigHandler) AddConfig(writer http.ResponseWriter, request *http.Request) {
	ctx, span := c.tracer.Start(request.Context(), "AddConfig")
	defer span.End()

	defer request.Body.Close()

	var config model.Config

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&config)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}
	if config.Name == "" || config.Version == 0 || len(config.Parameters) == 0 || len(config.Labels) == 0 {
		http.Error(writer, "Error: 'name', 'version', 'params', and 'labels' fields are required and cannot be empty", http.StatusBadRequest)
		return
	}

	existingConfig, err := c.service.GetConfig(ctx, config.Name, config.Version)
	if err != nil && err.Error() != "config not found" {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if existingConfig.Name != "" {
		http.Error(writer, "Configuration with the given name and version already exists", http.StatusConflict)
		return
	}

	err = c.service.AddConfig(ctx, config)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(config)
}

// @Summary Get a configuration
// @Description Retrieves a configuration by name and version
// @Tags configs
// @Produce json
// @Param name path string true "Name of the configuration"
// @Param version path int true "Version of the configuration"
// @Success 200 {object} model.Config
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Configuration not found"
// @Failure 500 {string} string "Internal server error"
// @Router /configs/{name}/{version} [get]
func (c ConfigHandler) GetConfig(writer http.ResponseWriter, request *http.Request) {
	ctx, span := c.tracer.Start(request.Context(), "GetConfig")
	defer span.End()

	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, "Invalid version format", http.StatusBadRequest)
		return
	}

	config, err := c.service.GetConfig(ctx, name, versionInt)
	if err != nil {
		if err.Error() == "config not found" {
			http.Error(writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(config)
}

// @Summary Delete a configuration
// @Description Deletes a configuration by name and version
// @Tags configs
// @Param name path string true "Name of the configuration"
// @Param version path int true "Version of the configuration"
// @Success 200 {string} string "Successfully deleted configuration"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Configuration not found"
// @Failure 500 {string} string "Internal server error"
// @Router /configs/{name}/{version} [delete]
func (c ConfigHandler) DeleteConfig(writer http.ResponseWriter, request *http.Request) {
	ctx, span := c.tracer.Start(request.Context(), "DeleteConfig")
	defer span.End()

	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, "Invalid version format", http.StatusBadRequest)
		return
	}

	err = c.service.DeleteConfig(ctx, name, versionInt)
	if err != nil {
		if err.Error() == "config not found" {
			http.Error(writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{"message": "Configuration successfully deleted"}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}
