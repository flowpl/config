package config

import (
	"strings"
	"errors"
	"fmt"
	"strconv"
	"syscall"
)

var ACCEPTABLE_TRUETHY []string = []string{"yes", "1", "true", "y"}
var ACCEPTABLE_FALSY []string = []string{"no", "0", "false", "n"}

/**
 * for a list of path components load the value to be set that path or an error if the path cannot be resolved
 */
type Loader func(pathElements []string) (interface{}, error)

type Config struct {
	loader Loader
	pathElements []string
	errors []error
}

func (c *Config) GetErrors() []error {
	return c.errors
}

func (c *Config) TriggerErrorPanic() {
	if len(c.errors) > 0 {
		errorMsg := ""
		for _, err := range c.errors {
			errorMsg += err.Error() + "\n"
		}
		panic(strings.TrimSpace(errorMsg))
	}
}

func (c *Config) ChildConfig(name string) (*Config) {
	return &Config{
		loader: c.loader,
		pathElements: append(c.pathElements, name),
	}
}

func (c *Config) GetString(name string) (string, error) {
	value, loadErr := c.loader(append(c.pathElements, name))
	if loadErr != nil {
		c.errors = append(c.errors, loadErr)
		return "", loadErr
	}

	switch t := value.(type) {
	case string:
		return t, nil
	}

	err := errors.New(fmt.Sprintf("value is not a string for name %s", name))
	c.errors = append(c.errors, err)
	return "", err
}

func (c *Config) MayGetString(name string) string {
	res, _ := c.GetString(name)
	return res
}

func (c *Config) GetBool(name string) (bool, error) {
	value, loadErr := c.loader(append(c.pathElements, name))
	if loadErr != nil {
		c.errors = append(c.errors, loadErr)
		return false, loadErr
	}

	switch t := value.(type) {
	case string:
		if inArray(value.(string), ACCEPTABLE_TRUETHY) {
			return true, nil
		} else if inArray(value.(string), ACCEPTABLE_FALSY) {
			return false, nil
		}
	case bool:
		return t, nil
	}

	err := errors.New(fmt.Sprintf("value is not a valid boolean or parseable string for name %s", name))
	c.errors = append(c.errors, err)
	return false, err
}

func (c *Config) MayGetBool(name string) bool {
	res, _ := c.GetBool(name)
	return res
}

func (c *Config) GetInt(name string) (int, error) {
	value, loadErr := c.loader(append(c.pathElements, name))
	if loadErr != nil {
		c.errors = append(c.errors, loadErr)
		return 0, loadErr
	}

	switch t := value.(type) {
	case string:
		intval, parseErr := strconv.Atoi(value.(string))
		if parseErr != nil {
			c.errors = append(c.errors, parseErr)
			return 0, parseErr
		}
		return intval, nil
	case int:
		return t, nil
	}

	err := errors.New("value is not an int or parseable string for name " + name)
	c.errors = append(c.errors, err)
	return 0, err
}

func (c *Config) MayGetInt(name string) int {
	res, _ := c.GetInt(name)
	return res
}

func (c *Config) GetFloat(name string) (float64, error) {
	value, loadErr := c.loader(append(c.pathElements, name))
	if loadErr != nil {
		c.errors = append(c.errors, loadErr)
		return 0, loadErr
	}

	switch t := value.(type) {
	case string:
		floatVal, parseErr := strconv.ParseFloat(t, 64)
		if parseErr != nil {
			c.errors = append(c.errors, parseErr)
			return 0, parseErr
		}
		return floatVal, nil
	case float32:
		return float64(t), nil
	case float64:
		return t, nil
	case int:
		return float64(t), nil
	}

	err := errors.New("value is not a compatible float (32, 64), int or parseable string for name " + name)
	c.errors = append(c.errors, err)
	return 0, err
}

func (c *Config) MayGetFloat(name string) float64 {
	res, _ := c.GetFloat(name)
	return res
}

func NewConfig(applicationName string, loader Loader) *Config {
	return &Config{loader: loader, pathElements: []string{applicationName}}
}

func NewEnvironmentLoader() Loader {
	return func(pathElements []string) (interface{}, error) {
		variableName := strings.ToUpper(strings.Join(pathElements, "_"))
		value, exists := syscall.Getenv(variableName)
		if !exists {
			return "", errors.New(fmt.Sprintf("environment variable %s does not exist", variableName))
		}

		return value, nil
	}
}

//func NewJsonLoader(filePath string) (Loader, error) {
//	fileContents, loadErr := ioutil.ReadFile(filePath)
//	if loadErr != nil {
//		return nil, loadErr
//	}
//
//	parsedContents := new(interface{})
//	parseErr := json.Unmarshal(fileContents, parsedContents)
//	if parseErr != nil {
//		return nil, parseErr
//	}
//
//	return func(pathElements []string) (interface{}, error) {
//
//	}
//}

func inArray(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == strings.ToLower(needle) {
			return true
		}
	}
	return false
}
