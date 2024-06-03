package tests

import (
	"projekat/model"
	"testing"
)

func TestConfigFields(t *testing.T) {
	validConfig := model.Config{
		Name:       "config1",
		Version:    1,
		Parameters: map[string]string{"param1": "value1"},
		Labels:     map[string]string{"label1": "value1"},
	}

	if validConfig.Name != "config1" {
		t.Errorf("Expected Name to be 'config1', got '%s'", validConfig.Name)
	}

	if validConfig.Version != 1 {
		t.Errorf("Expected Version to be 1, got '%d'", validConfig.Version)
	}

	if len(validConfig.Parameters) != 1 || validConfig.Parameters["param1"] != "value1" {
		t.Errorf("Expected Parameters to have one entry with key 'param1' and value 'value1'")
	}

	if len(validConfig.Labels) != 1 || validConfig.Labels["label1"] != "value1" {
		t.Errorf("Expected Labels to have one entry with key 'label1' and value 'value1'")
	}
}

func TestInvalidConfig(t *testing.T) {
	invalidConfig := model.Config{
		Name:       "",
		Version:    0,
		Parameters: map[string]string{},
		Labels:     map[string]string{},
	}

	if invalidConfig.Name != "" {
		t.Errorf("Expected Name to be empty, got '%s'", invalidConfig.Name)
	}

	if invalidConfig.Version != 0 {
		t.Errorf("Expected Version to be 0, got '%d'", invalidConfig.Version)
	}

	if len(invalidConfig.Parameters) != 0 {
		t.Errorf("Expected Parameters to be empty")
	}

	if len(invalidConfig.Labels) != 0 {
		t.Errorf("Expected Labels to be empty")
	}
}

func TestConfigEquality(t *testing.T) {
	config1 := model.Config{
		Name:       "config1",
		Version:    1,
		Parameters: map[string]string{"param1": "value1"},
		Labels:     map[string]string{"label1": "value1"},
	}

	config2 := model.Config{
		Name:       "config1",
		Version:    1,
		Parameters: map[string]string{"param1": "value1"},
		Labels:     map[string]string{"label1": "value1"},
	}

	if config1.Name != config2.Name ||
		config1.Version != config2.Version ||
		!compareMaps(config1.Parameters, config2.Parameters) ||
		!compareMaps(config1.Labels, config2.Labels) {
		t.Errorf("Expected config1 and config2 to be equal")
	}
}

func compareMaps(map1, map2 map[string]string) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, value := range map1 {
		if map2[key] != value {
			return false
		}
	}
	return true
}

func TestConfigParameterUpdate(t *testing.T) {
	config := model.Config{
		Name:       "config1",
		Version:    1,
		Parameters: map[string]string{"param1": "value1"},
		Labels:     map[string]string{"label1": "value1"},
	}

	config.Parameters["param2"] = "value2"

	if len(config.Parameters) != 2 {
		t.Errorf("Expected Parameters to have two entries")
	}

	if config.Parameters["param2"] != "value2" {
		t.Errorf("Expected Parameters to have key 'param2' with value 'value2'")
	}
}

func TestConfigLabelUpdate(t *testing.T) {
	config := model.Config{
		Name:       "config1",
		Version:    1,
		Parameters: map[string]string{"param1": "value1"},
		Labels:     map[string]string{"label1": "value1"},
	}

	config.Labels["label2"] = "value2"

	if len(config.Labels) != 2 {
		t.Errorf("Expected Labels to have two entries")
	}

	if config.Labels["label2"] != "value2" {
		t.Errorf("Expected Labels to have key 'label2' with value 'value2'")
	}
}
