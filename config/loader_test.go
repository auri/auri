package config

import (
	"reflect"
	"testing"

	"github.com/gobuffalo/envy"
)

type TestConfig struct {
	TestString      string `envyname:"TEST_STRING" default:"DefaultTestString"`
	TestBool        bool
	TestInt         int      `envyname:"TEST_INT"`
	TestEmptyString string   `envyname:"TEST_EMPTY_STRING" default:""`
	TestUInt        uint     `envyname:"TEST_UINT"`
	TestSlice       []string `envyname:"TEST_SLICE"`
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		testConfig TestConfig
		envyVars   map[string]string
		wantErr    string // empty for no errors
	}{
		{
			"Valid values",
			TestConfig{"teststringValue", true, 20, "", 10, []string{"element1", "element2"}},
			map[string]string{
				"TEST_STRING": "teststringValue",
				"TESTBOOL":    "true",
				"TEST_INT":    "20",
				"TEST_UINT":   "10",
				"TEST_SLICE":  "element1, element2",
			},
			"",
		},
		{
			"Empty string",
			TestConfig{},
			map[string]string{
				"TEST_STRING": "",
				"TESTBOOL":    "true",
				"TEST_INT":    "20",
			},
			"Wrong value for config parameter TEST_STRING: '' - should not be empty",
		},
		{
			"Broken boolean",
			TestConfig{},
			map[string]string{
				"TEST_STRING": "teststringValue",
				"TESTBOOL":    "ABC",
				"TEST_INT":    "20",
			},
			"Wrong value for config parameter TESTBOOL: 'ABC' - should be true/false",
		},
		{
			"Broken int",
			TestConfig{},
			map[string]string{
				"TEST_STRING": "teststringValue",
				"TESTBOOL":    "true",
				"TEST_INT":    "ABC",
			},
			"Wrong value for config parameter TEST_INT: 'ABC' - should be a number",
		},
		{
			"Broken uint",
			TestConfig{},
			map[string]string{
				"TEST_STRING": "teststringValue",
				"TESTBOOL":    "true",
				"TEST_INT":    "20",
				"TEST_UINT":   "-10",
			},
			"Wrong value for config parameter TEST_UINT: '-10' - should be a positive number",
		},
		{
			"Config option without default value is missing",
			TestConfig{},
			map[string]string{
				"TEST_STRING": "teststringValue",
				"TEST_INT":    "20",
			},
			"Wrong value for config parameter TESTBOOL: '' - should be set",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance = &TestConfig{}
			envy.Reload()
			for k, v := range tt.envyVars {
				envy.Set(k, v)
			}

			err := LoadConfig()

			if (err != nil) != (tt.wantErr != "") {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != "" && (err != nil) {
				if !reflect.DeepEqual(err.Error(), tt.wantErr) {
					t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				configStruct := *(instance.(*TestConfig))
				if !reflect.DeepEqual(configStruct, tt.testConfig) {
					t.Errorf("configStruct = \n%#v, want \n%#v", configStruct, tt.testConfig)
				}
			}

		})
	}
}

func TestGetInstance(t *testing.T) {
	instance = &Config{}

	if got := GetInstance(); !reflect.DeepEqual(got, Config{}) {
		t.Errorf("GetInstance() = %v, want %v", got, Config{})
	}
}
