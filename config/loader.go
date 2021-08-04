package config

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
)

// errorLogWrongValue contains the the error message
// returned for wrong configuration values
const errorLogWrongValue = "Wrong value for config parameter %v: '%v' - %s"

// ConfigDir defines the path to the configuration directory of auri
var ConfigDir = "/etc/auri"

// auriConfigFile is the filename of auri configuration
var auriConfigFile = "config.env"

// prodBuild is set to yes during the build process of production builds
var prodBuild = "no"
var instance interface{}

// IsProdBuild returns true if this is a prod build
func IsProdBuild() bool {
	return prodBuild == "yes"
}

// GetInstance provides the current validated configuration
func GetInstance() Config {
	if instance == nil {
		if err := LoadConfig(); err != nil {
			panic(fmt.Sprintf("Fatal error: %v", err))
		}
	}
	// copy of instance value to avoid any manipulations of config
	return *(instance.(*Config))
}

// LoadConfig loads the configuration
func LoadConfig() error {
	if prodBuild == "yes" {
		if err := envy.Load(filepath.Join(ConfigDir, auriConfigFile)); err != nil {
			return errors.WithMessage(err, "Failed to load configuration")
		}
	}
	envy.Load()

	if instance == nil {
		instance = &Config{}
	}

	v := reflect.ValueOf(instance).Elem()
	t := reflect.TypeOf(instance).Elem()

	for i := 0; i < v.NumField(); i++ {
		fieldTag := t.Field(i).Tag
		// get the envy config name from tag property, use uppercase notation otherwise
		var configOpt string
		if envynameTag, ok := fieldTag.Lookup("envyname"); ok {
			configOpt = envynameTag
		} else {
			configOpt = strings.ToUpper(t.Field(i).Name)
		}

		fieldVal := v.Field(i)
		defaultValPresent := false
		var valueDefault string
		var value string

		// load default values
		valueDefault, defaultValPresent = fieldTag.Lookup("default")

		// get the config value from envy if possible, set default value otherwise
		if valueEnvy, err := envy.MustGet(configOpt); err == nil {
			// verify if we got an empty value for a required option (without default value)
			if len(valueEnvy) == 0 && (!defaultValPresent || len(valueDefault) != 0) {
				return errors.Errorf(errorLogWrongValue, configOpt, valueEnvy, "should not be empty")
			}

			value = valueEnvy
		} else {
			if defaultValPresent {
				value = valueDefault
			} else {
				return errors.Errorf(errorLogWrongValue, configOpt, value, "should be set")
			}
		}

		// verify if we get a value corresponding to our expected type
		fieldValKind := fieldVal.Kind()
		switch fieldValKind {
		case reflect.String:
			fieldVal.SetString(value)
		case reflect.Bool:
			s, err := strconv.ParseBool(value)
			if err != nil {
				return errors.Errorf(errorLogWrongValue, configOpt, value, "should be true/false")
			}

			fieldVal.SetBool(s)
		case reflect.Int:
			s, err := strconv.ParseInt(value, 10, 0)
			if err != nil {
				return errors.Errorf(errorLogWrongValue, configOpt, value, "should be a number")
			}

			fieldVal.SetInt(s)
		case reflect.Uint:
			s, err := strconv.ParseUint(value, 10, 0)
			if err != nil {
				return errors.Errorf(errorLogWrongValue, configOpt, value, "should be a positive number")
			}

			fieldVal.SetUint(s)
		case reflect.Slice:
			elemKind := fieldVal.Type().Elem().Kind()
			switch elemKind {
			case reflect.String:
				var res []string
				//do some cleanup on the resulting slice
				for _, el := range strings.Split(value, ",") {
					res = append(res, strings.Trim(el, " "))
				}
				fieldVal.Set(reflect.ValueOf(res))
			default:
				return errors.Errorf("Validation failed: unhandled element type %v for config parameter %v", elemKind, configOpt)
			}
		default:
			return errors.Errorf("Validation failed: unhandled type %v for config parameter %v", fieldValKind, configOpt)
		}
	}

	return nil
}
