package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"projekat/handlers"
	"projekat/model"
	"projekat/service"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
)

// Testni repozitorijum za testiranje
type TestConfigRepository struct{}

func (r *TestConfigRepository) AddConfig(ctx context.Context, config model.Config) error {
	return nil
}

func (r *TestConfigRepository) GetConfig(ctx context.Context, name string, version int) (model.Config, error) {
	if name == "testConfig" && version == 1 {
		return model.Config{
			Name:       "testConfig",
			Version:    1,
			Parameters: map[string]string{"param1": "value1"},
			Labels:     map[string]string{"label1": "value1"},
		}, nil
	}
	return model.Config{}, fmt.Errorf("config not found")
}

func (r *TestConfigRepository) DeleteConfig(ctx context.Context, name string, version int) error {
	return nil
}

func TestGetConfig(t *testing.T) {
	// Postavljanje testnog repozitorijuma i servisa
	repo := &TestConfigRepository{}
	configService := service.NewConfigService(repo, otel.Tracer("test-tracer"))
	handler := handlers.NewConfigHandler(configService, otel.Tracer("test-tracer"))

	// Priprema request-a
	req, err := http.NewRequest("GET", "/configs/testConfig/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Postavljanje mux varijabli u request
	req = mux.SetURLVars(req, map[string]string{
		"name":    "testConfig",
		"version": strconv.Itoa(1),
	})

	// Priprema response recorder-a
	rr := httptest.NewRecorder()

	// Pozivanje handler-a
	handler.GetConfig(rr, req)

	// Provera rezultata
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	var returnedConfig model.Config
	err = json.NewDecoder(rr.Body).Decode(&returnedConfig)
	if err != nil {
		t.Fatal(err)
	}
	expectedConfig := model.Config{
		Name:       "testConfig",
		Version:    1,
		Parameters: map[string]string{"param1": "value1"},
		Labels:     map[string]string{"label1": "value1"},
	}
	assert.Equal(t, expectedConfig, returnedConfig, "Expected returned config to match input config")
}

func TestDeleteConfig(t *testing.T) {
	// Postavljanje testnog repozitorijuma i servisa
	repo := &TestConfigRepository{}
	configService := service.NewConfigService(repo, otel.Tracer("test-tracer"))
	handler := handlers.NewConfigHandler(configService, otel.Tracer("test-tracer"))
	t.Fail()

	// Priprema request-a
	req, err := http.NewRequest("DELETE", "/configs/testConfig/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Postavljanje mux varijabli u request
	req = mux.SetURLVars(req, map[string]string{
		"name":    "testConfig",
		"version": strconv.Itoa(1),
	})

	// Priprema response recorder-a
	rr := httptest.NewRecorder()

	// Pozivanje handler-a
	handler.DeleteConfig(rr, req)

	// Provera rezultata
	assert.Equal(t, http.StatusOK, rr.Code, "Greska")
	assert.Contains(t, rr.Body.String(), "Configuration successfully deleted", "Expected success message in response")
}

func TestAddConfig(t *testing.T) {
	// Postavljanje testnog repozitorijuma i servisa
	repo := &TestConfigRepository{}
	configService := service.NewConfigService(repo, otel.Tracer("test-tracer"))
	handler := handlers.NewConfigHandler(configService, otel.Tracer("test-tracer"))

	// Priprema podataka za testiranje
	config := model.Config{
		Name:       "testConfig2",
		Version:    11,
		Parameters: map[string]string{"param1": "value1"},
		Labels:     map[string]string{"label1": "value1"},
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	// Priprema HTTP POST zahteva
	req, err := http.NewRequest("POST", "/configs", bytes.NewBuffer(configJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Priprema response recorder-a
	rr := httptest.NewRecorder()

	// Pozivanje handler-a
	handler.AddConfig(rr, req)

	// Provera rezultata
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	var returnedConfig model.Config
	err = json.NewDecoder(rr.Body).Decode(&returnedConfig)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, config, returnedConfig, "Expected returned config to match input config")
}
