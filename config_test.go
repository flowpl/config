package config_test

import (
	"testing"
	"github.com/flowpl/config"
	"errors"
	"reflect"
	"os"
)

var getBooleanValidValues = []struct{
	input string
	expected bool
} {
	{input: "true", expected: true},
	{input: "TRUE", expected: true},
	{input: "1", expected: true},
	{input: "yes", expected: true},
	{input: "YES", expected: true},
	{input: "y", expected: true},
	{input: "Y", expected: true},

	{input: "false", expected: false},
	{input: "FALSE", expected: false},
	{input: "0", expected: false},
	{input: "no", expected: false},
	{input: "NO", expected: false},
	{input: "n", expected: false},
	{input: "N", expected: false},
}

func TestConfig_GetBool_ShouldReturnValidValuesForValidInputs(t *testing.T) {
	for _, testData := range getBooleanValidValues {
		c := config.NewConfig("app", func(name []string) (interface{}, error) {
			return testData.input, nil
		})

		result, _ := c.GetBool("someName")

		if result != testData.expected {
			t.Errorf("expected bool result for %s to be %t. Actual: %t", testData.input, testData.expected, result)
		}
	}
}

func TestConfig_GetBool_ShouldReturnAnErrorIfTheLoaderReturnsAnError(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return "", errors.New("testerror")
	})

	_, err := c.GetBool("someName")

	if err == nil {
		t.Error("expected an error to be returned, but it wasn't")
	}

	if len(c.GetErrors()) == 0 {
		t.Error("expected error to be appended to the Config error list")
	}
}

func TestConfig_GetBool_ShouldReturnAnErrorIfTheConfigStringCannotBeParsedToBool(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return "NotABool", nil
	})

	_, err := c.GetBool("someName")

	if err == nil {
		t.Error("expected an error to be returned, but it wasn't")
	}

	if len(c.GetErrors()) == 0 {
		t.Error("expected error to be appended to the Config error list")
	}
}

func TestConfig_GetBool_ShouldReturnTheValueDirectlyIfLoaderReturnsBool(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return true, nil
	})

	result, _ := c.GetBool("someName")

	if result != true {
		t.Errorf("expected bool result for true to be true. Actual: %t", result)
	}
}

func TestConfig_GetBool_ShouldReturnAnErrorIfLoaderNeitherReturnsStringNorBool(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return 65, nil
	})

	_, err := c.GetBool("someName")

	if err == nil {
		t.Error("expected an error to be returned, but it wasn't")
	}

	if len(c.GetErrors()) == 0 {
		t.Error("expected error to be appended to the Config error list")
	}
}

func TestConfig_GetBool_ShouldPassTheCorrectPathComponentsToTheLoader(t *testing.T) {
	var passedPath []string
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		passedPath = name
		return true, nil
	})

	c.GetBool("someName")

	if !reflect.DeepEqual(passedPath, []string{"app", "someName"}) {
		t.Error("expected path to be []string{\"app\",\"someName\"}, but it wasn't")
	}
}

func TestConfig_MayGetBool_ShouldReturnTheSameValueAsGetBool(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return "true", nil
	})

	result, _ := c.GetBool("someName")
	result2 := c.MayGetBool("someName")

	if result != result2 {
		t.Error("expected return values from GetBool() and MayGetBool() to be identical")
	}
}

func TestConfig_GetString_ShouldReturnAValueDirectlyFromLoader(t *testing.T) {
	const testValue = "someValue"
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return testValue, nil
	})

	result, _ := c.GetString("someName")

	if result != testValue {
		t.Errorf("expected result to be \"%s\", actual: \"%s\"", testValue, result)
	}
}

func TestConfig_GetString_ShouldReturnAnErrorIfLoaderReturnsAnError(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return "", errors.New("testerror")
	})

	_, err := c.GetString("someName")

	if err == nil {
		t.Error("expected an error to be returned, but it wasn't")
	}

	if len(c.GetErrors()) == 0 {
		t.Error("expected error to be appended to the Config error list")
	}
}

func TestConfig_GetString_ShouldReturnAnErrorIfLoaderDoesNotReturnString(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return 65, nil
	})

	_, err := c.GetString("someName")

	if err == nil {
		t.Error("expected an error to be returned, but it wasn't")
	}

	if len(c.GetErrors()) == 0 {
		t.Error("expected error to be appended to the Config error list")
	}
}

func TestConfig_GetString_ShouldPassTheCorrectPathComponentsToTheLoader(t *testing.T) {
	var passedPath []string
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		passedPath = name
		return "string", nil
	})

	c.GetString("someName")

	if !reflect.DeepEqual(passedPath, []string{"app", "someName"}) {
		t.Error("expected path to be []string{\"app\",\"someName\"}, but it wasn't")
	}
}

func TestConfig_MayGetString_ShouldReturnTheSameValueAsGetString(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return "true", nil
	})

	result, _ := c.GetString("someName")
	result2 := c.MayGetString("someName")

	if result != result2 {
		t.Error("expected return values from GetString() and MayGetString() to be identical")
	}
}

var getIntValidStrings = []struct{
	input string
	expected int
} {
	{input: "1", expected: 1},
	{input: "0", expected: 0},
	{input: "-1", expected: -1},
	{input: "2000000000", expected: 2000000000},
	{input: "-2000000000", expected:-2000000000},
}

func TestConfig_GetInt_ShouldReturnAValidValueForValidInputString(t *testing.T) {
	for _, testData := range getIntValidStrings {
		c := config.NewConfig("app", func(name []string) (interface{}, error) {
			return testData.input, nil
		})

		result, _ := c.GetInt("someName")

		if result != testData.expected {
			t.Errorf("expected result for %s to be %d, actual %d", testData.input, testData.expected, result)
		}
	}
}

func TestConfig_GetInt_ShouldReturnAValueDirectlyFromLoaderIfItsInt(t *testing.T) {
	const testValue = 32
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return testValue, nil
	})

	result, _ := c.GetInt("someName")

	if result != testValue {
		t.Errorf("expected result for %d to be %d, actual %d", testValue, testValue, result)
	}
}

func TestConfig_GetInt_ShouldReturnErrorIfLoaderReturnsError(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return "", errors.New("testerror")
	})

	_, err := c.GetInt("someName")

	if err == nil {
		t.Error("expected an error to be returned, but it wasn't")
	}

	if len(c.GetErrors()) == 0 {
		t.Error("expected error to be appended to the Config error list")
	}
}

func TestConfig_GetInt_ShouldReturnErrorIfLoaderReturnsAnUnsupportedType(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return true, nil
	})

	_, err := c.GetInt("someName")

	if err == nil {
		t.Error("expected an error to be returned, but it wasn't")
	}

	if len(c.GetErrors()) == 0 {
		t.Error("expected error to be appended to the Config error list")
	}
}

func TestConfig_GetInt_ShouldReturnErrorIfLoaderReturnsUnparseableString(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return "true", nil
	})

	_, err := c.GetInt("someName")

	if err == nil {
		t.Error("expected an error to be returned, but it wasn't")
	}

	if len(c.GetErrors()) == 0 {
		t.Error("expected error to be appended to the Config error list")
	}
}

func TestConfig_GetInt_ShouldPassTheCorrectPathComponentsToTheLoader(t *testing.T) {
	var passedPath []string
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		passedPath = name
		return 32, nil
	})

	c.GetInt("someName")

	if !reflect.DeepEqual(passedPath, []string{"app", "someName"}) {
		t.Error("expected path to be []string{\"app\",\"someName\"}, but it wasn't")
	}
}

func TestConfig_MayGetInt_ShouldReturnTheSameValueAsGetInt(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return "1", nil
	})

	result, _ := c.GetInt("someName")
	result2 := c.MayGetInt("someName")

	if result != result2 {
		t.Error("expected return values from GetInt() and MayGetInt() to be identical")
	}
}

var getFloatValidStrings = []struct{
	input string
	expected float64
} {
	{input: "1.0", expected: 1},
	{input: "0", expected: 0},
	{input: "-1.0", expected: -1},
	{input: "2000000000", expected: 2000000000},
	{input: "-2000000000", expected:-2000000000},
	{input: "15.98754379", expected:15.98754379},
	{input: "-15.98754379", expected:-15.98754379},
}

func TestConfig_GetFloat_ShouldReturnAValidValueForValidInputString(t *testing.T) {
	for _, testData := range getFloatValidStrings {
		c := config.NewConfig("app", func(name []string) (interface{}, error) {
			return testData.input, nil
		})

		result, _ := c.GetFloat("someName")

		if result != testData.expected {
			t.Errorf("expected result for %s to be %d, actual %d", testData.input, testData.expected, result)
		}
	}
}

func TestConfig_GetFloat_ShouldReturnAValueDirectlyFromLoaderIfItsFloat64(t *testing.T) {
	const testValue float64 = 4.5
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return testValue, nil
	})

	result, _ := c.GetFloat("someName")

	if result != testValue {
		t.Errorf("expected result for %d to be %d, actual %d", testValue, testValue, result)
	}
}

func TestConfig_GetFloat_ShouldCastToFloat64IfLoaderReturnsInt(t *testing.T) {
	const testValue int = 4
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return testValue, nil
	})

	result, _ := c.GetFloat("someName")

	floatTestValue := float64(testValue)
	if floatTestValue != result {
		t.Errorf("expected result for %d to be %d, actual %d", testValue, testValue, result)
	}
}

func TestConfig_GetFloat_ShouldCastToFloat64IfLoaderReturnsFloat32(t *testing.T) {
	const testValue float32 = 4
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return testValue, nil
	})

	result, _ := c.GetFloat("someName")

	floatTestValue := float64(testValue)
	if floatTestValue != result {
		t.Errorf("expected result for %d to be %d, actual %d", testValue, testValue, result)
	}
}

func TestConfig_GetFloat_ShouldReturnErrorIfLoaderReturnsError(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return "", errors.New("testerror")
	})

	_, err := c.GetFloat("someName")

	if err == nil {
		t.Error("expected an error to be returned, but it wasn't")
	}

	if len(c.GetErrors()) == 0 {
		t.Error("expected error to be appended to the Config error list")
	}
}

func TestConfig_GetFloat_ShouldReturnErrorIfLoaderReturnsAnUnsupportedType(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return true, nil
	})

	_, err := c.GetFloat("someName")

	if err == nil {
		t.Error("expected an error to be returned, but it wasn't")
	}

	if len(c.GetErrors()) == 0 {
		t.Error("expected error to be appended to the Config error list")
	}
}

func TestConfig_GetFloat_ShouldReturnErrorIfLoaderReturnsUnparseableString(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return "true", nil
	})

	_, err := c.GetFloat("someName")

	if err == nil {
		t.Error("expected an error to be returned, but it wasn't")
	}

	if len(c.GetErrors()) == 0 {
		t.Error("expected error to be appended to the Config error list")
	}
}

func TestConfig_GetFloat_ShouldPassTheCorrectPathComponentsToTheLoader(t *testing.T) {
	var passedPath []string
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		passedPath = name
		return 3.4, nil
	})

	c.GetFloat("someName")

	if !reflect.DeepEqual(passedPath, []string{"app", "someName"}) {
		t.Error("expected path to be []string{\"app\",\"someName\"}, but it wasn't")
	}
}

func TestConfig_MayGetFloat_ShouldReturnTheSameValueAsGetFloat(t *testing.T) {
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		return "true", nil
	})

	result, _ := c.GetFloat("someName")
	result2 := c.MayGetFloat("someName")

	if result != result2 {
		t.Error("expected return values from GetFloat() and MayGetFloat() to be identical")
	}
}

func TestConfig_ShouldPassAFourLevelPathIfGetChildHasBeenCalledTwice(t *testing.T) {
	var passedPath []string
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		passedPath = name
		return "0", nil
	})
	child1 := c.ChildConfig("one")
	child2 := child1.ChildConfig("two")

	child2.GetBool("someName")
	if !reflect.DeepEqual(passedPath, []string{"app", "one", "two", "someName"}) {
		t.Error("expected path to be []string{\"app\",\"one\",\"two\",\"someName\"} in GetBool(), but it wasn't")
	}

	child2.GetString("someName")
	if !reflect.DeepEqual(passedPath, []string{"app", "one", "two", "someName"}) {
		t.Error("expected path to be []string{\"app\",\"one\",\"two\",\"someName\"} in GetString(), but it wasn't")
	}

	child2.GetInt("someName")
	if !reflect.DeepEqual(passedPath, []string{"app", "one", "two", "someName"}) {
		t.Error("expected path to be []string{\"app\",\"one\",\"two\",\"someName\"} in GetInt(), but it wasn't")
	}

	child2.GetFloat("someName")
	if !reflect.DeepEqual(passedPath, []string{"app", "one", "two", "someName"}) {
		t.Error("expected path to be []string{\"app\",\"one\",\"two\",\"someName\"} in GetFloat(), but it wasn't")
	}
}

func TestConfig_ShouldPassATwoLevelPathIfChildHasBeenCalledOnceButConfigLookupIsOnTheOriginalConfig(t *testing.T) {
	var passedPath []string
	c := config.NewConfig("app", func(name []string) (interface{}, error) {
		passedPath = name
		return "0", nil
	})
	child1 := c.ChildConfig("one")

	child1.GetBool("someName")
	if !reflect.DeepEqual(passedPath, []string{"app", "one", "someName"}) {
		t.Error("expected path to be []string{\"app\",\"one\",\"someName\"} in GetBool(), but it wasn't")
	}

	c.GetBool("someName")
	if !reflect.DeepEqual(passedPath, []string{"app", "someName"}) {
		t.Error("expected path to be []string{\"app\",\"someName\"} in GetBool(), but it wasn't")
	}
}

func TestConfig_ShouldLoadAnActualVariableFromTheEnvironment(t *testing.T) {
	os.Setenv("APP_ONE_VALUE", "success")
	c := config.NewConfig("app", config.NewEnvironmentLoader())
	child1 := c.ChildConfig("one")

	result, err := child1.GetString("value")
	if err != nil {
		t.Error("an error was not expected to be returned but it was")
	}

	if result != "success" {
		t.Errorf("expected value for name %s to be %s, actual %s", "value", "success", result)
	}
}
